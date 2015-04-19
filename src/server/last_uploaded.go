// Route giving the last uploaded files.
// Copyright © 2015 - Rémy MATHIEU
package server

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
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

	lastUploaded := make([]LastUploadedResponse, 0)
	l.Server.Database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("LastUploaded"))
		c := b.Cursor()

		min := []byte("1990-01-01T00:00:00Z")
		max := []byte(time.Now().String())

		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) >= 0; k, v = c.Next() {
			// unmarshal
			var metadata Metadata
			err := json.Unmarshal(v, &metadata)
			if err != nil {
				log.Println("[err] Can't read a metadata:", err.Error())
				continue
			}

			lastUploaded = append(lastUploaded, LastUploadedResponse{
				Name:         metadata.Filename,
				Original:     metadata.Original,
				DeleteKey:    metadata.DeleteKey,
				CreationTime: metadata.CreationTime,
			})
		}
		return nil
	})

	bytes, err := json.Marshal(lastUploaded)
	if err != nil {
		log.Println("[err] Can't marshal the list of last uploaded:", err.Error())
		w.WriteHeader(500)
	}

	w.Write(bytes)
}
