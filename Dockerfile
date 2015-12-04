FROM golang:1.5

ENV GO15VENDOREXPERIMENT 1

COPY 	. /go/src/github.com/remeh/upd
WORKDIR /go/src/github.com/remeh/upd

RUN go build  -v -o /upd-server bin/server/server.go

EXPOSE 9000

ENTRYPOINT ["/upd-server"]
CMD ["-c", "/etc/upd/server.conf"]
