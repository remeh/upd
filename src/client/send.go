// Client - Sending file.
// Copyright © 2015 - Rémy MATHIEU

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"server"
)

// sendData sends the data to the upd server.
func (c *Client) sendData(filename string, data []byte) error {
	// Prepare the multipart content
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("data", "file")
	if err != nil {
		log.Println("[err] Unable to prepare the multipart content (CreateFormFile)")
		return err
	}

	_, err = io.Copy(part, bytes.NewReader(data))
	if err != nil {
		log.Println("[err] Unable to prepare the multipart content (Copy)")
		return err
	}

	err = writer.Close()
	if err != nil {
		log.Println("[err] Unable to prepare the multipart content (Close)")
		return err
	}

	// create the request
	client := c.createHttpClient()

	uri := c.Flags.ServerUrl + ROUTE_SEND

	params := make(map[string]string)
	params["ttl"] = c.Flags.TTL
	params["name"] = filename

	uri = c.buildParams(uri, params, c.Flags.Tags)

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	if err != nil {
		log.Println("[err] Unable to create the request to send the file.")
		return err
	}

	// adds the secret key if any
	if len(c.Flags.SecretKey) > 0 {
		req.Header.Set(server.SECRET_KEY_HEADER, c.Flags.SecretKey)
	}

	// execute
	resp, err := client.Do(req)
	if err != nil {
		log.Println("[err] Unable to execut the request to send the file.")
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("[err] Received a %d while sending: %s", resp.StatusCode, filename)
	}

	// read the name given by the server
	defer resp.Body.Close()
	readBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[err] Unable to read the body returned by the server.")
		return err
	}

	// decodes the json
	var sendResponse server.SendResponse
	err = json.Unmarshal(readBody, &sendResponse)
	if err != nil {
		log.Println("[err] Unable to read the returned JSON.")
	}

	fmt.Println("For file :", filename)
	fmt.Println("URL:", c.Flags.ServerUrl+"/"+sendResponse.Name)
	fmt.Println("Delete URL:", c.Flags.ServerUrl+"/"+sendResponse.Name+"/"+sendResponse.DeleteKey)

	// compute until when it'll be available
	if sendResponse.DeletionTime.IsZero() {
		fmt.Println("Available forever.")
	} else {
		fmt.Println("Available until:", sendResponse.ExpirationTime)
	}
	fmt.Println("--")

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
