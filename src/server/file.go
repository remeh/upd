// Methods dealing with files writing/reading
// from either the FS or a distant service (s3, ...)
// Copyright © 2015 - Rémy MATHIEU
package server

import (
	"ioutil"
	"log"
	"os"
)

// writeFile deals with writing the file, using the flags
// to know where and the filename / data to store it.
func writeFile(serverFlags Flags, filename string, data []byte) error {
	// TODO check the backend to use
	file, err := os.Create(serverFlags.OutputDirectory + "/" + filename)
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

// readFile is the method to read the file from wherever it
// is stored. The serverFlags are used to know where to read,
// the filename is used to know what to read.
func readFile(serverFlags Flags, filename string) ([]byte, error) {
	// TODO check the backend to use
	file, err := os.Open(s.Server.Flags.OutputDirectory + "/" + entry.Filename)
	if err != nil {
		w.WriteHeader(500)
		log.Println("[err] While requesting:", entry.Filename)
		log.Println(err)
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(500)
		log.Println("[err] While reading:", entry.Filename)
		log.Println(err)
	}
}
