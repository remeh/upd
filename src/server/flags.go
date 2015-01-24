// Server executable to receive/host files.
// Copyright © 2015 - Rémy MATHIEU

package server

// Flags for server configuration
type Flags struct {
	Addr            string // Address to listen to
	SecretKey       string // Secret between the client and the server
	OutputDirectory string // Where the server can write the files.
	Route           string // Route served by the webserver
}
