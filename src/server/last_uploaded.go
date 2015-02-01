// Route giving the last uploaded files.
// Copyright © 2015 - Rémy MATHIEU
package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type LastUploadedHandler struct {
	Server *Server // pointer to the started server
}

// Json returned to the client
type LastUploadedResponse struct {
	Name         string    `json:"name"`
	DeleteKey    string    `json:"delete_key"`
	CreationTime time.Time `json:"creation_time"`
}

func (l *LastUploadedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// checks the secret key
	key := r.Header.Get(SECRET_KEY_HEADER)
	if l.Server.Config.SecretKey != "" && key != l.Server.Config.SecretKey {
		w.WriteHeader(403)
		return
	}

	lastUploaded := make([]LastUploadedResponse, len(l.Server.Metadata.LastUploaded))
	for i, v := range l.Server.Metadata.LastUploaded {
		entry := l.Server.Metadata.Data[v]
		lastUploaded[i] = LastUploadedResponse{
			Name:         entry.Filename,
			DeleteKey:    entry.DeleteKey,
			CreationTime: entry.CreationTime,
		}
	}

	bytes, err := json.Marshal(lastUploaded)
	if err != nil {
		log.Println("[err] Can't marshal the list of last uploaded:", err.Error())
		w.WriteHeader(500)
	}

	w.Write(bytes)
}
