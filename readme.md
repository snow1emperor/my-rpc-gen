# Generate grpc service

```text
.
└── feeds
    ├── client
    │   └── feeds_client.go
    ├── cmd
    │   └── feeds
    │    │   └── main.go
    ├── etc
    │   └── feeds.yaml
    ├── feeds
    │   └── feeds.pb.go
    ├── helper.go
    └── internal
        ├── config
        │   └── config.go
        ├── core
        │   └── core.go
        │   └── feeds.getfeedlist_handler.go
        │   └── feeds.gethistorycounter_handler.go
        │   └── feeds.readhistory_handler.go
        │   └── feeds.updatefeedlist_handler.go
        ├── dao
        │   └── dao.go
        ├── server
        │   └── server.go
        │   └── grpc
        │       └── grpc.go
        │       └── service
        │           └── service.go
        │           └── feeds_service_impl.go
        └── svc   
            └── servicecontext.go
```

### Example:
```shell
PWD=`pwd`
GO_CTR_PATH=./goctl.go

GOPATH=$HOME/go
GOGOPROTO_PATH=$GOPATH/pkg/mod/github.com/gogo/protobuf@v1.3.2/protobuf

CUSTOM_PATH=example/.

SRC_PATH=example/feeds.proto
DST_DIR=./example/feeds

go run $GO_CTR_PATH rpc protoc $SRC_PATH \
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
```

### Clear
```shell
 bash example/clear.sh
```
