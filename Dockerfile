FROM golang:latest

RUN apt-get update

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest