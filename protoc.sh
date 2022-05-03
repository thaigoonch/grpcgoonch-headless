#!/bin/sh
set -x
set -e

protoc \
  -I. \
  $@