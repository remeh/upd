// Client executable to send file to the upd daemon.
// Copyright © 2015 - Rémy MATHIEU

package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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

func (c *Client) createHttpClient() *http.Client {
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
	} else {
		// add the tags.
		pTags := "tags="
		for i := range tags {
			pTags = pTags + url.QueryEscape(tags[i])
			if i < len(tags)-1 {
				pTags = pTags + url.QueryEscape(",")
			}
			atLeastOne = true
		}
		ret = ret + pTags
	}

	// there were parameters but they're all empty
	if !atLeastOne {
		return uri
	}

	return ret
}
