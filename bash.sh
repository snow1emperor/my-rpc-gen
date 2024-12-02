PWD=`pwd`
GO_CTR_PATH=./bin/linuxapp

GOPATH=$HOME/go
GOGOPROTO_PATH=$GOPATH/pkg/mod/github.com/gogo/protobuf@v1.3.2/protobuf

CUSTOM_PATH=example/.

SRC_PATH=example/payments.tl.proto
DST_DIR=./example/payments


$GO_CTR_PATH rpc protoc $SRC_PATH \
  -I=$PWD \
  -I=$GOGOPROTO_PATH \
  -I=$CUSTOM_PATH \
  --gogo_dst=$DST_DIR \
  --gogo_out="plugins=grpc\,Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/api.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/field_mask.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/source_context.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/type.proto=github.com/gogo/protobuf/types,Mmy_svc.proto=github.com/snow1emperor/my-rpc-gen/example/my_svc/my_svc"\
  --zrpc_out=$DST_DIR \
  --commands_pkg="github.com/teamgram/marmota/pkg/commands" \
  --types="github.com/gogo/protobuf/types" \
  --type_map="my_svc=github.com/snow1emperor/my-rpc-gen/example/my_svc/my_svc" \
  --verbose

#rm -rf example/feeds
#
#SRC_DIR=./example
#DST_DIR=./example/feeds
#
#GOPATH=$HOME/go
#GOGOPROTO_PATH=$GOPATH/pkg/mod/github.com/gogo/protobuf@v1.3.2/protobuf
#
#protoc --proto_path=$GOGOPROTO_PATH:./ \
#    --gogo_out=plugins=grpc\,Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/api.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/field_mask.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/source_context.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/type.proto=github.com/gogo/protobuf/types,:$DST_DIR \
#        $SRC_DIR/*.proto

#MTPROTO_PATH=$GOPATH/pkg/mod/github.com/teamgram/proto@v0.150.0/mtproto

#protoc -I=$SRC_DIR:$MTPROTO_PATH --proto_path=$GOPATH/pkg/mod:$GOGOPROTO_PATH:./ \
#    --gogo_out=plugins=grpc\,Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/api.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/field_mask.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/source_context.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/type.proto=github.com/gogo/protobuf/types,:$DST_DIR \
#    $SRC_DIR/*.proto
#
