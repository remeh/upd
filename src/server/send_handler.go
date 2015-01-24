package server

import (
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type SendHandler struct {
	Server *Server // pointer to the started server
}

const (
	MAX_MEMORY = 1024 * 1024
	DICTIONARY = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func (s *SendHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reader, _, err := r.FormFile("data")

	if err != nil {
		w.WriteHeader(500)
		log.Println("[err] Error while receiving data (FormFile).")
		log.Println(err)
		return
	}

	// read the data
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		w.WriteHeader(500)
		log.Println("[err] Error while receiving data (ReadAll).")
		log.Println(err)
		return
	}

	// write the file in the directory
	name := s.randomString(8)
	s.writeFile(name, data)

	// add to metadata
	s.addMetadata(name)
	s.Server.writeMetadata() // TODO do it regularly instead of one big time.

	w.Write([]byte(name))
}

// TODO less ugly
func (s *SendHandler) randomString(size int) string {
	result := ""

	for i := 0; i < size; i++ {
		result += string(DICTIONARY[rand.Int31n(int32(len(DICTIONARY)))])
	}

	return result
}

func (s *SendHandler) addMetadata(name string) {
	metadata := Metadata{
		Filename:     name,
		CreationTime: time.Now(),
	}
	s.Server.Metadata.Data[name] = metadata
}

func (s *SendHandler) writeFile(filename string, data []byte) error {
	file, err := os.Create(s.Server.Flags.OutputDirectory + "/" + filename)
	if err != nil {
		log.Println("[err] Can't create the file to write: ", filename)
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		log.Println("[err] Can't write the file to write: ", filename)
		return err
	}

	err = file.Close()
	if err != nil {
		log.Println("[err] Can't close the file to write: ", filename)
		return err
	}

	return nil
}
