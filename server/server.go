// Server executable to receive/host files.
// Copyright © 2015 - Rémy MATHIEU

package main

import (
	"flag"

	"server"
)

func parseFlags() server.Flags {
	var flags server.Flags

	// Declare the flags

	flag.StringVar(&(flags.Addr), "addr", ":9000", "The address to listen to with the server.")
	flag.StringVar(&(flags.SecretKey), "key", "", "The secret key to identify the client.")
	flag.StringVar(&(flags.OutputDirectory), "out", "./", "Directory in which the server can write the data.")
	flag.StringVar(&(flags.Route), "route", "/clioud", "Route served by the server.")
	flag.StringVar(&(flags.CertificateFile), "cfile", "", "Path to a TLS certificate. Ex: ./certs/cert.pem")
	flag.StringVar(&(flags.CertificateKey), "ckey", "", "Path to a TLS key file. Ex: ./certs/key.pem")

	// Read them
	flag.Parse()

	// Ensure the validity of the flags
	if flags.Route[0] != '/' {
		flags.Route = "/" + flags.Route
	}
	if flags.Route[len(flags.Route)-1] == '/' {
		flags.Route = flags.Route[:len(flags.Route)-1]
	}

	return flags
}

func main() {
	flags := parseFlags()

	app := server.NewServer(flags)
	app.Start()
}
