{{.head}}

package {{.packageName}}

import (
	"context"
	"github.com/gogo/protobuf/proto"
	"github.com/zeromicro/go-zero/core/logx"

	{{.imports}}

)


type newRPCReplyFunc func() proto.Message

type RPCContextTuple struct {
	Method       string
	NewReqFunc   newRPCReplyFunc
	NewReplyFunc newRPCReplyFunc
}

func Fn(m proto.Message) newRPCReplyFunc {
	return func() proto.Message {
		return m
	}
}

func FindRPCContextTuple(constructor int64) *RPCContextTuple {
	m, ok := rpcContextRegisters[TLConstructor(constructor)]
	if !ok {
		logx.Errorf("Can't find constructor: %s", constructor)
		return nil
	}
	return &m
}

func GetRPCContextRegisters() map[TLConstructor]RPCContextTuple {
	return rpcContextRegisters
}

var rpcContextRegisters = map[TLConstructor]RPCContextTuple{
	{{.paths}}
}
