// Client executable to send file to the upd daemon.
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

	"server"
)

const (
	ROUTE_SEND = "/1.0/send"
)

type Client struct {
	Flags Flags
}

func NewClient(flags Flags) *Client {
	return &Client{Flags: flags}
}

// Send sends the given file to the upd server.
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
	client := c.createClient()

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
		fmt.Println("Available until:", sendResponse.DeletionTime)
	}
	fmt.Println("--")

	return nil
}

// buildParams adds the GET parameters to the given uri.
func (c *Client) buildParams(uri string, params map[string]string, tags []string) string {
	if len(params) == 0 && len(tags) == 0 {
		return uri
	}

	atLeastOne := false

	ret := uri
	ret += "?"

	for k, v := range params {
		if len(v) > 0 {
			ret = fmt.Sprintf("%s%s=%s&", ret, k, url.QueryEscape(v))
			atLeastOne = true
		}
	}

	// remove last & if no tags
	if len(tags) == 0 {
		ret = ret[0 : len(ret)-1]
	}

	pTags := "tags="
	for i := range tags {
		pTags = pTags + url.QueryEscape(tags[i])
		if i < len(tags)-1 {
			pTags = pTags + url.QueryEscape(",")
		}
		atLeastOne = true
	}

	ret = ret + pTags

	// there were parameters but they're all empty
	if !atLeastOne {
		return uri
	}

	println(ret)
	return ret
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
