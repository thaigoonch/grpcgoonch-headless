#!/bin/sh

set -e

go mod vendor

SERVICE="grpcgoonch"

protoc.sh \
    --go_out=service \
    --proto_path=${ROOT_DIR}/vendor \
    --go_opt=paths=source_relative \
    --go-grpc_out=service \
    --go-grpc_opt=paths=source_relative \
    service.proto