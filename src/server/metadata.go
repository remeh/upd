// Saving information on the hosted files.
// Copyright © 2015 - Rémy MATHIEU
package server

import (
	"time"
)

type Metadatas struct {
	CreationTime time.Time           `json:"creation_time"`
	Storage      string              `json:"storage"` // Storage used with this metadata file.
	Data         map[string]Metadata `json:"metadatas"`
	LastUploaded []string            `json:"last_uploaded"` // stores the 20 last updated files id.
}

type Metadata struct {
	Original     string    `json:"original"`   // original name of the file.
	Filename     string    `json:"filename"`   // name of the file on the FS
	Tags         []string  `json:"tags"`       // tags attached to the uploaded file
	TTL          string    `json:"ttl"`        // time.Duration representing the lifetime of the file.
	DeleteKey    string    `json:"delete_key"` // The key to delete this file.
	CreationTime time.Time `json:"creation_time"`
}
