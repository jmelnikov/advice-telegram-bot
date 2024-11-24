FROM golang:1.23.3

COPY app /go/src/app

WORKDIR /go/src/app

RUN go get github.com/mattn/go-sqlite3
