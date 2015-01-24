// Route to server the files.
// Copyright © 2015 - Rémy MATHIEU

package server

import (
	"net/http"
)

type ServingHandler struct {
	Server *Server // pointer to the started server
}

func (s *ServingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
