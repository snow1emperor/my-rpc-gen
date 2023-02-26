package generator

import (
	_ "embed"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/emicklei/proto"
	"github.com/snow1emperor/my-rpc-gen/rpc/parser"
	"github.com/zeromicro/go-zero/core/collection"
	conf "github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"github.com/zeromicro/go-zero/tools/goctl/util/stringx"
)

const (
	callInterfaceFunctionTemplate = `{{if .hasComment}}{{.comment}}
{{end}}{{.method}}(ctx context.Context{{if .hasReq}}, in *{{.pbRequest}}{{end}}, opts ...grpc.CallOption) ({{if .notStream}}*{{.pbResponse}}, {{else}}{{.streamBody}},{{end}} error)`

	callFunctionTemplate = `
{{if .hasComment}}{{.comment}}{{end}}
func (m *default{{.serviceName}}Client) {{.method}}(ctx context.Context{{if .hasReq}}, in *{{.pbRequest}}{{end}}, opts ...grpc.CallOption) ({{if .notStream}}*{{.pbResponse}}, {{else}}{{.streamBody}},{{end}} error) {
	client := {{if .isCallPkgSameToGrpcPkg}}{{else}}{{.package}}.{{end}}New{{.rpcServiceName}}Client(m.cli.Conn())
	return client.{{.method}}(ctx{{if .hasReq}}, in{{end}}, opts...)
}
`
)

//go:embed call.tpl
var callTemplateText string

// GenCall generates the rpc client code, which is the entry point for the rpc service call.
// It is a layer of encapsulation for the rpc client and shields the details in the pb.
func (g *Generator) GenCall(ctx DirContext, proto parser.Proto, cfg *conf.Config,
	c *ZRpcContext) error {
	if !c.Multiple {
		return g.genCallInCompatibility(ctx, proto, cfg, c)
	}

	return g.genCallGroup(ctx, proto, cfg, c)
}

func (g *Generator) genCallGroup(ctx DirContext, proto parser.Proto, cfg *conf.Config, c *ZRpcContext) error {
	dir := ctx.GetCall()
	head := util.GetHead(proto.Name)
	for _, service := range proto.Service {
		childPkg, err := dir.GetChildPackage(service.Name)
		if err != nil {
			return err
		}

		callFilename, err := format.FileNamingFormat(cfg.NamingFormat, service.Name)
		if err != nil {
			return err
		}

		childDir := filepath.Base(childPkg)
		filename := filepath.Join(dir.Filename, childDir, fmt.Sprintf("%s.go", callFilename))
		isCallPkgSameToPbPkg := childDir == ctx.GetProtoGo().Filename
		isCallPkgSameToGrpcPkg := childDir == ctx.GetProtoGo().Filename

		functions, err := g.genFunction(proto.PbPackage, service, isCallPkgSameToGrpcPkg, c.VarStringTypeMap)
		if err != nil {
			return err
		}

		iFunctions, err := g.getInterfaceFuncs(proto.PbPackage, service, isCallPkgSameToGrpcPkg, c.VarStringTypeMap)
		if err != nil {
			return err
		}

		text, err := pathx.LoadTemplate(category, callTemplateFile, callTemplateText)
		if err != nil {
			return err
		}

		imports := collection.NewSet()
		alias := collection.NewSet()
		if !isCallPkgSameToPbPkg {
			var impList []string
			var typeMap = c.VarStringTypeMap
			for _, rpc := range proto.Service[0].RPC {

				request := func() string {
					var mess = rpc.RequestType
					if strings.Contains(mess, "google.protobuf") {
						if path, ok := typeMap["types"]; ok {
							impList = append(impList, path)
						}
						typeName := parser.CamelCase(strings.Trim(mess, "google.protobuf."))
						return fmt.Sprintf("%s = %s", typeName,
							fmt.Sprintf("%s.%s", "types", typeName))
					} else if strings.Contains(mess, ".") && !strings.HasPrefix(mess, proto.PbPackage+".") {
						values := strings.Split(mess, ".")
						pkgKey, typeName := values[0], values[1]
						if path, ok := typeMap[pkgKey]; ok {
							impList = append(impList, path)
							return fmt.Sprintf("%s = %s", typeName, mess)
						} else {
							err = errors.New(fmt.Sprintf("request type %s must defined in flags type_map", pkgKey))
							return ""
						}
					} else if strings.HasPrefix(mess, proto.PbPackage+".") {
						mess = strings.Split(mess, ".")[1]
					}
					typeName := fmt.Sprintf("%s.%s", proto.PbPackage, parser.CamelCase(mess))
					return fmt.Sprintf("%s = %s", parser.CamelCase(mess), typeName)
				}()
				if err != nil {
					return err
				}
				alias.AddStr(request)

				response := func() string {
					var mess = rpc.ReturnsType
					if strings.Contains(mess, "google.protobuf") {
						if path, ok := typeMap["types"]; ok {
							impList = append(impList, path)
						}
						typeName := parser.CamelCase(strings.Trim(mess, "google.protobuf."))
						return fmt.Sprintf("%s = %s", typeName,
							fmt.Sprintf("%s.%s", "types", typeName))
					} else if strings.Contains(mess, ".") && !strings.HasPrefix(mess, proto.PbPackage+".") {
						values := strings.Split(mess, ".")
						pkgKey, typeName := values[0], values[1]
						if path, ok := typeMap[pkgKey]; ok {
							impList = append(impList, path)
							return fmt.Sprintf("%s = %s", typeName, mess)
						} else {
							err = errors.New(fmt.Sprintf("request type %s must defined in flags type_map", pkgKey))
							return ""
						}
					} else if strings.HasPrefix(mess, proto.PbPackage+".") {
						mess = strings.Split(mess, ".")[1]
					}
					typeName := fmt.Sprintf("%s.%s", proto.PbPackage, parser.CamelCase(mess))
					return fmt.Sprintf("%s = %s", parser.CamelCase(mess), typeName)
				}()
				if err != nil {
					return err
				}
				alias.AddStr(response)
				//msgName := getMessageName(*item.Message)
			}
			for _, item := range impList {
				imports.AddStr(fmt.Sprintf(`"%v"`, item))
			}
		}

		pbPackage := fmt.Sprintf(`"%s"`, ctx.GetPb().Package)
		protoGoPackage := fmt.Sprintf(`"%s"`, ctx.GetProtoGo().Package)
		if isCallPkgSameToGrpcPkg {
			pbPackage = ""
			protoGoPackage = ""
		}

		aliasKeys := alias.KeysStr()
		sort.Strings(aliasKeys)
		if err = util.With("shared").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
			"name":           callFilename,
			"alias":          strings.Join(aliasKeys, pathx.NL),
			"head":           head,
			"filePackage":    dir.Base,
			"pbPackage":      pbPackage,
			"protoGoPackage": protoGoPackage,
			"serviceName":    stringx.From(service.Name).ToCamel(),
			"functions":      strings.Join(functions, pathx.NL),
			"interface":      strings.Join(iFunctions, pathx.NL),
			"imports":        strings.Join(imports.KeysStr(), pathx.NL),
		}, filename, true); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) genCallInCompatibility(ctx DirContext, proto parser.Proto, cfg *conf.Config, c *ZRpcContext) error {
	dir := ctx.GetCall()
	service := proto.Service[0]
	head := util.GetHead(proto.Name)
	isCallPkgSameToPbPkg := ctx.GetCall().Filename == ctx.GetPb().Filename
	isCallPkgSameToGrpcPkg := ctx.GetCall().Filename == ctx.GetProtoGo().Filename

	callFilename := ctx.GetServiceName().Lower() + "_client"

	filename := filepath.Join(dir.Filename, fmt.Sprintf("%s.go", callFilename))
	functions, err := g.genFunction(proto.PbPackage, service, isCallPkgSameToGrpcPkg, c.VarStringTypeMap)
	if err != nil {
		return err
	}

	iFunctions, err := g.getInterfaceFuncs(proto.PbPackage, service, isCallPkgSameToGrpcPkg, c.VarStringTypeMap)
	if err != nil {
		return err
	}

	text, err := pathx.LoadTemplate(category, callTemplateFile, callTemplateText)
	if err != nil {
		return err
	}

	imports := collection.NewSet()
	alias := collection.NewSet()
	if !isCallPkgSameToPbPkg {
		var impList []string
		var typeMap = c.VarStringTypeMap
		for _, rpc := range proto.Service[0].RPC {

			request := func() string {
				var mess = rpc.RequestType
				if strings.Contains(mess, "google.protobuf") {
					if path, ok := typeMap["types"]; ok {
						impList = append(impList, path)
					}
					typeName := parser.CamelCase(strings.Trim(mess, "google.protobuf."))
					return fmt.Sprintf("%s = %s", typeName,
						fmt.Sprintf("%s.%s", "types", typeName))
				} else if strings.Contains(mess, ".") && !strings.HasPrefix(mess, proto.PbPackage+".") {
					values := strings.Split(mess, ".")
					pkgKey, typeName := values[0], values[1]
					if path, ok := typeMap[pkgKey]; ok {
						impList = append(impList, path)
						return fmt.Sprintf("%s = %s", typeName, mess)
					} else {
						err = errors.New(fmt.Sprintf("request type %s must defined in flags type_map", pkgKey))
						return ""
					}
				} else if strings.HasPrefix(mess, proto.PbPackage+".") {
					mess = strings.Split(mess, ".")[1]
				}
				typeName := fmt.Sprintf("%s.%s", proto.PbPackage, parser.CamelCase(mess))
				return fmt.Sprintf("%s = %s", parser.CamelCase(mess), typeName)
			}()
			if err != nil {
				return err
			}
			alias.AddStr(request)

			response := func() string {
				var mess = rpc.ReturnsType
				if strings.Contains(mess, "google.protobuf") {
					if path, ok := typeMap["types"]; ok {
						impList = append(impList, path)
					}
					typeName := parser.CamelCase(strings.Trim(mess, "google.protobuf."))
					return fmt.Sprintf("%s = %s", typeName,
						fmt.Sprintf("%s.%s", "types", typeName))
				} else if strings.Contains(mess, ".") && !strings.HasPrefix(mess, proto.PbPackage+".") {
					values := strings.Split(mess, ".")
					pkgKey, typeName := values[0], values[1]
					if path, ok := typeMap[pkgKey]; ok {
						impList = append(impList, path)
						return fmt.Sprintf("%s = %s", typeName, mess)
					} else {
						err = errors.New(fmt.Sprintf("request type %s must defined in flags type_map", pkgKey))
						return ""
					}
				} else if strings.HasPrefix(mess, proto.PbPackage+".") {
					mess = strings.Split(mess, ".")[1]
				}
				typeName := fmt.Sprintf("%s.%s", proto.PbPackage, parser.CamelCase(mess))
				return fmt.Sprintf("%s = %s", parser.CamelCase(mess), typeName)
			}()
			if err != nil {
				return err
			}
			alias.AddStr(response)
			//msgName := getMessageName(*item.Message)
		}
		for _, item := range impList {
			imports.AddStr(fmt.Sprintf(`"%v"`, item))
		}
	}

	pbPackage := fmt.Sprintf(`"%s"`, ctx.GetPb().Package)
	protoGoPackage := fmt.Sprintf(`"%s"`, ctx.GetProtoGo().Package)
	if isCallPkgSameToGrpcPkg {
		pbPackage = ""
		protoGoPackage = ""
	}
	aliasKeys := alias.KeysStr()
	sort.Strings(aliasKeys)
	return util.With("shared").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
		"name":           callFilename,
		"alias":          strings.Join(aliasKeys, pathx.NL),
		"head":           head,
		"filePackage":    ctx.GetServiceName().Lower() + "_client",
		"pbPackage":      pbPackage,
		"protoGoPackage": protoGoPackage,
		"serviceName":    ctx.GetServiceName().Title(),
		"functions":      strings.Join(functions, pathx.NL),
		"interface":      strings.Join(iFunctions, pathx.NL),
		"imports":        strings.Join(imports.KeysStr(), pathx.NL),
	}, filename, true)
}

func getMessageName(msg proto.Message) string {
	list := []string{msg.Name}

	for {
		parent := msg.Parent
		if parent == nil {
			break
		}

		parentMsg, ok := parent.(*proto.Message)
		if !ok {
			break
		}

		tmp := []string{parentMsg.Name}
		list = append(tmp, list...)
		msg = *parentMsg
	}

	return strings.Join(list, "_")
}

func (g *Generator) genFunction(goPackage string, service parser.Service, isCallPkgSameToGrpcPkg bool, typeMap map[string]string) ([]string, error) {
	functions := make([]string, 0)

	for _, rpc := range service.RPC {
		text, err := pathx.LoadTemplate(category, callFunctionTemplateFile, callFunctionTemplate)
		if err != nil {
			return nil, err
		}

		comment := parser.GetComment(rpc.Doc())
		streamServer := fmt.Sprintf("%s.%s_%s%s", goPackage, parser.CamelCase(service.Name),
			parser.CamelCase(rpc.Name), "Client")
		if isCallPkgSameToGrpcPkg {
			streamServer = fmt.Sprintf("%s_%s%s", parser.CamelCase(service.Name),
				parser.CamelCase(rpc.Name), "Client")
		}

		request := func() string {
			var mess = rpc.RequestType
			if strings.Contains(mess, "google.protobuf") {
				return parser.CamelCase(strings.Trim(mess, "google.protobuf."))
			} else if strings.Contains(mess, ".") && !strings.HasPrefix(mess, goPackage+".") {
				pkgKey := strings.Split(mess, ".")[0]
				if _, ok := typeMap[pkgKey]; ok {
					return parser.CamelCase(strings.Trim(mess, pkgKey+"."))
				} else {
					err = errors.New(fmt.Sprintf("request type %s must defined in flags type_map", pkgKey))
					return ""
				}
			} else if strings.HasPrefix(mess, goPackage+".") {
				mess = strings.Split(mess, ".")[1]
			}
			return parser.CamelCase(mess)
		}()
		if err != nil {
			return nil, err
		}

		response := func() string {
			var mess = rpc.ReturnsType
			if strings.Contains(mess, "google.protobuf") {
				return parser.CamelCase(strings.Trim(mess, "google.protobuf."))
			} else if strings.Contains(mess, ".") && !strings.HasPrefix(mess, goPackage+".") {
				pkgKey := strings.Split(mess, ".")[0]
				if _, ok := typeMap[pkgKey]; ok {
					return parser.CamelCase(strings.Trim(mess, pkgKey+"."))
				} else {
					err = errors.New(fmt.Sprintf("request type %s must defined in flags type_map", pkgKey))
					return ""
				}
			} else if strings.HasPrefix(mess, goPackage+".") {
				mess = strings.Split(mess, ".")[1]
			}
			return parser.CamelCase(mess)
		}()
		if err != nil {
			return nil, err
		}

		buffer, err := util.With("sharedFn").Parse(text).Execute(map[string]interface{}{
			"serviceName":            stringx.From(goPackage).ToCamel(),
			"rpcServiceName":         parser.CamelCase(service.Name),
			"method":                 parser.CamelCase(rpc.Name),
			"package":                goPackage,
			"pbRequest":              request,
			"pbResponse":             response,
			"hasComment":             len(comment) > 0,
			"comment":                comment,
			"hasReq":                 !rpc.StreamsRequest,
			"notStream":              !rpc.StreamsRequest && !rpc.StreamsReturns,
			"streamBody":             streamServer,
			"isCallPkgSameToGrpcPkg": isCallPkgSameToGrpcPkg,
		})
		if err != nil {
			return nil, err
		}

		functions = append(functions, buffer.String())
	}

	return functions, nil
}

func (g *Generator) getInterfaceFuncs(goPackage string, service parser.Service, isCallPkgSameToGrpcPkg bool, typeMap map[string]string) ([]string, error) {
	functions := make([]string, 0)

	for _, rpc := range service.RPC {
		text, err := pathx.LoadTemplate(category, callInterfaceFunctionTemplateFile,
			callInterfaceFunctionTemplate)
		if err != nil {
			return nil, err
		}

		comment := parser.GetComment(rpc.Doc())
		streamServer := fmt.Sprintf("%s.%s_%s%s", goPackage, parser.CamelCase(service.Name),
			parser.CamelCase(rpc.Name), "Client")
		if isCallPkgSameToGrpcPkg {
			streamServer = fmt.Sprintf("%s_%s%s", parser.CamelCase(service.Name),
				parser.CamelCase(rpc.Name), "Client")
		}
		request := func() string {
			var mess = rpc.RequestType
			if strings.Contains(mess, "google.protobuf") {
				return parser.CamelCase(strings.Trim(mess, "google.protobuf."))
			} else if strings.Contains(mess, ".") && !strings.HasPrefix(mess, goPackage+".") {
				pkgKey := strings.Split(mess, ".")[0]
				if _, ok := typeMap[pkgKey]; ok {
					return parser.CamelCase(strings.Trim(mess, pkgKey+"."))
				} else {
					err = errors.New(fmt.Sprintf("request type %s must defined in flags type_map", pkgKey))
					return ""
				}
			} else if strings.HasPrefix(mess, goPackage+".") {
				mess = strings.Split(mess, ".")[1]
			}
			return parser.CamelCase(mess)
		}()
		if err != nil {
			return nil, err
		}

		response := func() string {
			var mess = rpc.ReturnsType
			if strings.Contains(mess, "google.protobuf") {
				return parser.CamelCase(strings.Trim(mess, "google.protobuf."))
			} else if strings.Contains(mess, ".") && !strings.HasPrefix(mess, goPackage+".") {
				pkgKey := strings.Split(mess, ".")[0]
				if _, ok := typeMap[pkgKey]; ok {
					return parser.CamelCase(strings.Trim(mess, pkgKey+"."))
				} else {
					err = errors.New(fmt.Sprintf("request type %s must defined in flags type_map", pkgKey))
					return ""
				}
			} else if strings.HasPrefix(mess, goPackage+".") {
				mess = strings.Split(mess, ".")[1]
			}
			return parser.CamelCase(mess)
		}()
		if err != nil {
			return nil, err
		}

		buffer, err := util.With("interfaceFn").Parse(text).Execute(
			map[string]interface{}{
				"hasComment": len(comment) > 0,
				"comment":    comment,
				"method":     parser.CamelCase(rpc.Name),
				"hasReq":     !rpc.StreamsRequest,
				"pbRequest":  request,
				"notStream":  !rpc.StreamsRequest && !rpc.StreamsReturns,
				"pbResponse": response,
				"streamBody": streamServer,
			})
		if err != nil {
			return nil, err
		}

		functions = append(functions, buffer.String())
	}

	return functions, nil
}
