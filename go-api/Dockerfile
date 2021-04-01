FROM golang:1.16-stretch

WORKDIR /usr/app

ENV GO111MODULE=on
RUN go get github.com/cosmtrek/air@v1.15.1
RUN go get -u github.com/pressly/goose/cmd/goose

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

CMD air


