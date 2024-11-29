package generator

import (
	_ "embed"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/snow1emperor/my-rpc-gen/rpc/parser"
	"github.com/zeromicro/go-zero/core/collection"
	conf "github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"github.com/zeromicro/go-zero/tools/goctl/util/stringx"
)

const logicFunctionTemplate = `{{if .hasComment}}{{.comment}}{{end}}
func (c *{{.packageName}}Core) {{.method}} ({{if .hasReq}}in {{.request}}{{if .stream}},stream {{.streamBody}}{{end}}{{else}}stream {{.streamBody}}{{end}}) ({{if .hasReply}}res {{.response}},{{end}} err error) {
	// todo: add your logic here and delete this line
	{{if .stream}}return errors.New("Unimplemented"){{else}}
	return nil, errors.New("Unimplemented")
	return {{if .hasReply}}&{{.responseType}}{},{{end}} nil{{end}}
}
`

//go:embed core.tpl
var coreTemplate string

//go:embed logic.tpl
var logicTemplate string

// GenLogic generates the logic file of the rpc service, which corresponds to the RPC definition items in proto.
func (g *Generator) GenLogic(ctx DirContext, proto parser.Proto, cfg *conf.Config,
	c *ZRpcContext) error {
	if !c.Multiple {

		return g.genLogicInCompatibility(ctx, proto, cfg, c)
	}

	return g.genLogicGroup(ctx, proto, cfg, c)
}

func (g *Generator) genLogicInCompatibility(ctx DirContext, proto parser.Proto, cfg *conf.Config, c *ZRpcContext) error {
	dir := ctx.GetLogic()
	service := proto.Service[0].Service.Name

	// core.go
	if err := func() error {
		coreFilename := "core"
		filename := filepath.Join(dir.Filename, coreFilename+".go")
		imports := collection.NewSet()
		imports.AddStr(fmt.Sprintf(`"%v"`, ctx.GetSvc().Package))
		//imports.AddStr(fmt.Sprintf(`"%v"`, ctx.GetPb().Package))
		//imports.AddStr(fmt.Sprintf(`"%s"`, c.VarStringMTProtPkg))
		text, err := pathx.LoadTemplate(category, logicTemplateCoreFile, coreTemplate)
		if err != nil {
			return err
		}
		return util.With("logic").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
			"logicName":   fmt.Sprintf("%s", stringx.From(ctx.GetServiceName().Title()).ToCamel()),
			"functions":   "",
			"packageName": "core",
			"imports":     strings.Join(imports.KeysStr(), pathx.NL),
		}, filename, false)
	}(); err != nil {
		return err
	}

	for _, rpc := range proto.Service[0].RPC {
		logicName := fmt.Sprintf("%s", stringx.From(rpc.Name).ToCamel())
		logicFilename, err := format.FileNamingFormat(cfg.NamingFormat, rpc.Name)
		if err != nil {
			return err
		}
		if strings.Contains(logicFilename, ctx.GetServiceName().Source()) {
			logicFilename = strings.Replace(logicFilename, ctx.GetServiceName().Source(), "", 1)
			logicFilename = firstLetterToLower(logicFilename)
		}

		logicFilename = fmt.Sprintf("%s.%s_handler", ctx.GetServiceName().Source(), logicFilename)

		filename := filepath.Join(dir.Filename, logicFilename+".go")
		functions, impList, err := g.genLogicFunction(service, proto.PbPackage, logicName, rpc, c.VarStringTypeMap, ctx.GetPb().Package)
		if err != nil {
			return err
		}

		imports := collection.NewSet()
		//imports.AddStr(fmt.Sprintf(`"%v"`, ctx.GetSvc().Package))
		//imports.AddStr(fmt.Sprintf(`"%v"`, ctx.GetPb().Package))
		for _, item := range impList {
			imports.AddStr(fmt.Sprintf(`"%v"`, item))
		}
		text, err := pathx.LoadTemplate(category, logicTemplateFileFile, logicTemplate)
		if err != nil {
			return err
		}
		err = util.With("logic").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
			"logicName":   fmt.Sprintf("%sLogic", stringx.From(rpc.Name).ToCamel()),
			"functions":   functions,
			"packageName": "core",
			"imports":     strings.Join(imports.KeysStr(), pathx.NL),
		}, filename, false)
		if err != nil {
			return err
		}
	}
	return nil
}

func firstLetterToLower(s string) string {

	if len(s) == 0 {
		return s
	}

	r := []rune(s)
	r[0] = unicode.ToLower(r[0])

	return string(r)
}

func (g *Generator) genLogicGroup(ctx DirContext, proto parser.Proto, cfg *conf.Config, c *ZRpcContext) error {
	dir := ctx.GetLogic()
	for _, item := range proto.Service {
		serviceName := item.Name
		for _, rpc := range item.RPC {
			var (
				err           error
				filename      string
				logicName     string
				logicFilename string
				packageName   string
			)

			logicName = fmt.Sprintf("%sLogic", stringx.From(rpc.Name).ToCamel())
			childPkg, err := dir.GetChildPackage(serviceName)
			if err != nil {
				return err
			}

			serviceDir := filepath.Base(childPkg)
			nameJoin := fmt.Sprintf("%s_logic", serviceName)
			packageName = strings.ToLower(stringx.From(nameJoin).ToCamel())
			logicFilename, err = format.FileNamingFormat(cfg.NamingFormat, rpc.Name+"_logic")
			if err != nil {
				return err
			}

			pbImport := fmt.Sprintf(`"%v"`, ctx.GetPb().Package)

			filename = filepath.Join(dir.Filename, serviceDir, logicFilename+".go")
			functions, impList, err := g.genLogicFunction(serviceName, proto.PbPackage, logicName, rpc, c.VarStringTypeMap, pbImport)
			if err != nil {
				return err
			}

			imports := collection.NewSet()
			imports.AddStr(fmt.Sprintf(`"%v"`, ctx.GetSvc().Package))
			//imports.AddStr(fmt.Sprintf(`"%v"`, ctx.GetPb().Package))
			for _, item := range impList {
				imports.AddStr(fmt.Sprintf(`"%v"`, item))
			}
			text, err := pathx.LoadTemplate(category, logicTemplateFileFile, logicTemplate)
			if err != nil {
				return err
			}

			if err = util.With("logic").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
				"logicName":   logicName,
				"functions":   functions,
				"packageName": packageName,
				"imports":     strings.Join(imports.KeysStr(), pathx.NL),
			}, filename, false); err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *Generator) genLogicFunction(serviceName, goPackage, logicName string, rpc *parser.RPC, typeMap map[string]string, pbImport string) (string, []string, error) {
	var impList []string

	functions := make([]string, 0)
	text, err := pathx.LoadTemplate(category, logicFuncTemplateFileFile, logicFunctionTemplate)
	if err != nil {
		return "", impList, err
	}

	comment := parser.GetComment(rpc.Doc())
	streamServer := fmt.Sprintf("%s.%s_%s%s", goPackage, parser.CamelCase(serviceName),
		parser.CamelCase(rpc.Name), "Server")

	request := func() string {
		var mess = rpc.RequestType
		if strings.Contains(mess, "google.protobuf") {
			if path, ok := typeMap["types"]; ok {
				impList = append(impList, path)
			}
			values := strings.Split(mess, ".")
			return fmt.Sprintf("*%s.%s", "types", values[len(values)-1])
		} else if strings.Contains(mess, ".") && !strings.HasPrefix(mess, goPackage+".") {
			pkgKey := strings.Split(mess, ".")[0]
			if path, ok := typeMap[pkgKey]; ok {
				impList = append(impList, path)
				return fmt.Sprintf("*%s", mess)
			} else {
				err = errors.New(fmt.Sprintf("request type %s must defined in flags type_map", pkgKey))
				return ""
			}
		} else if strings.HasPrefix(mess, goPackage+".") {
			mess = strings.Split(mess, ".")[1]
		}
		impList = append(impList, pbImport)
		return fmt.Sprintf("*%s.%s", goPackage, parser.CamelCase(mess))
	}()
	if err != nil {
		return "", nil, err
	}

	response := func() string {
		var mess = rpc.ReturnsType
		if strings.Contains(mess, "google.protobuf") {
			if path, ok := typeMap["types"]; ok {
				impList = append(impList, path)
			}
			values := strings.Split(mess, ".")
			return fmt.Sprintf("%s.%s", "types", values[len(values)-1])
		} else if strings.Contains(mess, ".") && !strings.HasPrefix(mess, goPackage+".") {
			pkgKey := strings.Split(mess, ".")[0]
			if path, ok := typeMap[pkgKey]; ok {
				impList = append(impList, path)
				return mess
			} else {
				err = errors.New(fmt.Sprintf("request package %s must defined in flags type_map", pkgKey))
				return ""
			}
		} else if strings.HasPrefix(mess, goPackage+".") {
			mess = strings.Split(mess, ".")[1]
		}
		impList = append(impList, pbImport)
		return fmt.Sprintf("%s.%s", goPackage, parser.CamelCase(mess))
	}()
	if err != nil {
		return "", nil, err
	}

	buffer, err := util.With("fun").Parse(text).Execute(map[string]interface{}{
		"packageName":  parser.CamelCase(goPackage),
		"logicName":    logicName,
		"method":       parser.CamelCase(rpc.Name),
		"hasReq":       !rpc.StreamsRequest,
		"request":      request,
		"hasReply":     !rpc.StreamsRequest && !rpc.StreamsReturns,
		"response":     fmt.Sprintf("*%s", response),
		"responseType": response,
		"stream":       rpc.StreamsRequest || rpc.StreamsReturns,
		"streamBody":   streamServer,
		"hasComment":   len(comment) > 0,
		"comment":      comment,
	})
	if err != nil {
		return "", impList, err
	}

	functions = append(functions, buffer.String())
	return strings.Join(functions, pathx.NL), impList, nil
}
