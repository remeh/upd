// Instance of the server.
// Copyright © 2015 - Rémy MATHIEU

package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	Config   Config    // Configuration
	Metadata Metadatas // Link to the read metadata
}

func NewServer(config Config) *Server {
	// init the random
	rand.Seed(time.Now().Unix())

	return &Server{
		Config:   config,
		Metadata: Metadatas{CreationTime: time.Now(), Data: make(map[string]Metadata)},
	}
}

// Starts the listening daemon.
func (s *Server) Start() {
	router := s.prepareRouter()

	// Setup the router on the net/http stack
	http.Handle("/", router)

	// Read the existing metadata.
	s.readMetadata()

	go s.StartCleanJob()

	// Listen
	if len(s.Config.CertificateFile) != 0 && len(s.Config.CertificateKey) != 0 {
		log.Println("[info] Start secure listening on", s.Config.Addr)
		err := http.ListenAndServeTLS(s.Config.Addr, s.Config.CertificateFile, s.Config.CertificateKey, nil)
		log.Println("[err]", err.Error())
	} else {
		log.Println("[info] Start listening on", s.Config.Addr)
		err := http.ListenAndServe(s.Config.Addr, nil)
		log.Println("[err]", err.Error())
	}
}

// Starts the Clean Job
func (s *Server) StartCleanJob() {
	timer := time.NewTicker(60 * time.Second)
	for _ = range timer.C {
		job := CleanJob{s}
		job.Run()
	}
}

// Reads the stored metadata.
func (s *Server) readMetadata() {
	file, err := os.Open(s.Config.OutputDirectory + "/metadata.json")
	create := false
	if err != nil {
		create = true
	}

	if create {
		// Create the file
		log.Println("[info] Creating metadata.json")
		s.writeMetadata(true)
	} else {
		// Read the file
		log.Println("[info] Reading metadata.json")

		readData, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println("[err] The existing metadata.json seems corrupted. Exiting.")
			log.Println(err)
			os.Exit(1)
		}

		var data Metadatas
		json.Unmarshal(readData, &data)
		s.Metadata = data
		log.Printf("[info] %d metadata read.\n", len(s.Metadata.Data))
		file.Close()
	}
}

func (s *Server) writeMetadata(printLog bool) {
	// TODO mutex!!
	file, err := os.Create(s.Config.OutputDirectory + "/metadata.json")
	if err != nil {
		log.Println("[err] Can't write in the output directory:", s.Config.OutputDirectory)
		log.Println(err)
		os.Exit(1)
	}
	data, _ := json.Marshal(s.Metadata)
	file.Write(data)
	file.Close()

	if printLog {
		log.Printf("[info] %d metadatas written.\n", len(s.Metadata.Data))
	}
}

// Prepares the route
func (s *Server) prepareRouter() *mux.Router {
	r := mux.NewRouter()

	sendHandler := &SendHandler{s}
	r.Handle(s.Config.Route+"/1.0/send", sendHandler)

	lastUploadeHanlder := &LastUploadedHandler{s}
	r.Handle(s.Config.Route+"/1.0/list", lastUploadeHanlder)

	deleteHandler := &DeleteHandler{s}
	r.Handle(s.Config.Route+"/{file}/{key}", deleteHandler)

	sh := &ServingHandler{s}
	r.Handle(s.Config.Route+"/{file}", sh) // Serving route.

	return r
}

// Expire expires a file : delete it from the metadata
// and from the FS.
func (s *Server) Expire(m Metadata) error {
	delete(s.Metadata.Data, m.Filename)
	return os.Remove(s.Config.OutputDirectory + "/" + m.Filename)
}
