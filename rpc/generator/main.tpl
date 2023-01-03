package {{.packageName}}_helper

import (
	{{.imports}}
)

type (
	Config = config.Config
)

func New(c Config) *service.Service {
	return service.New(svc.NewServiceContext(c))
}
