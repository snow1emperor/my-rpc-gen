GO_CTR_PATH=/Volumes/EXTERNAL_DISK/GolandProjects/my-rpc-gen/goctl.go

go run $GO_CTR_PATH rpc protoc example/feeds.proto \
  --proto_path=${PWD} \
  --go_out=./example/feeds \
  --go-grpc_out=./example/feeds \
  --zrpc_out=./example/feeds

#rm -rf example/feeds
#
#SRC_DIR=./example
#DST_DIR=./example/feeds
#
#GOGOPROTO_PATH=$GOPATH/pkg/mod/github.com/gogo/protobuf@v1.3.2/protobuf
#MTPROTO_PATH=$GOPATH/pkg/mod/github.com/teamgram/proto@v0.150.0/mtproto
#
#protoc -I=$SRC_DIR:$MTPROTO_PATH --proto_path=$GOPATH/pkg/mod:$GOGOPROTO_PATH:./ \
#    --gogo_out=plugins=grpc,Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types,:$DST_DIR \
#    $SRC_DIR/*.proto
