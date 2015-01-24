// Route to server the files.
// Copyright © 2015 - Rémy MATHIEU

package server

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type ServingHandler struct {
	Server *Server // pointer to the started server
}

func (s *ServingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Parse the route parameters
	vars := mux.Vars(r)

	id := vars["file"]

	// Some check on the file id
	if len(id) == 0 {
		w.WriteHeader(404)
		return
	}

	// Existing file ?
	entry := s.Server.Metadata.Data[id]
	if entry.Filename == "" {
		w.WriteHeader(404)
		return
	}

	// Existing, serve the file !

	// read it
	file, err := os.Open(s.Server.Flags.OutputDirectory + "/" + entry.Filename)
	if err != nil {
		w.WriteHeader(500)
		log.Println("[err] While requesting:", entry.Filename)
		log.Println(err)
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(500)
		log.Println("[err] While reading:", entry.Filename)
		log.Println(err)
	}

	contentType := http.DetectContentType(data)
	w.Header().Set("Content-Type", contentType)
	w.Write(data)
}
