// Route to search by tags
// Copyright © 2015 - Rémy MATHIEU
package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

type SearchTagsHandler struct {
	Server *Server // pointer to the started server
}

// Json returned to the client
type SearchTagsResponse struct {
	Results []SearchTagsEntryResponse `json:"results"`
}

// actually contains everything in Metadata but eh,
// looks more clean to do so if oneee daaay...
type SearchTagsEntryResponse struct {
	Filename       string    `json:"filename"`        // name attributed by upd
	Original       string    `json:"original"`        // original name of the file
	DeleteKey      string    `json:"delete_key"`      // the delete key
	CreationTime   time.Time `json:"creation_time"`   // creation time of the given file
	ExpirationTime time.Time `json:"expiration_time"` // When this file expired
	Tags           []string  `json:"tags"`            // Tags attached to this file.
}

func (l *SearchTagsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// checks the secret key
	key := r.Header.Get(SECRET_KEY_HEADER)
	if l.Server.Config.SecretKey != "" && key != l.Server.Config.SecretKey {
		w.WriteHeader(403)
		return
	}

	r.ParseForm()
	if len(r.Form["tags"]) == 0 || len(r.Form["tags"][0]) == 0 {
		println("!")
		w.WriteHeader(400)
		return
	}

	tagParam := r.Form["tags"][0]
	tags := strings.Split(tagParam, ",")
	for i := range tags {
		tags[i] = strings.Trim(tags[i], " ")
	}

	// At the moment, without a 'database' system, we must
	// look in every entries to find if some have the tags.
	response := SearchTagsResponse{Results: make([]SearchTagsEntryResponse, 0)}

	l.Server.Database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Metadata"))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			// unmarshal
			var metadata Metadata
			err := json.Unmarshal(v, &metadata)
			if err != nil {
				log.Println("[err] Can't read a metadata:", err.Error())
				continue
			}

			if stringArrayContainsOne(metadata.Tags, tags) {
				entry := SearchTagsEntryResponse{
					Filename:       metadata.Filename,
					Original:       metadata.Original,
					CreationTime:   metadata.CreationTime,
					DeleteKey:      metadata.DeleteKey,
					ExpirationTime: metadata.ExpirationTime,
					Tags:           metadata.Tags,
				}
				response.Results = append(response.Results, entry)
			}
		}
		return nil
	})

	bytes, err := json.Marshal(response)
	if err != nil {
		log.Println("[err] Can't marshal the list of last uploaded:", err.Error())
		w.WriteHeader(500)
	}

	w.Write(bytes)
}

// stringArrayContains returns true if the array contains at least one of the values.
func stringArrayContainsOne(array []string, tags []string) bool {
	for i := range array {
		for j := range tags {
			if array[i] == tags[j] {
				return true
			}
		}
	}
	return false
}
