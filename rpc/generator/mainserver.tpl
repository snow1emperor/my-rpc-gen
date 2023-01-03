{{.head}}

package server

import (
    "flag"

    {{.imports}}

    "github.com/zeromicro/go-zero/core/conf"
    "github.com/zeromicro/go-zero/core/logx"
    "github.com/zeromicro/go-zero/zrpc"
)

var configFile = flag.String("f", "etc/{{.filename}}.yaml", "the config file")

type Server struct {
    grpcSrv *zrpc.RpcServer
}

func New() *Server {
    return new(Server)
}

func (s *Server) Initialize() error {
    var c config.Config
    conf.MustLoad(*configFile, &c)

    logx.Infov(c)
    ctx := svc.NewServiceContext(c)
    s.grpcSrv = grpc.New(ctx, c.RpcServerConf)

    go func() {
    go s.grpcSrv.Start()
    }()
    return nil
}

func (s *Server) RunLoop() {
}

func (s *Server) Destroy() {
    s.grpcSrv.Stop()
}
