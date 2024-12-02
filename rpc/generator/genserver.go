package generator

import (
	_ "embed"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/snow1emperor/my-rpc-gen/rpc/parser"
	"github.com/zeromicro/go-zero/core/collection"
	conf "github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"github.com/zeromicro/go-zero/tools/goctl/util/stringx"
)

const functionTemplate = `
{{if .hasComment}}{{.comment}}{{end}}
func (s *Service) {{.method}} ({{if .notStream}}ctx context.Context,{{if .hasReq}} in {{.request}}{{end}}{{else}}{{if .hasReq}} in {{.request}},{{end}}stream {{.streamBody}}{{end}}) ({{if .notStream}}{{.response}},{{end}}error) {
	c := {{.logicPkg}}.New{{.logicName}}({{if .notStream}}ctx,{{else}}stream.Context(),{{end}}s.svcCtx)

	c.Logger.Debugf("{{.handler}} - metadata: %s, request: {{if .notStream}}%s{{else}}stream{{end}}", DebugString(c.MD){{if .notStream}}, DebugString(in){{end}})  
	{{if .notStream}}r, {{end}}err := c.{{.method}}({{if .hasReq}}in{{if .stream}} ,stream{{end}}{{else}}{{if .stream}}stream{{end}}{{end}})
	if err != nil {
		c.Logger.Debugf("{{.handler}} - error: %s", err.Error())
		return {{if .notStream}}nil, {{end}}err
	}

	{{if .notStream}}c.Logger.Debugf("{{.handler}} - reply: %s", DebugString(r)){{end}}
	return {{if .notStream}}r, {{end}}err
}
`

//go:embed server.tpl
var serverTemplate string

//go:embed mainserver.tpl
var serverMainTemplate string

//go:embed service.tpl
var serviceTemplate string

//go:embed grpc.tpl
var grpcTemplate string

// GenServer generates rpc server file, which is an implementation of rpc server
func (g *Generator) GenServer(ctx DirContext, proto parser.Proto, cfg *conf.Config,
	c *ZRpcContext) error {
	if !c.Multiple {
		return g.genServerInCompatibility(ctx, proto, cfg, c)
	}

	return g.genServerGroup(ctx, proto, cfg, c)
}

func (g *Generator) genServerGroup(ctx DirContext, proto parser.Proto, cfg *conf.Config, c *ZRpcContext) error {
	dir := ctx.GetServer()
	for _, service := range proto.Service {
		var (
			serverFile  string
			logicImport string
		)

		serverFilename, err := format.FileNamingFormat(cfg.NamingFormat, service.Name+"_server")
		if err != nil {
			return err
		}

		serverChildPkg, err := dir.GetChildPackage(service.Name)
		if err != nil {
			return err
		}

		logicChildPkg, err := ctx.GetLogic().GetChildPackage(service.Name)
		if err != nil {
			return err
		}

		serverDir := filepath.Base(serverChildPkg)
		logicImport = fmt.Sprintf(`"%v"`, logicChildPkg)
		serverFile = filepath.Join(dir.Filename, serverDir, serverFilename+".go")

		svcImport := fmt.Sprintf(`"%v"`, ctx.GetSvc().Package)
		pbImport := fmt.Sprintf(`"%v"`, ctx.GetPb().Package)

		imports := collection.NewSet()
		imports.AddStr(logicImport, svcImport)

		head := util.GetHead(proto.Name)

		funcList, impList, err := g.genFunctions(proto.PbPackage, service, true, c.VarStringTypeMap, pbImport)
		if err != nil {
			return err
		}

		for _, item := range impList {
			imports.AddStr(fmt.Sprintf(`"%v"`, item))
		}

		text, err := pathx.LoadTemplate(category, serverTemplateFile, serverTemplate)
		if err != nil {
			return err
		}

		notStream := false
		for _, rpc := range service.RPC {
			if !rpc.StreamsRequest && !rpc.StreamsReturns {
				notStream = true
				break
			}
		}

		if err = util.With("server").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
			"head": head,
			"unimplementedServer": fmt.Sprintf("%s.Unimplemented%sServer", proto.PbPackage,
				stringx.From(service.Name).ToCamel()),
			"server":    stringx.From(service.Name).ToCamel(),
			"imports":   strings.Join(imports.KeysStr(), pathx.NL),
			"funcs":     strings.Join(funcList, pathx.NL),
			"notStream": notStream,
		}, serverFile, true); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) genServerInCompatibility(ctx DirContext, proto parser.Proto,
	cfg *conf.Config, c *ZRpcContext) error {
	dir := ctx.GetService()
	logicImport := fmt.Sprintf(`"%v"`, ctx.GetLogic().Package)
	svcImport := fmt.Sprintf(`"%v"`, ctx.GetSvc().Package)
	serviceImport := fmt.Sprintf(`"%v"`, dir.Package)
	pbImport := fmt.Sprintf(`"%v"`, ctx.GetPb().Package)

	imports := collection.NewSet()
	imports.AddStr(logicImport) //svcImport,

	head := util.GetHead(proto.Name)
	service := proto.Service[0]

	// service.go
	if err := func() error {
		imports := collection.NewSet()
		imports.AddStr(svcImport)
		serverFilename := "service"
		serverFile := filepath.Join(dir.Filename, serverFilename+".go")
		text, err := pathx.LoadTemplate(category, serverSrvTemplateFile, serviceTemplate)
		if err != nil {
			return err
		}
		notStream := false

		return util.With("server").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
			"head": head,
			"unimplementedServer": fmt.Sprintf("%s.Unimplemented%sServer", proto.PbPackage,
				stringx.From(service.Name).ToCamel()),
			"server":    stringx.From(service.Name).ToCamel(),
			"imports":   strings.Join(imports.KeysStr(), pathx.NL),
			"funcs":     "",
			"notStream": notStream,
		}, serverFile, true)
	}(); err != nil {
		return err
	}
	// grpc.go
	if err := func() error {
		imports := collection.NewSet()
		imports.AddStr(pbImport, serviceImport, svcImport)
		serverFilename := "grpc"
		serverFile := filepath.Join(ctx.GetGRPC().Filename, serverFilename+".go")
		text, err := pathx.LoadTemplate(category, serverRpcTemplateFile, grpcTemplate)
		if err != nil {
			return err
		}
		notStream := false

		return util.With("server").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
			"head": head,
			"unimplementedServer": fmt.Sprintf("%s.Unimplemented%sServer", proto.PbPackage,
				stringx.From(service.Name).ToCamel()),
			"grpcServer": fmt.Sprintf("%s.Register%sServer", ctx.GetServiceName().ToCamel(), service.Name),
			"imports":    strings.Join(imports.KeysStr(), pathx.NL),
			"funcs":      "",
			"notStream":  notStream,
		}, serverFile, true)
	}(); err != nil {
		return err
	}

	// server.go
	if err := func() error {
		grpcImport := fmt.Sprintf(`"%v"`, ctx.GetGRPC().Package)
		confImport := fmt.Sprintf(`"%v"`, ctx.GetConfig().Package)
		imports := collection.NewSet()
		imports.AddStr(confImport, grpcImport, svcImport)
		serverFilename := "server"
		serverFile := filepath.Join(ctx.GetServer().Filename, serverFilename+".go")
		text, err := pathx.LoadTemplate(category, serverMainTemplateFile, serverMainTemplate)
		if err != nil {
			return err
		}
		notStream := false

		return util.With("server").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
			"head": head,
			"unimplementedServer": fmt.Sprintf("%s.Unimplemented%sServer", proto.PbPackage,
				stringx.From(service.Name).ToCamel()),
			"filename":  ctx.GetServiceName().Lower(),
			"imports":   strings.Join(imports.KeysStr(), pathx.NL),
			"funcs":     "",
			"notStream": notStream,
		}, serverFile, false)
	}(); err != nil {
		return err
	}
	//serverFilename, err := format.FileNamingFormat(cfg.NamingFormat, service.Name+"_server")
	//if err != nil {
	//	return err
	//}
	serverFilename := ctx.GetServiceName().Lower() + "_service_impl"
	serverFile := filepath.Join(dir.Filename, serverFilename+".go")
	funcList, impList, err := g.genFunctions(proto.PbPackage, service, false, c.VarStringTypeMap, ctx.GetPb().Package)
	if err != nil {
		return err
	}

	for _, item := range impList {
		imports.AddStr(fmt.Sprintf(`"%v"`, item))
	}

	text, err := pathx.LoadTemplate(category, serverTemplateFile, serverTemplate)
	if err != nil {
		return err
	}

	notStream := false
	for _, rpc := range service.RPC {
		if !rpc.StreamsRequest && !rpc.StreamsReturns {
			notStream = true
			break
		}
	}

	return util.With("server").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
		"head": head,
		"unimplementedServer": fmt.Sprintf("%s.Unimplemented%sServer", proto.PbPackage,
			stringx.From(service.Name).ToCamel()),
		"server":    stringx.From(service.Name).ToCamel(),
		"imports":   strings.Join(imports.KeysStr(), pathx.NL),
		"funcs":     strings.Join(funcList, pathx.NL),
		"notStream": notStream,
	}, serverFile, true)
}

func (g *Generator) genFunctions(goPackage string, service parser.Service, multiple bool, typeMap map[string]string, pbImport string) ([]string, []string, error) {
	var (
		functionList []string
		logicPkg     string
		impList      []string
	)
	for _, rpc := range service.RPC {
		text, err := pathx.LoadTemplate(category, serverFuncTemplateFile, functionTemplate)
		if err != nil {
			return nil, impList, err
		}

		//var logicName string
		if !multiple {
			logicPkg = "core"
			//logicName = fmt.Sprintf("%s", stringx.From(rpc.Name).ToCamel())
		} else {
			logicPkg = "core"
			//nameJoin := fmt.Sprintf("%s_logic", service.Name)
			//logicPkg = strings.ToLower(stringx.From(nameJoin).ToCamel())
			//logicName = fmt.Sprintf("%s", stringx.From(rpc.Name).ToCamel())
		}

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
			return nil, nil, err
		}

		response := func() string {
			var mess = rpc.ReturnsType
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
					err = errors.New(fmt.Sprintf("request package %s must defined in flags type_map", pkgKey))
					return ""
				}
			} else if strings.HasPrefix(mess, goPackage+".") {
				mess = strings.Split(mess, ".")[1]
			}
			impList = append(impList, pbImport)
			return fmt.Sprintf("*%s.%s", goPackage, parser.CamelCase(mess))
		}()
		if err != nil {
			return nil, nil, err
		}

		comment := parser.GetComment(rpc.Doc())
		streamServer := fmt.Sprintf("%s.%s_%s%s", goPackage, parser.CamelCase(service.Name),
			parser.CamelCase(rpc.Name), "Server")
		buffer, err := util.With("func").Parse(text).Execute(map[string]interface{}{
			"server":     stringx.From(service.Name).ToCamel(),
			"logicName":  stringx.From(goPackage).ToCamel(),
			"method":     parser.CamelCase(rpc.Name),
			"handler":    goPackage + "." + rpc.Name,
			"request":    request,
			"response":   response,
			"hasComment": len(comment) > 0,
			"comment":    comment,
			"hasReq":     !rpc.StreamsRequest,
			"stream":     rpc.StreamsRequest || rpc.StreamsReturns,
			"notStream":  !rpc.StreamsRequest && !rpc.StreamsReturns,
			"streamBody": streamServer,
			"logicPkg":   logicPkg,
		})
		if err != nil {
			return nil, impList, err
		}

		functionList = append(functionList, buffer.String())
	}
	return functionList, impList, nil
}
