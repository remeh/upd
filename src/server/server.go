// Instance of the server.
// Copyright © 2015 - Rémy MATHIEU

package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	Flags Flags // Configuration
}

// Starts the listening daemon.
func (s *Server) Start() {
	router := s.prepareRouter()

	// Setup the router on the net/http stack
	http.Handle("/", router)

	// Listen
	// TODO TLS support
	http.ListenAndServe(s.Flags.Addr, nil)
}

// Prepares the route
func (s *Server) prepareRouter() *mux.Router {
	r := mux.NewRouter()

	sh := &ServingHandler{s}
	r.Handle(s.Flags.Route+"/{file}", sh) // Serving route.
	return r
}
