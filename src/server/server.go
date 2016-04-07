// Instance of the server.
// Copyright © 2015 - Rémy MATHIEU

package server

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
)

const (
	LAST_UPLOADED_KEY = "LastUploaded"
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
	for range timer.C {
		job := CleanJob{s}
		job.Run()
	}
}

// writeBoltMetadata stores the metadata in a BoltDB file.
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
		_, err = tx.CreateBucketIfNotExists([]byte("Runtime"))
		if err != nil {
			log.Println("Can't create the bucket 'Runtime'")
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
	err := s.Database.Update(func(tx *bolt.Tx) error {
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
	var metadata *Metadata
	err := s.Database.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Metadata"))
		v := bucket.Get([]byte(id))
		if v == nil {
			return nil
		}

		metadata = new(Metadata)

		// unmarshal the bytes
		err := json.Unmarshal(v, &metadata)
		if err != nil {
			return err
		}

		return nil
	})

	return metadata, err
}

// GetLastUploaded reads into BoltDB the array
// of last uploaded entries.
func (s *Server) GetLastUploaded() ([]string, error) {
	values := make([]string, 0)
	err := s.Database.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Runtime"))
		v := bucket.Get([]byte(LAST_UPLOADED_KEY))

		if v == nil {
			return nil
		}

		err := json.Unmarshal(v, &values)
		return err
	})

	return values, err
}

func (s *Server) SetLastUploaded(lastUploaded []string) error {
	return s.Database.Update(func(tx *bolt.Tx) error {
		d, err := json.Marshal(lastUploaded)
		if err != nil {
			return err
		}

		bucket := tx.Bucket([]byte("Runtime"))
		return bucket.Put([]byte(LAST_UPLOADED_KEY), d)
	})
}

// Prepares the route
func (s *Server) prepareRouter() *mux.Router {
	r := mux.NewRouter()

	println(s.Config.Route)
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
