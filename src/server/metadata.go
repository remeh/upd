// Saving information on the hosted files.
// Copyright © 2015 - Rémy MATHIEU
package server

import (
	"time"
)

type Metadatas struct {
	CreationTime time.Time
	Metadata     []Metadata
}

type Metadata struct {
	Filename     string // name of the file on the FS
	RealFilename string // Origin filename
	CreationTime time.Time
}
