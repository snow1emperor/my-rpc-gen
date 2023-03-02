#!/bin/sh
SRC_DIR=.
DST_DIR=.

GOPATH=$HOME/go
GOGOPROTO_PATH=$GOPATH/pkg/mod/github.com/gogo/protobuf@v1.3.2/protobuf
GOGOPROTO_PATH2=$GOPATH/pkg/mod/github.com/gogo/protobuf@v1.3.2/gogoproto


protoc -I=$SRC_DIR:./ \
    --proto_path=$GOGOPROTO_PATH:./ \
    --proto_path=$GOGOPROTO_PATH2:./ \
    --gogo_out=plugins=grpc,Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types,Mgogo.proto=github.com/gogo/protobuf/gogoproto,:$DST_DIR \
    $SRC_DIR/proxy_constructors.proto
