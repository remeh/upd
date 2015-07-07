// Client executable to upload data on a upd daemon.
// Copyright © 2015 - Rémy MATHIEU

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/vrischmann/userdir"

	"client"
)

type config struct {
	ServerURL string `toml:"server_url"`
	SecretKey string `toml:"secret_key"`
}

var conf config
var flags client.Flags

const (
	defaultServerURL = "http://localhost:9000/upd"
)

func init() {
	flag.StringVar(&(flags.CA), "ca", "none", "For HTTPS support: none / filename of an accepted CA / unsafe (doesn't check the CA)")
	flag.StringVar(&(flags.ServerUrl), "url", defaultServerURL, "The server to contact")
	flag.StringVar(&(flags.SecretKey), "key", "", "A shared secret key to identify the client.")
	flag.StringVar(&(flags.TTL), "ttl", "", `TTL after which the file expires, ex: 30m. Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h"`)
	flag.StringVar(&(flags.SearchTags), "search-tags", "", "Search by tags. If many, must be separated by a comma, an 'or' operator is used. Ex: \"may,screenshot\".")
	flag.Var(&flags.Tags, "tags", "Tags to attach to the file, separated by a comma. Ex: \"screenshot,may\"")
}

func parseConfig() error {
	path := filepath.Join(userdir.GetConfigHome(), "upd", "client.conf")

	fi, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if fi.IsDir() {
		return fmt.Errorf("%s is a directory, not a file", path)
	}

	// If the file does not exist quit as it's not an error
	if os.IsNotExist(err) {
		return nil
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	_, err = toml.Decode(string(data), &conf)
	return err
}

func parseFlags() error {
	flag.Parse()

	// remove / on url if necessary
	if flags.ServerUrl[len(flags.ServerUrl)-1] == '/' {
		flags.ServerUrl = flags.ServerUrl[:len(flags.ServerUrl)-1]
	}

	// checks that the given ttl is correct
	if flags.TTL != "" {
		_, err := time.ParseDuration(flags.TTL)
		return err
	}

	return nil
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

func replaceDefaultFlagsWithConfig() {
	if flags.ServerUrl == defaultServerURL {
		flags.ServerUrl = conf.ServerURL
	}
	if flags.SecretKey == "" {
		flags.SecretKey = conf.SecretKey
	}
}

func main() {
	if err := parseFlags(); err != nil {
		fmt.Println(`Wrong duration format, it should be such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h"`)
		os.Exit(1)
	}

	if err := parseConfig(); err != nil {
		fmt.Println("Unable to read configuration file, error was %s", err)
		os.Exit(1)
	}

	replaceDefaultFlagsWithConfig()

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
