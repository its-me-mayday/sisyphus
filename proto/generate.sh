#!/bin/bash
set -e

PROTO_DIR="$(cd "$(dirname "$0")" && pwd)"
OUT_DIR="$PROTO_DIR"

echo "→ generating tasks..."
protoc \
  --proto_path="$PROTO_DIR" \
  --go_out="$OUT_DIR" \
  --go_opt=paths=source_relative \
  --go-grpc_out="$OUT_DIR" \
  --go-grpc_opt=paths=source_relative \
  tasks/tasks.proto

echo "→ generating acl..."
protoc \
  --proto_path="$PROTO_DIR" \
  --go_out="$OUT_DIR" \
  --go_opt=paths=source_relative \
  --go-grpc_out="$OUT_DIR" \
  --go-grpc_opt=paths=source_relative \
  acl/acl.proto

echo "✓ done"