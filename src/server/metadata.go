// Saving information on the hosted files.
// Copyright © 2015 - Rémy MATHIEU
package server

import (
	"time"
)

type Metadata struct {
	Original       string    `json:"original"`        // original name of the file.
	Filename       string    `json:"filename"`        // name of the file on the FS
	Tags           []string  `json:"tags"`            // tags attached to the uploaded file
	TTL            string    `json:"ttl"`             // time.Duration representing the lifetime of the file.
	ExpirationTime time.Time `json:"expiration_time"` // at which time this file should expire.
	DeleteKey      string    `json:"delete_key"`      // The key to delete this file.
	CreationTime   time.Time `json:"creation_time"`
}
