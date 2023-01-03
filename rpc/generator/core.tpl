package {{.packageName}}

import (
    "context"

    {{.imports}}

    "github.com/zeromicro/go-zero/core/logx"
)

type {{.logicName}}Core struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    logx.Logger
    MD *metadata.RpcMetadata
}

func New{{.logicName}}(ctx context.Context, svcCtx *svc.ServiceContext) *{{.logicName}}Core {
    return &{{.logicName}}Core{
        ctx:    ctx,
        svcCtx: svcCtx,
        Logger: logx.WithContext(ctx),
        MD:     metadata.RpcMetadataFromIncoming(ctx),
    }
}
{{.functions}}
