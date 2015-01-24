// Server flags.
// Copyright © 2015 - Rémy MATHIEU

package server

// Flags for server configuration
type Flags struct {
	Addr            string // Address to listen to
	SecretKey       string // Secret between the client and the server
	OutputDirectory string // Where the server can write the files.
	Route           string // Route served by the webserver
	CertificateFile string // Filepath to an tls certificate
	CertificateKey  string // Filepath to the key part of a certificate
}
