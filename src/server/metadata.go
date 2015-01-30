// Saving information on the hosted files.
// Copyright © 2015 - Rémy MATHIEU
package server

import (
	"time"
)

type Metadatas struct {
	CreationTime time.Time           `json:"creation_time"`
	Data         map[string]Metadata `json:"metadatas"`
	LastUploaded []string            `json:"last_uploaded"` // stores the 20 last updated files id.
}

type Metadata struct {
	Original     string       `json:"original"`      // original name of the file.
	Filename     string       `json:"filename"`      // name of the file on the FS
	TTL          string       `json:"ttl"`           // time.Duration representing the lifetime of the file.
	DeleteKey    string       `json:"delete_key"`    // The key to delete this file.
	BackendInfos BackendInfos `json:"backend_infos"` // The used backend
	CreationTime time.Time    `json:"creation_time"`
}

// When wrote on a backend (S3, GCS), we'll need some more info.
type BackendInfos struct {
	Type   string // possible values : s3
	Bucket string // used by : s3
}
