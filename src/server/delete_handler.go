// Route to delete a file from the server
// Copyright © 2015 - Rémy MATHIEU

package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type DeleteHandler struct {
	Server *Server // pointer to the started server
}

func (s *DeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Parse the route parameters
	vars := mux.Vars(r)

	id := vars["file"] // file id
	key := vars["key"] // delete key

	if len(id) == 0 || len(key) == 0 {
		w.WriteHeader(400)
		return
	}

	// Existing file ?
	entry, err := s.Server.GetEntry(id)
	if err != nil {
		log.Println("Can't use the database:", err.Error())
		w.WriteHeader(500)
		return
	}
	if entry == nil {
		w.WriteHeader(404)
		return
	}

	// checks that the key is correct
	if entry.DeleteKey != key {
		w.WriteHeader(403)
		return
	}

	// deletes the file
	err = s.Server.Expire(*entry)
	if err != nil {
		log.Println("[err] While deleting the entry:", entry.Filename)
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	// we must remove the line from the LastUploaded
	lastUploaded, err := s.Server.GetLastUploaded()
	if err != nil {
		log.Println("[err] Can't retrieve the last uploaded when deleting.")
		w.WriteHeader(500)
		return
	}

	lastUploaded = s.removeFromLastUploaded(lastUploaded, id)

	s.Server.SetLastUploaded(lastUploaded)

	w.WriteHeader(200)
	w.Write([]byte("File deleted."))
}

func (s *DeleteHandler) removeFromLastUploaded(lastUploaded []string, id string) []string {
	result := make([]string, 0)

	for _, lu := range lastUploaded {
		if lu == id {
			continue
		}
		result = append(result, lu)
	}

	return result
}
