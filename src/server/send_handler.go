package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type SendHandler struct {
	Server *Server // pointer to the started server
}

const (
	SECRET_KEY_HEADER = "X-Cloudia-Key"
)

// Json returned to the client
type SendResponse struct {
	Name         string `json:"name"`
	DeleteKey    string `json:"delete_key"`
	DeletionTime string `json:"availaible_until"`
}

const (
	MAX_MEMORY = 1024 * 1024
	DICTIONARY = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func (s *SendHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// checks the secret key
	key := r.Header.Get(SECRET_KEY_HEADER)
	if s.Server.Flags.SecretKey != "" && key != s.Server.Flags.SecretKey {
		w.WriteHeader(403)
		return
	}

	// parse the form
	reader, _, err := r.FormFile("data")

	if err != nil {
		w.WriteHeader(500)
		log.Println("[err] Error while receiving data (FormFile).")
		log.Println(err)
		return
	}

	// read the data
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		w.WriteHeader(500)
		log.Println("[err] Error while receiving data (ReadAll).")
		log.Println(err)
		return
	}

	// write the file in the directory
	name := ""

	// name
	if len(r.Form["name"]) > 0 {
		name = filepath.Base(r.Form["name"][0])
		// Reserved name.
		if name == "metadata.json" {
			w.WriteHeader(400)
			w.Write([]byte("'metadata.json' : reserved name."))
			return
		}
	} else {
		for {
			name = s.randomString(8)
			if s.Server.Metadata.Data[name].Filename == "" {
				break
			}
		}
		s.writeFile(name, data)
	}

	// reads the TTL
	var ttl string
	if len(r.Form["ttl"]) > 0 {
		ttl = r.Form["ttl"][0]
		// check that the value is a correct duration
		_, err := time.ParseDuration(ttl)
		if err != nil {
			w.WriteHeader(400)
			return
		}
	}

	// add to metadata
	now := time.Now()
	deleteKey := s.randomString(16)
	s.addMetadata(name, ttl, deleteKey, now)
	s.Server.writeMetadata() // TODO do it regularly instead of here.

	// encode the response json
	response := SendResponse{
		Name:         name,
		DeleteKey:    deleteKey,
		DeletionTime: s.computeEndOfLife(ttl, now),
	}

	resp, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

// computeEndOfLife return as a string the end of life of the new file.
func (s *SendHandler) computeEndOfLife(ttl string, now time.Time) string {
	if len(ttl) == 0 {
		return "Forever."
	}
	duration, _ := time.ParseDuration(ttl) // no error possible 'cause already checked in the controller
	t := now.Add(duration)
	return fmt.Sprintf("%s\n", t)
}

// randomString generates a random valid URL string of the given size
func (s *SendHandler) randomString(size int) string {
	result := ""

	for i := 0; i < size; i++ {
		result += string(DICTIONARY[rand.Int31n(int32(len(DICTIONARY)))])
	}

	return result
}

// addMetadata adds the given entry to the Server metadata information.
func (s *SendHandler) addMetadata(name string, ttl string, key string, now time.Time) {
	metadata := Metadata{
		Filename:     name,
		TTL:          ttl,
		DeleteKey:    key,
		CreationTime: now,
	}
	s.Server.Metadata.Data[name] = metadata
}

// writeFile uses the Server flags to save the given data as a file
// on the FS.
func (s *SendHandler) writeFile(filename string, data []byte) error {
	file, err := os.Create(s.Server.Flags.OutputDirectory + "/" + filename)
	if err != nil {
		log.Println("[err] Can't create the file to write: ", filename)
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		log.Println("[err] Can't write the file to write: ", filename)
		return err
	}

	err = file.Close()
	if err != nil {
		log.Println("[err] Can't close the file to write: ", filename)
		return err
	}

	return nil
}
