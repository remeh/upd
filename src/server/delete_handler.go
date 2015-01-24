// Route to delete a file from the server
// Copyright © 2015 - Rémy MATHIEU

package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type DeleteHandler struct {
	Server *Server // pointer to the started server
}

func (s *DeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Parse the route parameters
	vars := mux.Vars(r)

	id := vars["file"] // file id
	key := vars["key"] // delete key

	if len(id) == 0 || len(key) == 0 {
		w.WriteHeader(400)
		return
	}

	// Existing file ?
	entry := s.Server.Metadata.Data[id]
	if entry.Filename == "" {
		w.WriteHeader(404)
		return
	}

	// checks that the key is correct
	if entry.DeleteKey != key {
		w.WriteHeader(403)
		return
	}

	// deletes the file
	err := s.Server.Expire(entry)
	if err != nil {
		log.Println("[err] While deleting the entry:", entry.Filename)
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	// Re-save the metadata
	s.Server.writeMetadata()

	w.WriteHeader(200)
	w.Write([]byte("File deleted."))
}
