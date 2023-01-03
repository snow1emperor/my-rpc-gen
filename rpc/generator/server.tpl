{{.head}}

package service

import (
	{{if .notStream}}"context"{{end}}

	{{.imports}}
)

{{.funcs}}
