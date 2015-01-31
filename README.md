# upd

Upload from CLI, share with browsers.

## About

upd is file upload service to quickly share files through http(s) supporting different storage backend (filesystem and Amazon S3).

The server provides a simple API to easily allow the creation of clients.

A command-line client is also provided in this repo.


## Features

  * Storages backend : Filesystem, Amazon S3
  * Daemon listening to receive files 
  * Daemon serving files
  * TTL for expiration of files.
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
-keep=false: Whether or not we must keep the filename
-key="": A shared secret key to identify the client.
-ttl="": TTL after which the file expires, ex: 30m. Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h"
-url="http://localhost:9000/upd": The server to contact.
```
