// Server config.
// Copyright © 2015 - Rémy MATHIEU

package server

// Server flags
type Flags struct {
	ConfigFile string // the file configuration to use for the server
}

// Server configuration
type Config struct {
	Addr            string `toml:"listen_addr"`     // Address to listen to
	SecretKey       string `toml:"secret_key"`      // Secret between the client and the server
	OutputDirectory string `toml:"output_dir"`      // Where the server can write the files.
	Route           string `toml:"route"`           // Route served by the webserver
	CertificateFile string `toml:"certificate"`     // Filepath to an tls certificate
	CertificateKey  string `toml:"certificate_key"` // Filepath to the key part of a certificate
}
