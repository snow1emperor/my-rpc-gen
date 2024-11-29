PWD=`pwd`
GO_CTR_PATH=./bin/linuxapp

CUSTOM_PATH=example/.

SRC_PATH=example/payments.tl.proto
DST_DIR=./example/payments


$GO_CTR_PATH rpc protoc $SRC_PATH \
  -I=$PWD \
  -I=$CUSTOM_PATH \
  --gogo_dst=$DST_DIR \
  --zrpc_out=$DST_DIR \
  --commands_pkg="github.com/teamgram/marmota/pkg/commands" \
  --verbose
