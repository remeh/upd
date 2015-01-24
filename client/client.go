// Client executable to upload data on a clioud daemon.
// Copyright © 2015 - Rémy MATHIEU

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"client"
)

type Client struct {
}

func parseFlags() client.Flags {
	var flags client.Flags

	// Declare the flags
	flag.StringVar(&(flags.ServerUrl), "url", "http://localhost:9000/files", "The server to contact")
	flag.StringVar(&(flags.SecretKey), "key", "", "The secret key to identify the client.")

	// Read them
	flag.Parse()

	// remove / on url if necessary
	if flags.ServerUrl[len(flags.ServerUrl)-1] == '/' {
		flags.ServerUrl = flags.ServerUrl[:len(flags.ServerUrl)-1]
	}

	return flags
}

// sendFile uses the client to send the data to the clioud server.
func sendFile(wg *sync.WaitGroup, client *client.Client, filename string) {
	defer wg.Done()

	err := client.Send(filename)
	if err != nil {
		log.Println("[err] While sending:", filename)
		log.Println(err)
	}
}

func main() {
	flags := parseFlags()

	// Looks for the file to send
	// TODO directory
	if len(os.Args) == 1 {
		fmt.Printf("Usage: %s [flags] file1 file2\n", os.Args[0])
		flag.PrintDefaults()
	}

	c := client.NewClient(flags)

	var wg sync.WaitGroup
	// Send each file.
	for _, filename := range os.Args[1:] {
		wg.Add(1)
		go sendFile(&wg, c, filename)
	}

	wg.Wait() // Wait for all routine to stop
}
