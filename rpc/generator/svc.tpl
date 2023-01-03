package svc

import (
	{{.imports}}
)

type ServiceContext struct {
	Config config.Config
	Dao   *dao.Dao
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Dao:   dao.New(c),
	}
}
