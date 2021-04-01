FROM golang:1.16-alpine

WORKDIR /usr/app

ENV GO111MODULE=on
RUN go get github.com/cosmtrek/air@v1.15.1

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

CMD air


