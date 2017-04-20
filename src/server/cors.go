package server

import "net/http"

// CorsHandler adds the required CORS headers, and forwards the request to the real handler
// (with the notable exception of OPTIONS requests, that it will eat)
type CorsHandler struct {
	h http.Handler
}

func (c *CorsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	headers.Set("Allow", "GET, POST, OPTIONS")
	headers.Set("Access-Control-Allow-Headers", "Content-Type, Accept, X-upd-key")
	headers.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	headers.Set("Access-Control-Allow-Origin", "*")

	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
	} else {
		c.h.ServeHTTP(w, r)
	}
}
