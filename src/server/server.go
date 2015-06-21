// Instance of the server.
// Copyright © 2015 - Rémy MATHIEU

package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
)

type Server struct {
	Config   Config   // Configuration
	Database *bolt.DB // opened bolt db
	Storage  string   // Storage used with this metadata file.
}

func NewServer(config Config) *Server {
	// init the random
	rand.Seed(time.Now().Unix())

	return &Server{
		Config:  config,
		Storage: config.Storage,
	}
}

// Starts the listening daemon.
func (s *Server) Start() {
	router := s.prepareRouter()

	// Setup the router on the net/http stack
	http.Handle("/", router)

	// Open the database
	s.openBoltDatabase()

	s.convertMetadata()

	go s.StartCleanJob()

	// Listen
	if len(s.Config.CertificateFile) != 0 && len(s.Config.CertificateKey) != 0 {
		log.Println("[info] Start secure listening on", s.Config.Addr)
		err := http.ListenAndServeTLS(s.Config.Addr, s.Config.CertificateFile, s.Config.CertificateKey, nil)
		log.Println("[err]", err.Error())
	} else {
		log.Println("[info] Start listening on", s.Config.Addr)
		err := http.ListenAndServe(s.Config.Addr, nil)
		log.Println("[err]", err.Error())
	}
}

// Starts the Clean Job
func (s *Server) StartCleanJob() {
	timer := time.NewTicker(60 * time.Second)
	for _ = range timer.C {
		job := CleanJob{s}
		job.Run()
	}
}

// readMetadata reads an old .json metadata file.
// Deprecataed.
func (s *Server) convertMetadata() {
	file, err := os.Open(s.Config.RuntimeDir + "/metadata.json")
	if err != nil {
		log.Printf("[err] Can't find %s", s.Config.RuntimeDir+"/metadata.json")
		os.Exit(1)
	}

	defer file.Close()

	readData, err := ioutil.ReadAll(file)
	var data Metadatas
	json.Unmarshal(readData, &data)

	// save each entry in the boltdb database.
	sendHandler := SendHandler{Server: s}

	log.Println("[info] Start of .json -> .db processing.")

	for _, m := range data.Data {
		log.Printf("[info] %s processed.", m.Original)
		sendHandler.addMetadata(m.Filename, m.Original, m.Tags, m.ExpirationTime, m.TTL, m.DeleteKey, m.CreationTime)
	}

	log.Println("[info] End of .json -> .db processing.")
}

// openBoltDatabase opens an existing database.
func (s *Server) openBoltDatabase() {
	db, err := bolt.Open(s.Config.RuntimeDir+"/metadata.db", 0600, nil)
	if err != nil {
		log.Println("[err] Can't open the metadata.db file in :", s.Config.RuntimeDir)
		log.Println(err)
		os.Exit(1)
	}

	log.Printf("[info] %s opened.", s.Config.RuntimeDir+"/metadata.db")

	s.Database = db

	// creates the bucket if needed
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Metadata"))
		if err != nil {
			log.Println("Can't create the bucket 'Metadata'")
			log.Println(err)
		}
		_, err = tx.CreateBucketIfNotExists([]byte("LastUploaded"))
		if err != nil {
			log.Println("Can't create the bucket 'LastUploaded'")
			log.Println(err)
		}
		_, err = tx.CreateBucketIfNotExists([]byte("Config"))
		if err != nil {
			log.Println("Can't create the bucket 'LastUploaded'")
			log.Println(err)
		}
		return err
	})

	// test that the storage is still the same
	var mustSave bool
	s.Database.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Config"))
		v := bucket.Get([]byte("storage"))
		if v == nil {
			mustSave = true
			return nil
		}

		if string(v) != s.Config.Storage {
			log.Printf("The database use the storage %s, can't start with the storage %s\n", string(v), s.Config.Storage)
			os.Exit(1)
		}
		return nil
	})

	// save the storage
	if mustSave {
		s.Database.Update(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte("Config"))
			bucket.Put([]byte("storage"), []byte(s.Config.Storage))
			return nil
		})
	}
}

func (s *Server) deleteMetadata(name string) error {
	// remove from last uploaded
	entry, err := s.GetEntry(name)
	if err != nil {
		return err
	}
	err = s.Database.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("LastUploaded"))
		return bucket.Delete([]byte(entry.CreationTime.String()))
	})

	err = s.Database.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Metadata"))
		return bucket.Delete([]byte(name))
	})
	if err != nil {
		log.Println("Can't delete some metadata from the database:")
		log.Println(err)
	}

	return err
}

// getEntry looks in the Bolt DB whether this entry exists and returns it
// if found, otherwise, nil is returned.
func (s *Server) GetEntry(id string) (*Metadata, error) {
	var metadata Metadata
	err := s.Database.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Metadata"))
		v := bucket.Get([]byte(id))
		if v == nil {
			return nil
		}

		// unmarshal the bytes
		err := json.Unmarshal(v, &metadata)
		if err != nil {
			return err
		}

		return nil
	})

	return &metadata, err
}

// Prepares the route
func (s *Server) prepareRouter() *mux.Router {
	r := mux.NewRouter()

	sendHandler := &SendHandler{s}
	r.Handle(s.Config.Route+"/1.0/send", sendHandler)

	lastUploadeHandler := &LastUploadedHandler{s}
	r.Handle(s.Config.Route+"/1.0/list", lastUploadeHandler)

	searchTagsHandler := &SearchTagsHandler{s}
	r.Handle(s.Config.Route+"/1.0/search_tags", searchTagsHandler)

	deleteHandler := &DeleteHandler{s}
	r.Handle(s.Config.Route+"/{file}/{key}", deleteHandler)

	sh := &ServingHandler{s}
	r.Handle(s.Config.Route+"/{file}", sh) // Serving route.

	return r
}
