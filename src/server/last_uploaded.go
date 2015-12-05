// Route giving the last uploaded files.
// Copyright © 2015 - Rémy MATHIEU
package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const (
	MAX_LAST_UPLOADED = 20
)

type LastUploadedHandler struct {
	Server *Server // pointer to the started server
}

// Json returned to the client
type LastUploadedResponse struct {
	Name         string    `json:"name"`
	Original     string    `json:"original"`
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

	lastUploadedResp := make([]LastUploadedResponse, 0)
	lastUploaded, err := l.Server.GetLastUploaded()
	if err != nil {
		log.Println("[err] can't retrieve the last uploaded ids:", err.Error())
		w.WriteHeader(500)
		return
	}

	for _, id := range lastUploaded {
		metadata, err := l.Server.GetEntry(id)
		if err != nil {
			log.Println("[err] Error while retrieving an entry in LastUploaded handler:", err.Error())
			continue
		}

		if metadata == nil {
			continue
		}

		lastUploadedResp = append(lastUploadedResp, LastUploadedResponse{
			Name:         metadata.Filename,
			Original:     metadata.Original,
			DeleteKey:    metadata.DeleteKey,
			CreationTime: metadata.CreationTime,
		})
	}

	bytes, err := json.Marshal(lastUploadedResp)
	if err != nil {
		log.Println("[err] Can't marshal the list of last uploaded:", err.Error())
		w.WriteHeader(500)
	}

	w.Write(bytes)
}
