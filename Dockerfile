FROM golang:1.17-alpine AS builder
WORKDIR /app
COPY . /app
RUN chmod +x ./generate.sh && \
    /bin/sh ./generate.sh && \
    CGO_ENABLED=0 GOOS=linux \
    go build -a -o binary ./cmd/grpcgoonch
CMD ["./grpcgoonch"]