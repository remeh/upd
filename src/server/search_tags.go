// Route to search by tags
// Copyright © 2015 - Rémy MATHIEU
package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

type SearchTagsHandler struct {
	Server *Server // pointer to the started server
}

// Json returned to the client
type SearchTagsResponse struct {
	Results []SearchTagsEntryResponse `json:"results"`
}

type SearchTagsEntryResponse struct {
	Filename       string    `json:"filename"`        // name attributed by upd
	Original       string    `json:"original"`        // original name of the file
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
	for _, v := range l.Server.Metadata.Data {
		if stringArrayContainsOne(v.Tags, tags) {
			entry := SearchTagsEntryResponse{
				Filename:     v.Filename,
				Original:     v.Original,
				CreationTime: v.CreationTime,
				Tags:         v.Tags,
			}
			response.Results = append(response.Results, entry)
		}
	}

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
