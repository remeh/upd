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
	Flags    Flags     // Configuration
	Metadata Metadatas // Link to the read metadata
}

func NewServer(flags Flags) *Server {
	// init the random
	rand.Seed(time.Now().Unix())

	return &Server{
		Flags:    flags,
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

	// Listen

	if len(s.Flags.CertificateFile) != 0 && len(s.Flags.CertificateKey) != 0 {
		log.Println("[info] Start secure listening on", s.Flags.Addr)
		err := http.ListenAndServeTLS(s.Flags.Addr, s.Flags.CertificateFile, s.Flags.CertificateKey, nil)
		log.Println("[err]", err.Error())
	} else {
		log.Println("[info] Start listening on", s.Flags.Addr)
		err := http.ListenAndServe(s.Flags.Addr, nil)
		log.Println("[err]", err.Error())
	}
}

// Reads the stored metadata.
func (s *Server) readMetadata() {
	file, err := os.Open(s.Flags.OutputDirectory + "/metadata.json")
	create := false
	if err != nil {
		create = true
	}

	if create {
		// Create the file
		log.Println("[info] Creating metadata.json")
		s.writeMetadata()
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

func (s *Server) writeMetadata() {
	file, err := os.Create(s.Flags.OutputDirectory + "/metadata.json")
	if err != nil {
		log.Println("[err] Can't write in the output directory:", s.Flags.OutputDirectory)
		log.Println(err)
		os.Exit(1)
	}
	data, _ := json.Marshal(s.Metadata)
	file.Write(data)
	file.Close()

	log.Printf("[info] %d metadatas written.\n", len(s.Metadata.Data))
}

// Prepares the route
func (s *Server) prepareRouter() *mux.Router {
	r := mux.NewRouter()

	sendHandler := &SendHandler{s}
	r.Handle(s.Flags.Route+"/1.0/send", sendHandler)

	deleteHandler := &DeleteHandler{s}
	r.Handle(s.Flags.Route+"/{file}/{key}", deleteHandler)

	sh := &ServingHandler{s}
	r.Handle(s.Flags.Route+"/{file}", sh) // Serving route.

	return r
}

// Expire expires a file : delete it from the metadata
// and from the FS.
func (s *Server) Expire(m Metadata) error {
	delete(s.Metadata.Data, m.Filename)
	return os.Remove(s.Flags.OutputDirectory + "/" + m.Filename)
}
