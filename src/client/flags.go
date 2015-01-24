// Client flags.
// Copyright © 2015 - Rémy MATHIEU

package client

// Flags for client configuration
type Flags struct {
	ServerUrl string // Address to send to
	SecretKey string // Secret between the client and the server
	Keepname  bool   // Whether or not we must keep the name
	TTL       string // when a ttl is given for a file
}
