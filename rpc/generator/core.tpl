package {{.packageName}}

import (
    "context"

    {{.imports}}

    "github.com/zeromicro/go-zero/core/logx"
    "google.golang.org/grpc/metadata"
)

type {{.logicName}}Core struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    logx.Logger
    MD metadata.MD
}

func New{{.logicName}}(ctx context.Context, svcCtx *svc.ServiceContext) *{{.logicName}}Core {
    md, _ := metadata.FromIncomingContext(ctx)

    return &{{.logicName}}Core{
        ctx:    ctx,
        svcCtx: svcCtx,
        Logger: logx.WithContext(ctx),
        MD:     md,
    }
}
{{.functions}}
