#!/bin/bash

set -e

go mod vendor

protoc \
    -I. \
    -I/include/proto \
    --go_out=service \
    --proto_path=$(pwd) \
    --go_opt=paths=source_relative \
    --go-grpc_out=service \
    --go-grpc_opt=paths=source_relative \
    service.proto
