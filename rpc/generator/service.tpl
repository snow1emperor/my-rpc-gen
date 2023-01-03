{{.head}}

package service

import (
    {{if .notStream}}"context"{{end}}

    {{.imports}}
)

type Service struct {
    svcCtx *svc.ServiceContext
}

func New(svcCtx *svc.ServiceContext) *Service {
    return &Service{
        svcCtx: svcCtx,
    }
}
