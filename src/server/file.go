// Methods dealing with files writing/reading
// from either the FS or a distant service (s3, ...)
// Copyright © 2015 - Rémy MATHIEU
package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/gen/s3"
)

// writeFile deals with writing the file, using the flags
// to know where and the filename / data to store it.
func (s *Server) WriteFile(filename string, data []byte) error {
	if s.Config.Storage == FS_STORAGE {
		file, err := os.Create(s.Config.FSConfig.OutputDirectory + "/" + filename)
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
	} else if s.Config.Storage == S3_STORAGE {
		// S3 connection
		creds := aws.Creds(s.Config.S3Config.AccessKey, s.Config.S3Config.AccessSecret, "")
		client := s3.New(creds, s.Config.S3Config.Region, nil)
		body := ioutil.NopCloser(bytes.NewBuffer(data))

		// TODO TTL

		// creates the S3 put request
		por := &s3.PutObjectRequest{
			Body:          body,
			Key:           aws.String(filename),
			ContentLength: aws.Long(int64(len(data))),
			Bucket:        aws.String(s.Config.S3Config.Bucket),
		}

		_, err := client.PutObject(por)
		if err != nil {
			return err
		}
	}

	return fmt.Errorf("[err] Unsupported storage: %s", s.Config.Storage)
}

// readFile is the method to read the file from wherever it
// is stored. The serverFlags are used to know where to read,
// the filename is used to know what to read.
func (s *Server) ReadFile(filename string) ([]byte, error) {
	if s.Config.Storage == FS_STORAGE {
		file, err := os.Open(s.Config.FSConfig.OutputDirectory + "/" + filename)
		if err != nil {
			return nil, err
		}

		data, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}

		return data, nil
	} else if s.Config.Storage == S3_STORAGE {
		// S3 connection
		creds := aws.Creds(s.Config.S3Config.AccessKey, s.Config.S3Config.AccessSecret, "")
		client := s3.New(creds, s.Config.S3Config.Region, nil)

		gor := &s3.GetObjectRequest{
			Key:    aws.String(filename),
			Bucket: aws.String(s.Config.S3Config.Bucket),
		}

		resp, err := client.GetObject(gor)
		if err != nil {
			return nil, err
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("[err] Can't read the body of a GetObjectOutput from AWS")
			log.Println(err)
			return nil, err
		}

		return data, nil
	}

	return nil, fmt.Errorf("[err] Unsupported storage: %s", s.Config.Storage)
}

// Expire expires a file : delete it from the metadata
// and from the FS.
func (s *Server) Expire(m Metadata) error {
	delete(s.Metadata.Data, m.Filename)
	if s.Config.Storage == FS_STORAGE {
		return os.Remove(s.Config.RuntimeDir + "/" + m.Filename)
	}

	return fmt.Errorf("[err] Unsupported storage: %s", s.Config.Storage)
}
