{{.head}}

package service

import (
	{{if .notStream}}"context"{{end}}
	"github.com/zeromicro/go-zero/core/jsonx"

	{{.imports}}
)

func DebugString(in interface{}) string {
	s, _ := jsonx.MarshalToString(in)
	return s
}

{{.funcs}}
