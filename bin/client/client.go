// Client executable to upload data on a upd daemon.
// Copyright © 2015 - Rémy MATHIEU

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"client"
)

type Client struct {
}

func parseFlags() (client.Flags, error) {
	var flags client.Flags
	var tags string

	// Declare the flags
	flag.StringVar(&(flags.CA), "ca", "none", "For HTTPS support: none / filename of an accepted CA / unsafe (doesn't check the CA)")
	flag.StringVar(&(flags.ServerUrl), "url", "http://localhost:9000/upd", "The server to contact")
	flag.StringVar(&(flags.SecretKey), "key", "", "A shared secret key to identify the client.")
	flag.StringVar(&(flags.TTL), "ttl", "", `TTL after which the file expires, ex: 30m. Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h"`)
	flag.StringVar(&(flags.SearchTags), "search-tags", "", "Search by tags. If many, must be separated by a comma, an 'or' operator is used. Ex: \"may,screenshot\".")
	flag.StringVar(&tags, "tags", "", "Tags to attach to the file, separated by a comma. Ex: \"screenshot,may\"")

	// Read them
	flag.Parse()

	// remove / on url if necessary
	if flags.ServerUrl[len(flags.ServerUrl)-1] == '/' {
		flags.ServerUrl = flags.ServerUrl[:len(flags.ServerUrl)-1]
	}

	// checks that the given ttl is correct
	if flags.TTL != "" {
		_, err := time.ParseDuration(flags.TTL)
		if err != nil {
			return flags, err
		}
	}

	// split the tags in an array and trim them
	if len(tags) > 0 {
		flags.Tags = strings.Split(tags, ",")
		for i := range flags.Tags {
			flags.Tags[i] = strings.Trim(flags.Tags[i], " ")
		}
	}

	return flags, nil
}

// sendFile uses the client to send the data to the upd server.
func sendFile(wg *sync.WaitGroup, client *client.Client, filename string) {
	defer wg.Done()

	err := client.Send(filename)
	if err != nil {
		log.Println("[err] While sending:", filename)
		log.Println(err)
	}
}

func main() {
	flags, err := parseFlags()
	if err != nil {
		fmt.Println(`Wrong duration format, it should be such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h"`)
		os.Exit(1)
	}

	c := client.NewClient(flags)

	// Looks for tags to search
	if len(flags.SearchTags) > 0 {
		tags := strings.Split(flags.SearchTags, ",")
		for i := range tags {
			tags[i] = strings.Trim(tags[i], " ")
		}
		c.SearchTags(tags)
	} else {
		// Looks for the file to send
		// TODO directory
		if len(flag.Args()) < 1 {
			fmt.Printf("Usage: %s [flags] file1 file2\n", os.Args[0])
			flag.PrintDefaults()
		}

		var wg sync.WaitGroup
		// Send each file.
		for _, filename := range flag.Args() {
			wg.Add(1)
			go sendFile(&wg, c, filename)
		}

		wg.Wait() // Wait for all routine to stop
	}
}
