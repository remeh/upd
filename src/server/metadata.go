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
	Filename     string    `json:"filename"`      // name of the file on the FS
	RealFilename string    `json:"real_filename"` // Origin filename
	CreationTime time.Time `json:"creation_time"`
}
