// Client executable to send file to the clioud daemon.
// Copyright © 2015 - Rémy MATHIEU

package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Client struct {
	Flags Flags
}

func NewClient(flags Flags) *Client {
	return &Client{Flags: flags}
}

// Send sends the given file to the clioud server.
func (c *Client) Send(filename string) error {
	// first, we need to read the data
	data, err := c.readFile(filename)
	if err != nil {
		return err
	}

	// and now to send it the server
	return c.sendData(filename, data)
}

// sendData sends the data to the clioud server.
func (c *Client) sendData(filename string, data []byte) error {
	body := bytes.NewReader(data)

	contentType := http.DetectContentType(data)

	resp, err := http.Post(c.Flags.ServerUrl+"/send", contentType, body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Received a %d while sending: %s", resp.StatusCode, filename)
	}

	return nil
}

// readFile reads the content of the given file.
func (c *Client) readFile(filename string) ([]byte, error) {
	result := make([]byte, 0)

	// opening
	file, err := os.Open(filename)
	if err != nil {
		return result, err
	}

	// reading
	result, err = ioutil.ReadAll(file)
	if err != nil {
		return result, err
	}

	return result, nil
}
