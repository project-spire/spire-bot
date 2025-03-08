#!/bin/bash
#set -e -x

export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

GEN_PATH="gen/protocol"
PROTO_PATH="protocol"
PROTOS=$(find ${PROTO_PATH} -name "*.proto")
cmd="bin/protoc -I=${PROTO_PATH} --go_out=${GEN_PATH} --go_opt=paths=source_relative"

mkdir -p ${GEN_PATH}

for proto in ${PROTOS}; do
  relative_proto=$(realpath --relative-to="${PROTO_PATH}" "${proto}")
  module_proto=$(dirname "${relative_proto}")
  if [ "${module_proto}" = "." ] ; then
    module_proto=""
  fi

  cmd+=" --go_opt=M${relative_proto}=spire/bot/${GEN_PATH}/${module_proto}"
done

for proto in ${PROTOS}; do
  cmd+=" ${proto}"
done

echo "${PROTOS}"
#echo "${cmd}"

eval "${cmd}"