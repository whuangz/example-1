FROM golang:1.16-alpine

WORKDIR /usr/app

ENV GO111MODULE=on
RUN go get github.com/githubnemo/CompileDaemon

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

ENTRYPOINT CompileDaemon --build="go build server.go"

