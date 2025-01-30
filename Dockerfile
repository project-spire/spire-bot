FROM golang:latest

RUN apt-get update -y && \
    apt-get install -y unzip

WORKDIR /root
RUN wget https://github.com/protocolbuffers/protobuf/releases/download/v29.3/protoc-29.3-linux-x86_64.zip && \
    unzip protoc-29.3-linux-x86_64.zip -d protoc && \
    mv protoc/bin/protoc /usr/local/bin/

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

WORKDIR /workspace
COPY . .

RUN mkdir gen && \
    make gen