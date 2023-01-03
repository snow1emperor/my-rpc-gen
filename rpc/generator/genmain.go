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

//go:embed main.tpl
var mainTemplate string

//go:embed cmd.tpl
var cmdTemplate string

type MainServiceTemplateData struct {
	Service   string
	ServerPkg string
	Pkg       string
}

// GenMain generates the main file of the rpc service, which is an rpc service program call entry
func (g *Generator) GenMain(ctx DirContext, proto parser.Proto, cfg *conf.Config,
	c *ZRpcContext) error {

	// helper.go
	if err := func() error {
		head := util.GetHead(proto.Name)

		serviceImport := fmt.Sprintf(`"%v"`, ctx.GetService().Package)
		confImport := fmt.Sprintf(`"%v"`, ctx.GetConfig().Package)
		svcImport := fmt.Sprintf(`"%v"`, ctx.GetSvc().Package)

		imports := collection.NewSet()
		imports.AddStr(confImport, serviceImport, svcImport)
		serverFilename := "helper"
		serverFile := filepath.Join(ctx.GetMain().Filename, serverFilename+".go")
		text, err := pathx.LoadTemplate(category, mainTemplateFile, mainTemplate)
		if err != nil {
			return err
		}
		notStream := false

		return util.With("server").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
			"head":        head,
			"packageName": ctx.GetServiceName().Lower(),
			"imports":     strings.Join(imports.KeysStr(), pathx.NL),
			"funcs":       "",
			"notStream":   notStream,
		}, serverFile, false)
	}(); err != nil {
		return err
	}

	// main.go
	if err := func() error {
		head := util.GetHead(proto.Name)

		serverImport := fmt.Sprintf(`"%v"`, ctx.GetServer().Package)

		imports := collection.NewSet()
		imports.AddStr(serverImport, fmt.Sprintf(`"%s"`, c.VarStringCommandsPkg))
		serverFilename := "main"
		serverFile := filepath.Join(ctx.GetCmd().Filename, serverFilename+".go")
		text, err := pathx.LoadTemplate(category, cmdTemplateFile, cmdTemplate)
		if err != nil {
			return err
		}
		notStream := false

		return util.With("server").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
			"head":        head,
			"packageName": ctx.GetServiceName().Lower(),
			"imports":     strings.Join(imports.KeysStr(), pathx.NL),
			"funcs":       "",
			"notStream":   notStream,
		}, serverFile, false)
	}(); err != nil {
		return err
	}
	return nil
}
