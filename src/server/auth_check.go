package server

import "net/http"

// AuthCheckHandler simply tests whether the presented auth credentials are valid or not
// without doing any useless work
type AuthCheckHandler struct {
	Server *Server // pointer to the started server
}

func (l *AuthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// checks the secret key
	key := r.Header.Get(SECRET_KEY_HEADER)
	if l.Server.Config.SecretKey != "" && key != l.Server.Config.SecretKey {
		w.WriteHeader(403)
		w.Write([]byte("{auth_status:'invalid_credentials'}"))
		return
	}

	w.Write([]byte("{auth_status:'ok'}"))
}
