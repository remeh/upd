// Route to server the files.
// Copyright © 2015 - Rémy MATHIEU

package server

import (
	"fmt"
	"net/http"

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
	if s.Server.Metadata.Data[id].Filename == "" {
		w.WriteHeader(404)
		return
	}

	// Existing, serve the file !

}
