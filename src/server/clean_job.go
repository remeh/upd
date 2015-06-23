// Job launched and regularly executed to clean
// the expired files.
// Copyright © 2015 - Rémy MATHIEU

package server

import (
	"encoding/json"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

type CleanJob struct {
	server *Server
}

// Run deals with cleaning the expired files by
// checking their TTL.
func (j CleanJob) Run() {
	somethingChanged := false

	j.server.Database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Metadata"))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			// unmarshal
			var entry Metadata
			err := json.Unmarshal(v, &entry)
			if err != nil {
				log.Println("[err] Can't read a metadata:", err.Error())
				continue
			}

			if !entry.ExpirationTime.IsZero() && entry.ExpirationTime.Before(time.Now()) {
				// No longer alive!
				err := j.server.Expire(entry)
				somethingChanged = true
				if err != nil {
					log.Println("[warn] While deleting file:", entry.Filename)
					log.Println(err)
				} else {
					log.Println("[info] Deleted due to TTL:", entry.Filename)
				}
			}
		}

		return nil
	})
}
