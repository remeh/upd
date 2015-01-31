// Server config.
// Copyright © 2015 - Rémy MATHIEU

package server

// Server flags
type Flags struct {
	ConfigFile string // the file configuration to use for the server
}

const (
	FS_STORAGE = "fs"
	S3_STORAGE = "s3"
)

// Server configuration
type Config struct {
	Addr            string `toml:"listen_addr"`     // Address to listen to
	SecretKey       string `toml:"secret_key"`      // Secret between the client and the server
	RuntimeDir      string `toml:"run_dir"`         // Where the server can write the runtime files.
	Route           string `toml:"route"`           // Route served by the webserver
	CertificateFile string `toml:"certificate"`     // Filepath to an tls certificate
	CertificateKey  string `toml:"certificate_key"` // Filepath to the key part of a certificate

	Storage string `toml:"storage"` // possible values 'fs', 's3'

	FSConfig FSConfig `toml:"fsstorage"`
	S3Config S3Config `toml:"s3storage"`
}

type FSConfig struct {
	OutputDirectory string `toml:"output_dir"`
}

type S3Config struct {
	AccessKey    string `toml:"access_key"`
	AccessSecret string `toml:"access_secret"`
	Bucket       string `toml:"bucket"`
}
