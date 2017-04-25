// Server executable to receive/host files.
// Copyright © 2015 - Rémy MATHIEU

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/remeh/upd/src/server"

	"github.com/BurntSushi/toml"
)

// readFromFile looks whether a configuration file is provided
// to read the config into it.
func readFromFile(filename string) (server.Config, error) {
	// default hardcoded config when no configuration file is available
	config := server.Config{
		Addr:       ":9000",
		Storage:    "fs",
		RuntimeDir: "/tmp",
		FSConfig: server.FSConfig{
			OutputDirectory: "/tmp",
		},
		Route: "/upd",
	}

	// read the file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}

	// decode the file
	_, err = toml.Decode(string(data), &config)
	if err != nil {
		return config, err
	}

	// Ensure the validity of the config // TODO log
	if config.Route[0] != '/' {
		config.Route = "/" + config.Route
	}
	if config.Route[len(config.Route)-1] == '/' {
		config.Route = config.Route[:len(config.Route)-1]
	}

	if config.Storage != server.FS_STORAGE && config.Storage != server.S3_STORAGE {
		log.Println("[err] Unknown storage:", config.Storage)
		os.Exit(1)
	}

	return config, nil
}

func parseFlags() server.Flags {
	var flags server.Flags

	// Declare the flags
	flag.StringVar(&(flags.ConfigFile), "c", "upd.conf", "Configuration file to use.")

	// Read them
	flag.Parse()

	return flags
}

func main() {
	flags := parseFlags()

	config, err := readFromFile(flags.ConfigFile)
	if err != nil {
		log.Println("[warn] Can't read the configuration file")
		log.Println("[warn] Falling back on default values for configuration.")
	}

	app := server.NewServer(config)
	app.Start()
}
