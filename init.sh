#!/bin/bash
set -e

mkdir -p bin
wget https://github.com/protocolbuffers/protobuf/releases/download/v29.3/protoc-29.3-linux-x86_64.zip -P bin
unzip bin/protoc-29.3-linux-x86_64.zip bin/protoc

export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

git submodule update --remote --init
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

./gen.sh