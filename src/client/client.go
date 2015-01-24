// Client executable to send file to the clioud daemon.
// Copyright © 2015 - Rémy MATHIEU

package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"server"
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

	// and now to send it the servee
	return c.sendData(filename, data)
}

func (c *Client) createClient() *http.Client {
	if c.Flags.CA == "unsafe" {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		return &http.Client{Transport: tr}
	} else if len(c.Flags.CA) > 0 && c.Flags.CA != "none" {
		// reads the CA
		certs := x509.NewCertPool()
		pemData, err := ioutil.ReadFile(c.Flags.CA)
		if err != nil {
			log.Println("[err] Can't read the CA.")
		}
		certs.AppendCertsFromPEM(pemData)
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: certs},
		}
		return &http.Client{Transport: tr}
	}

	// No HTTPS support
	return &http.Client{}
}

// sendData sends the data to the clioud server.
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
	client := c.createClient()

	uri := c.Flags.ServerUrl + "/1.0/send"
	if len(c.Flags.TTL) > 0 {
		uri = uri + "?ttl=" + c.Flags.TTL
	}
	if c.Flags.Keepname {
		if len(c.Flags.TTL) != 0 {
			uri += "&"
		} else {
			uri += "?"
		}
		uri += "name=" + url.QueryEscape(filepath.Base(filename))
	}

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
	fmt.Println("Available until:", sendResponse.DeletionTime)
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
