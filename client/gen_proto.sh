#!/bin/bash

TS_PLUGIN="$(npm root)/.bin/protoc-gen-ts_proto"
OUT_DIR="./proto"


# generate TS codes
protoc \
--plugin="${TS_PLUGIN}" \
--ts_proto_out="${OUT_DIR}" \
--ts_proto_opt="esModuleInterop=true" \
--ts_proto_opt="outputServices=grpc-js" \
--proto_path ../proto \
../proto/*.proto

