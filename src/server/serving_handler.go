// Route to server the files.
// Copyright © 2015 - Rémy MATHIEU

package server

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type ServingHandler struct {
	Server *Server // pointer to the started server
}

func (s *ServingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Parse the route parameters
	vars := mux.Vars(r)

	id := vars["file"]

	// Some check on the file id
	if len(id) == 0 {
		w.WriteHeader(404)
		return
	}

	// Existing file ?
	entry := s.Server.Metadata.Data[id]
	if entry.Filename == "" {
		w.WriteHeader(404)
		return
	}

	// Existing, serve the file !

	// but first, check that it hasn't expired
	if entry.TTL != "" {
		duration, _ := time.ParseDuration(entry.TTL)
		now := time.Now()
		fileEndlife := entry.CreationTime.Add(duration)
		if fileEndlife.Before(now) {
			// No longer alive!
			err := s.Server.Expire(entry)
			if err != nil {
				log.Println("[warn] While deleting file:", entry.Filename)
				log.Println(err)
			} else {
				log.Println("[info] Deleted due to TTL:", entry.Filename)
				s.Server.writeMetadata()
			}

			w.WriteHeader(404)
			return
		}
	}

	// read it
	file, err := os.Open(s.Server.Flags.OutputDirectory + "/" + entry.Filename)
	if err != nil {
		w.WriteHeader(500)
		log.Println("[err] While requesting:", entry.Filename)
		log.Println(err)
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(500)
		log.Println("[err] While reading:", entry.Filename)
		log.Println(err)
	}

	contentType := http.DetectContentType(data)
	w.Header().Set("Content-Type", contentType)
	w.Write(data)
}
