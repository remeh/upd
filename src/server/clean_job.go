// Job launched and regularly executed to clean
// the expired files.
// Copyright © 2015 - Rémy MATHIEU

package server

import (
	"log"
	"time"
)

type CleanJob struct {
	server *Server
}

// Run deals with cleaning the expired files by
// checking their TTL.
func (j CleanJob) Run() {
	somethingChanged := false
	for _, entry := range j.server.Metadata.Data {
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
	if somethingChanged {
		j.server.writeMetadata(true)
	}
}
