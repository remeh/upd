package server

import "net/http"

// AuthCheckHandler simply tests whether the presented auth credentials are valid or not
// without doing any useless work
type AuthCheckHandler struct {
	Server *Server // pointer to the started server
}

func (a *AuthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !IsAuthValid(a.Server, r) {
		w.WriteHeader(403)
		w.Write([]byte("{auth_status:'invalid_credentials'}"))
		return
	}

	w.Write([]byte("{auth_status:'ok'}"))
}

// IsAuthValid returns whether the HTTP request contains the expected secret key, if the configuration
// requires one
func IsAuthValid(s *Server, r *http.Request) bool {
	key := r.Header.Get(SECRET_KEY_HEADER)
	return s.Config.SecretKey == "" || key == s.Config.SecretKey
}
