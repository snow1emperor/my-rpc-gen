PWD=`pwd`
GO_CTR_PATH=./bin/macapp

GOPATH=$HOME/go

DST_DIR=./example

$GO_CTR_PATH rpc proxy \
  --proxy_dst=$DST_DIR \
  --proxy_out="$PWD/example/my_svc.proto,$PWD/example/idgen.proto,$PWD/example/gateway.proto,$PWD/example/feeds.proto,$PWD/rpc/parser/stream.proto,$PWD/rpc/parser/test.proto,$PWD/rpc/parser/stream.proto,$PWD/rpc/parser/test.proto"\
  --proxy_types="github.com/gogo/protobuf/types" \
  --proxy_type_map="my_svc=github.com/snow1emperor/my-rpc-gen/example/my_svc/my_svc,feeds=github.com/snow1emperor/my-rpc-gen/example/feeds/feeds,gateway=github.com/snow1emperor/my-rpc-gen/example/gateway/gateway,idgen=github.com/snow1emperor/my-rpc-gen/example/idgen/idgen,test=github.com/snow1emperor/my-rpc-gen/rpc/parser"

