FROM golang

ADD . /go/src/github.com/remeh/clioud

RUN go get github.com/mattn/gom

RUN cd /go/src/github.com/remeh/clioud && gom install && gom build bin/server/server.go

ENTRYPOINT /go/src/github.com/remeh/clioud/server
