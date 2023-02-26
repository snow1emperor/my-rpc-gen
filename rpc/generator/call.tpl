{{.head}}

package {{.filePackage}}

import (
	"context"

	{{.imports}}

	{{.pbPackage}}
	{{if ne .pbPackage .protoGoPackage}}{{.protoGoPackage}}{{end}}

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	{{.alias}}

	{{.serviceName}}Client interface {
		{{.interface}}
	}

	default{{.serviceName}}Client struct {
		cli zrpc.Client
	}
)

func New{{.serviceName}}Client(cli zrpc.Client) {{.serviceName}}Client {
	return &default{{.serviceName}}Client{
		cli: cli,
	}
}

{{.functions}}
