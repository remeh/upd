// Instance of the server.
// Copyright © 2015 - Rémy MATHIEU

package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Server struct {
	Flags Flags               // Configuration
	Data  map[string]Metadata // Link to the read metadata
}

func NewServer(flags Flags) *Server {
	return &Server{
		Flags: flags,
		Data:  make(map[string]Metadata),
	}
}

// Starts the listening daemon.
func (s *Server) Start() {
	router := s.prepareRouter()

	// TODO check that we can write in the outputDirectory

	// Setup the router on the net/http stack
	http.Handle("/", router)

	// Read the existing metadata.
	s.readMetadata()

	// Listen
	// TODO TLS support
	http.ListenAndServe(s.Flags.Addr, nil)
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

		file, err = os.Create(s.Flags.OutputDirectory + "/metadata.json")
		if err != nil {
			log.Println("[err] Can't write in the output directory:", s.Flags.OutputDirectory)
			log.Println(err)
			os.Exit(1)
		}

		data, _ := json.Marshal(s.Data)
		file.Write(data)
		file.Close()
	} else {
		// Read the file
		log.Println("[info] Reading metadata.json")

		readData, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println("[err] The existing metadata.json seems corrupted. Exiting.")
			log.Println(err)
			os.Exit(1)
		}

		var data map[string]Metadata
		json.Unmarshal(readData, &data)
		s.Data = data

		log.Printf("[info] %d metadata read.\n", len(s.Data))
	}
}

// Prepares the route
func (s *Server) prepareRouter() *mux.Router {
	r := mux.NewRouter()

	sh := &ServingHandler{s}
	r.Handle(s.Flags.Route+"/{file}", sh) // Serving route.
	return r
}
