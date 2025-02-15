#!/bin/bash

PROTO_OUT_DIR="./proto"

protoc \
--go_out=${PROTO_OUT_DIR} \
--go_opt=paths=source_relative \
--go-grpc_out=${PROTO_OUT_DIR} \
--go-grpc_opt=paths=source_relative \
--proto_path ../proto \
../proto/*.proto

