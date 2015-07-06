// Client flags.
// Copyright © 2015 - Rémy MATHIEU

package client

import "strings"

type Tags []string

func (t *Tags) Set(s string) error {
	*t = append(*t, strings.TrimSpace(s))
	return nil
}

func (t Tags) String() string {
	return strings.Join(t, ",")
}

// Flags for client configuration
type Flags struct {
	ServerUrl  string // Address to send to
	SecretKey  string // Secret between the client and the server
	TTL        string // when a ttl is given for a file
	CA         string // Should we use HTTPS, and in which config "none", file to a CA or "unsafe"
	SearchTags string // if we wanna look for some files by tags

	Tags Tags // Array of tag to attach to the file
}
