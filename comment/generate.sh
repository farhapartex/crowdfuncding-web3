#!/bin/sh
set -e

cd "$(dirname "$0")"

protoc \
  --go_out=. --go_opt=module=comment \
  --go-grpc_out=. --go-grpc_opt=module=comment \
  proto/comment.proto
