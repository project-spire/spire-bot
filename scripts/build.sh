#!bin/bash

protoc --proto_path=src --go_out=out --go_opt=paths=source_relative foo.proto bar/baz.proto