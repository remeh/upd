# Clioud

Upload from CLI, share with browsers.

## How to use

### Basics

First, you need to use the excellent dependency system gom:

```
go get github.com/mattn/gom
```

Then, in the main directory:

```
gom install
```

to setup the dependencies.

### Start the daemon

To start the server:

```
gom run server/server.go
```

Available flags:

```
-addr=":9000": The address to listen to with the server.
-key="": The secret key to identify the client.
-out="./": Directory in which the server can write the data.
-route="/clioud": Route served by the server.
```

### Upload a file with the client

You can upload a files with this command:

```
gom run client/client.go file1 file2 file3
```

it'll return the URL to share/delete the uploaded files.

Available flags:

```
-key="": The secret key to identify the client.
-ttl="": TTL after which the file expires
-url="http://localhost:9000/clioud": The server to contact
```

## Roadmap

  * Daemon listening for files *[ok]*
  * Daemon serving files *[ok]*
  * Client uploading files *[ok]*
    * Keepname option
    * TTL option [*ok*]
  * Secret key *[ok]*
  * TTL *[ok]*
  * Delete link *[ok]*
  * TLS

