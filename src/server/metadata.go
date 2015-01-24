// Saving information on the hosted files.
// Copyright © 2015 - Rémy MATHIEU
package server

import (
	"time"
)

type Metadatas struct {
	CreationTime time.Time           `json:"creation_time"`
	Data         map[string]Metadata `json:"metadatas"`
}

type Metadata struct {
	Filename     string    `json:"filename"`   // name of the file on the FS
	TTL          string    `json:"ttl"`        // time.Duration representing the lifetime of the file.
	DeleteKey    string    `json:"delete_key"` // The key to delete this file.
	CreationTime time.Time `json:"creation_time"`
}
