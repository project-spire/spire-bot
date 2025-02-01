#!/bin/bash
set -e

mkdir -p gen

#PROTOCOLS=$(find protocol/msg -name "*.proto" -print)
#protoc -I=protocol/msg --go_out=gen --go_opt=paths=import "$PROTOCOLS"

find protocol/msg -name "*.proto" -print | while read -r proto; do
  bin/protoc -I=protocol/msg --go_out=gen --go_opt=paths=import "$proto"
done