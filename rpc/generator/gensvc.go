package generator

import (
	_ "embed"
	"fmt"
	"github.com/zeromicro/go-zero/core/collection"
	"path/filepath"
	"strings"

	"github.com/snow1emperor/my-rpc-gen/rpc/parser"
	conf "github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

//go:embed svc.tpl
var svcTemplate string

//go:embed dao.tpl
var daoTemplate string

// GenSvc generates the servicecontext.go file, which is the resource dependency of a service,
// such as rpc dependency, model dependency, etc.
func (g *Generator) GenSvc(ctx DirContext, _ parser.Proto, cfg *conf.Config) error {
	dir := ctx.GetSvc()
	//svcFilename, err := format.FileNamingFormat(cfg.NamingFormat, "service_context")
	//if err != nil {
	//	return err
	//}

	fileName := filepath.Join(dir.Filename, "service_context"+".go")
	text, err := pathx.LoadTemplate(category, svcTemplateFile, svcTemplate)
	if err != nil {
		return err
	}

	confImport := fmt.Sprintf(`"%v"`, ctx.GetConfig().Package)
	daoImport := fmt.Sprintf(`"%v"`, ctx.GetDao().Package)
	imports := collection.NewSet()
	imports.AddStr(confImport, daoImport)

	if err := func() error {
		fileName := filepath.Join(ctx.GetDao().Filename, "dao"+".go")
		text, err := pathx.LoadTemplate(category, daoTemplateFile, daoTemplate)
		if err != nil {
			return err
		}
		imports := collection.NewSet()
		imports.AddStr(confImport)
		return util.With("dao").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
			"imports": strings.Join(imports.KeysStr(), pathx.NL),
		}, fileName, false)
	}(); err != nil {
		return err
	}

	return util.With("svc").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
		"imports": strings.Join(imports.KeysStr(), pathx.NL),
	}, fileName, false)
}
