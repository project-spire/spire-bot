#!/bin/bash
#set -e -x

export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

mkdir -p gen/msg

PROTO_PATH="protocol/msg"
PROTOS=$(find ${PROTO_PATH} -name "*.proto")
cmd="bin/protoc -I=${PROTO_PATH} --go_out=gen/msg --go_opt=paths=source_relative"

for proto in ${PROTOS}; do
  relative_proto=$(realpath --relative-to="${PROTO_PATH}" "${proto}")
  module_proto=$(dirname "${relative_proto}")
  if [ "${module_proto}" = "." ] ; then
    module_proto=""
  fi

  cmd+=" --go_opt=M${relative_proto}=spire/bot/gen/msg/${module_proto}"
done

for proto in ${PROTOS}; do
  cmd+=" ${proto}"
done

echo "${PROTOS}"
#echo "${cmd}"

eval "${cmd}"