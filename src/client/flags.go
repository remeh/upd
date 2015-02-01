// Client flags.
// Copyright © 2015 - Rémy MATHIEU

package client

// Flags for client configuration
type Flags struct {
	ServerUrl string   // Address to send to
	SecretKey string   // Secret between the client and the server
	TTL       string   // when a ttl is given for a file
	CA        string   // Should we use HTTPS, and in which config "none", file to a CA or "unsafe"
	Tags      []string // Array of tag to attach to the file
}
