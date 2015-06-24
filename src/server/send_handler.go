// Route receiving the data when a file is uploaded.
// Copyright © 2015 - Rémy MATHIEU
package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

type SendHandler struct {
	Server *Server // pointer to the started server
}

const (
	SECRET_KEY_HEADER = "X-upd-key"
)

// Json returned to the client
type SendResponse struct {
	Name           string    `json:"name"`
	DeleteKey      string    `json:"delete_key"`
	ExpirationTime time.Time `json:"expiration_time"`
}

const (
	MAX_MEMORY = 1024 * 1024
	DICTIONARY = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func (s *SendHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// checks the secret key
	key := r.Header.Get(SECRET_KEY_HEADER)
	if s.Server.Config.SecretKey != "" && key != s.Server.Config.SecretKey {
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
	var name string
	var original string

	// name
	if len(r.Form["name"]) == 0 {
		w.WriteHeader(400)
		return
	} else {
		original = filepath.Base(r.Form["name"][0])
	}

	for {
		name = s.randomString(8)
		// test existence
		entry, err := s.Server.GetEntry(name)
		if err != nil {
			log.Println("[err] While reading the database:", err.Error())
			w.WriteHeader(500)
			return
		}
		if entry.Filename == "" {
			break
		}
	}

	var expirationTime time.Time
	now := time.Now()

	// reads the TTL
	var ttl string
	if len(r.Form["ttl"]) > 0 {
		ttl = r.Form["ttl"][0]
		// check that the value is a correct duration
		_, err := time.ParseDuration(ttl)
		if err != nil {
			println(err.Error())
			w.WriteHeader(400)
			return
		}

		// compute the expiration time
		expirationTime = s.Server.computeEndOfLife(ttl, now)
	}

	// reads the tags
	tags := make([]string, 0)
	if len(r.Form["tags"]) > 0 {
		tags = strings.Split(r.Form["tags"][0], ",")
	}

	// writes the data on the storage
	if err := s.Server.WriteFile(name, data); err != nil {
		log.Println(err)
	}

	// add to metadata
	deleteKey := s.randomString(16)
	s.addMetadata(name, original, tags, expirationTime, ttl, deleteKey, now)

	// encode the response json
	response := SendResponse{
		Name:           name,
		DeleteKey:      deleteKey,
		ExpirationTime: expirationTime,
	}

	resp, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
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
func (s *SendHandler) addMetadata(name string, original string, tags []string, expirationTime time.Time, ttl string, key string, now time.Time) {
	metadata := Metadata{
		Filename:       name,
		Original:       original,
		Tags:           tags,
		TTL:            ttl,
		ExpirationTime: expirationTime,
		DeleteKey:      key,
		CreationTime:   now,
	}

	// marshal the object
	data, err := json.Marshal(metadata)
	if err != nil {
		log.Println("[err] Can't marshal an object to store it")
		log.Println(err)
		return
	}

	// store into BoltDB
	err = s.Server.Database.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Metadata"))
		return bucket.Put([]byte(name), data)
	})

	if err != nil {
		log.Println("[err] Can't store")
		log.Println(string(data))
		log.Printf("[err] Reason: %s\n", err.Error())
		return
	}

	// store the LastUploaded infos
	lastUploaded, err := s.Server.GetLastUploaded()
	if err != nil {
		log.Println("[err] Can't read the last uploaded in send handler:", err.Error())
		return
	}

	lastUploaded = append([]string{name}, lastUploaded[:]...)
	if len(lastUploaded) > MAX_LAST_UPLOADED {
		lastUploaded = lastUploaded[:len(lastUploaded)-1]
	}

	s.Server.SetLastUploaded(lastUploaded)

	if err != nil {
		log.Printf("[err] Can't store the LastUploaded infos for %s, reason: %s\n", name, err.Error())
		return
	}
}
