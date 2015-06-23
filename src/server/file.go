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
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (s *Server) createS3Config(creds *credentials.Credentials, region string) *aws.Config {
	return &aws.Config{
		Credentials: creds,
		Region:      region,
	}
}

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
		creds := credentials.NewStaticCredentials(s.Config.S3Config.AccessKey, s.Config.S3Config.AccessSecret, "")
		client := s3.New(s.createS3Config(creds, s.Config.S3Config.Region))
		body := bytes.NewReader(data)

		// Creates the S3 put request
		por := &s3.PutObjectInput{
			Body:          body,
			Key:           aws.String(filename),
			ContentLength: aws.Long(int64(len(data))),
			Bucket:        aws.String(s.Config.S3Config.Bucket),
		}

		// Sends the S3 put request
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
		creds := credentials.NewStaticCredentials(s.Config.S3Config.AccessKey, s.Config.S3Config.AccessSecret, "")
		client := s3.New(s.createS3Config(creds, s.Config.S3Config.Region))

		// The get request
		gor := &s3.GetObjectInput{
			Key:    aws.String(filename),
			Bucket: aws.String(s.Config.S3Config.Bucket),
		}

		// Sends the request
		resp, err := client.GetObject(gor)
		if err != nil {
			return nil, err
		}

		// Reads the result
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
	filename := m.Filename

	// delete from the datbase
	s.deleteMetadata(filename)

	if s.Config.Storage == FS_STORAGE {
		return os.Remove(s.Config.RuntimeDir + "/" + filename)
	} else if s.Config.Storage == S3_STORAGE {
		// S3 connection
		creds := credentials.NewStaticCredentials(s.Config.S3Config.AccessKey, s.Config.S3Config.AccessSecret, "")
		client := s3.New(s.createS3Config(creds, s.Config.S3Config.Region))

		// The get request
		dor := &s3.DeleteObjectInput{
			Key:    aws.String(filename),
			Bucket: aws.String(s.Config.S3Config.Bucket),
		}

		_, err := client.DeleteObject(dor)
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("[err] Unsupported storage: %s", s.Config.Storage)
}

// computeEndOfLife return as a string the end of life of the new file.
func (s *Server) computeEndOfLife(ttl string, now time.Time) time.Time {
	if len(ttl) == 0 {
		return time.Time{}
	}
	duration, _ := time.ParseDuration(ttl) // no error possible 'cause already checked in the controller
	t := now.Add(duration)
	return t
}
