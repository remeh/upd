# Clioud

Upload from CLI, share with browsers.

## About

Clioud is the Rémy 'remeh' Mathieu's entry to the first GopherGala, 2015.

A client/server to upload and share files through http(s), with support of files auto-destruction with TTL.

I only had the possiblity to work on the first day of the competition, saturday, for about 9 hours.

## How to use

### Basics

First, you need to use the excellent dependency system gom:

```
go get github.com/mattn/gom
```

Then, in the main directory of `clioud`

```
gom install
```

to setup the dependencies.

### Start the daemon

To start the server:

```
gom build bin/server/server.go
./server
```

Available flags for the `server` executable:

```
-addr=":9000": The address to listen to with the server.
-cfile="": Path to a TLS certificate. Ex: ./certs/cert.pem
-ckey="": Path to a TLS key file. Ex: ./certs/key.pem
-key="": A shared secret key to identify the client.
-out="./": Directory in which the server can write the data.
-route="/clioud": Route served by the server.
```

### Upload a file with the client

Now that the server is up and running, you can upload files with this command:

```
gom build bin/client/client.go
./client file1 file2 file3 ...
```

it'll return the URL to share/delete the uploaded files. Example:
```
$ ./client --keep -ttl=4h README.md
For file : README.md
URL: http://localhost:9000/clioud/README.md
Delete URL: http://localhost:9000/clioud/README.md/ytGsotfcIUuZZ6eL
Available until: 2015-01-24 23:01:18.452801595 +0100 CET
```

Available flags for the `client` executable:

```
-ca="none": For HTTPS support: none / filename of an accepted CA / unsafe (doesn't check the CA)
-keep=false: Whether or not we must keep the filename
-key="": A shared secret key to identify the client.
-ttl="": TTL after which the file expires, ex: 30m. Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h"
-url="http://localhost:9000/clioud": The server to contact.
```

## Roadmap

  * Daemon listening for files *[ok]*
  * Daemon serving files *[ok]*
  * Client uploading files *[ok]*
    * Keepname option *[ok]*
    * TTL option *[ok]*
  * Secret key *[ok]*
  * TTL *[ok]*
  * Delete link *[ok]*
  * HTTPs *[ok]*

