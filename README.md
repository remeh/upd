# upd

Upload from CLI, share with browsers.

## About

upd is a file upload service to quickly share files through http(s) supporting different storage backend (filesystem and Amazon S3).

The server provides a simple API to easily allow the creation of clients and a command-line client is also provided in this repository.


## Features

  * Storages backend : Filesystem, Amazon S3
  * Daemon listening to receive files 
  * Daemon serving files
  * TTL for expiration of files.
  * Tags on files + search by tags API
  * Delete link 
  * HTTPs 
  * Secret shared key between client / server
  * Get last uploaded files
  * Routine job cleaning the expired files

## How to use

### Basics

First, you need to use the excellent dependency system gom:

```
go get github.com/mattn/gom
```

Then, in the main directory of `upd`

```
gom install
```

to setup the dependencies.

### Start the daemon

#### Normal server

To start the server:

```
gom build bin/server/server.go
./server
```

Available flags for the `server` executable:

```
-c="server.conf"": Path to a configuration file.
```

The configuration file is well-documented.

#### Docker server

upd daemon is ready to be launched with `Docker`. You must first build the docker container, in the upd directory :

```
docker build -t upd .
```

It'll build the upd server docker container. What you must know:
  * The `ENTRYPOINT` docker has been binded on the configuration file `/etc/upd/server.conf`, this way, by using a volume, you can provide your configuration file.
  * Don't forget to bind a volume for the data directory if you're using the filesystem storage backend. If you don't do so, you'll lost your data when the docker will be stopped/restarted.

Example of how to launch the upd container (with the `server.conf` in your host `/home/user/`) :

```
docker run --rm -ti -v /home/user:/etc/upd -p 9000:9000 upd
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
URL: http://localhost:9000/upd/README.md
Delete URL: http://localhost:9000/upd/README.md/ytGsotfcIUuZZ6eL
Available until: 2015-01-24 23:01:18.452801595 +0100 CET
```

Available flags for the `client` executable:

```
-ca="none": For HTTPS support: none / filename of an accepted CA / unsafe (doesn't check the CA)
-key="": A shared secret key to identify the client.
-tags="": Tag the files. Ex: -tags="screenshot,may"
-ttl="": TTL after which the file expires, ex: 30m. Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h"
-url="http://localhost:9000/upd": The upd server to contact.
```
