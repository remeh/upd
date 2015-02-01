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

type SearchHandler struct {
	Server *Server // pointer to the started server
}

// Json returned to the client
type SearchResponse struct {
	Results []SearchEntryResponse `json:"results"`
}

type SearchEntryResponse struct {
	Filename     string    `json:"filename"`      // name attributed by upd
	Original     string    `json:"original"`      // original name of the file
	CreationTime time.Time `json:"creation_time"` // creation time of the given file
}

func (l *SearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// checks the secret key
	key := r.Header.Get(SECRET_KEY_HEADER)
	if l.Server.Config.SecretKey != "" && key != l.Server.Config.SecretKey {
		w.WriteHeader(403)
		return
	}

	r.ParseForm()
	if len(r.Form["tag"]) == 0 || len(r.Form["tag"][0]) == 0 {
		w.WriteHeader(400)
		return
	}

	tagParam := r.Form["tag"][0]
	tags := strings.Split(tagParam, ",")
	for i := range tags {
		tags[i] = strings.Trim(tags[i], " ")
	}

	// At the moment, without a 'database' system, we must
	// look in every entries to find if some have the tags.
	response := SearchResponse{Results: make([]SearchEntryResponse, 0)}
	for _, v := range l.Server.Metadata.Data {
		if stringArrayContainsOne(v.Tags, tags) {
			entry := SearchEntryResponse{
				Filename:     v.Filename,
				Original:     v.Original,
				CreationTime: v.CreationTime,
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
