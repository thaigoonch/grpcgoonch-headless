#!/bin/sh

set -e

go mod vendor

path="$(pwd)"
if [[ ${HOST_DIR} ]]; then
   path=${HOST_DIR}
fi
ROOT_DIR=${ROOT_DIR:-${path}}
SERVICE="grpcgoonch"

protoc \
    --go_out=service \
    --proto_path=${ROOT_DIR} \
    --go_opt=paths=source_relative \
    --go-grpc_out=service \
    --go-grpc_opt=paths=source_relative \
    service.proto