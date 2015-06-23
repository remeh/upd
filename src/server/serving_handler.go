// Route to server the files.
// Copyright © 2015 - Rémy MATHIEU

package server

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/nfnt/resize"
)

const (
	HEADER_ORIGINAL_FILENAME = "X-Upd-Orig-Filename"
)

type ServingHandler struct {
	Server *Server // pointer to the started server
}

func (s *ServingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Parse the route parameters
	vars := mux.Vars(r)

	id := vars["file"]

	// Some check on the file id
	if len(id) == 0 {
		w.WriteHeader(404)
		return
	}

	// Look for the file in BoltDB
	entry, err := s.Server.GetEntry(id)
	if err != nil {
		log.Println("[err] Error while retrieving an entry:", err.Error())
		w.WriteHeader(500)
		return
	}

	// Existing file ?
	if entry.Filename == "" {
		w.WriteHeader(404)
		return
	}

	// Existing, serve the file !

	// but first, check that it hasn't expired
	if entry.TTL != "" {
		duration, _ := time.ParseDuration(entry.TTL)
		now := time.Now()
		fileEndlife := entry.CreationTime.Add(duration)
		if fileEndlife.Before(now) {
			// No longer alive!
			err := s.Server.Expire(*entry)
			if err != nil {
				log.Println("[warn] While deleting file:", entry.Filename)
				log.Println(err)
			} else {
				log.Println("[info] Deleted due to TTL:", entry.Filename)
			}

			w.WriteHeader(404)
			return
		}
	}

	// read it
	data, err := s.Server.ReadFile(entry.Filename)

	if err != nil {
		log.Println("[err] Can't read the file from the storage.")
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	// detect the content-type
	contentType := http.DetectContentType(data)

	// we'll see whether or not we want to generate a thumbnail
	r.ParseForm()
	width := r.Form.Get("w")
	height := r.Form.Get("h")
	if len(width) != 0 && len(height) != 0 {
		iwidth, err := strconv.Atoi(width)
		if err != nil {
			w.WriteHeader(400)
			return
		}
		iheight, err := strconv.Atoi(height)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		data = s.Resize(id, contentType, data, uint(iwidth), uint(iheight))
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set(HEADER_ORIGINAL_FILENAME, entry.Original)
	w.Header().Set("Content-Disposition", "inline; filename*=UTF-8''"+url.QueryEscape(entry.Original))
	w.Write(data)
}

func (s *ServingHandler) Resize(id string, contentType string, data []byte, width uint, height uint) []byte {
	if len(data) == 0 {
		return data
	}
	if contentType != "image/png" && contentType != "image/jpeg" {
		return data
	}

	var err error
	var img image.Image
	var buffer *bytes.Buffer

	// create the Image instance
	if contentType == "image/png" {
		img, err = png.Decode(bytes.NewReader(data))
		if err != nil {
			log.Println("[err] Can't resize png image with id:", id)
			return data
		}
	} else if contentType == "image/jpeg" {
		img, err = jpeg.Decode(bytes.NewReader(data))
		if err != nil {
			log.Println("[err] Can't resize jpg image with id:", id)
			return data
		}
	}

	// resize the image
	img = resize.Resize(width, height, img, resize.Lanczos3)

	// write the data
	buffer = bytes.NewBuffer(nil)
	if contentType == "image/png" {
		err = png.Encode(buffer, img)
		if err != nil {
			log.Println("[err] Can't encode png image with id:", id)
			return data
		}
	} else if contentType == "image/jpeg" {
		err = jpeg.Encode(buffer, img, nil)
		if err != nil {
			log.Println("[err] Can't encode jpg image with id:", id)
			return data
		}
	}

	return buffer.Bytes()
}
