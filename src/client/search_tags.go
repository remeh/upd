// Client - Searching tags.
// Copyright © 2015 - Rémy MATHIEU

package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"server"
)

const (
	ROUTE_SEARCH_TAGS = "/1.0/search_tags"
)

func (c *Client) SearchTags(tags []string) {
	// create the request
	client := c.createHttpClient()

	uri := c.Flags.ServerUrl + ROUTE_SEARCH_TAGS

	params := make(map[string]string)
	uri = c.buildParams(uri, params, tags)

	req, err := http.NewRequest("GET", uri, nil)

	// adds the secret key if any
	if len(c.Flags.SecretKey) > 0 {
		req.Header.Set(server.SECRET_KEY_HEADER, c.Flags.SecretKey)
	}

	// execute
	resp, err := client.Do(req)
	if err != nil {
		log.Println("[err] Unable to execute the request to search by tags:", err)
		os.Exit(1)
	}

	if resp.StatusCode != 200 {
		log.Printf("[err] Received a %d while searching by tags: %s", resp.StatusCode, tags)
		os.Exit(1)
	}

	// read the name given by the server
	defer resp.Body.Close()
	readBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[err] Unable to read the body returned by the server:", err)
		os.Exit(1)
	}

	var entries server.SearchTagsResponse
	err = json.Unmarshal(readBody, &entries)
	if err != nil {
		log.Println("[err] Unable to read the response:", err)
		os.Exit(1)
	}

	for i := range entries.Results {
		entry := entries.Results[i]
		fmt.Printf("Original name:\n  %s\n", entry.Original)
		fmt.Printf("Filename:\n  %s\n", entry.Filename)
		fmt.Printf("Creation time:\n  %s\n", entry.CreationTime)
		fmt.Printf("Tags:\n  %s\n", entry.Tags)
		fmt.Println("--")
	}
}
