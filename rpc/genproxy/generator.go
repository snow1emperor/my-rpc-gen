package genproxy

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/snow1emperor/my-rpc-gen/rpc/parser"
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"path/filepath"
	"sort"
	"strings"
)

const (
	PathTemplate        = `	%s: {"%s", Fn(new(%s)), Fn(new(%s))},`
	ConstructorTemplate = `	%s = %d;`

	constructorTemplateFile       = "constructor.tpl"
	proxyTemplateFile             = "proxy.tpl"
	initialConstructor      int64 = 100100101
)

//go:embed constructor.tpl
var constructorTemplate string

//go:embed proxy.tpl
var proxyTemplate string

type ProxyGenerator struct {
}

func NewProxyGenerator() *ProxyGenerator {
	return &ProxyGenerator{}
}

func (gen ProxyGenerator) Generate(ctx *ProxyContext) error {
	var (
		constructor  = initialConstructor
		regFile      = filepath.Join(ctx.Dst, "proxy_register.go")
		conFile      = filepath.Join(ctx.Dst, "proxy_constructors.proto")
		head         = util.GetHead("proxy_constructors.proto")
		protoArr     []parser.Proto
		p            = parser.NewDefaultProtoParser()
		constructors = collection.NewSet()
		paths        = collection.NewSet()
		imports      = collection.NewSet()
		filePackage  string
	)
	filePackage = func() string {
		if ctx.Pkg != "" {
			return ctx.Pkg
		}
		var path = ctx.Dst
		if strings.Contains(path, "/") {
			var values = strings.Split(path, "/")
			return values[len(values)-1]
		}
		return path
	}()
	//s, _ := jsonx.MarshalToString(ctx)
	//fmt.Println(s)

	abs, err := filepath.Abs(ctx.Dst)
	if err != nil {
		return err
	}

	err = pathx.MkdirIfNotExist(abs)
	if err != nil {
		return err
	}

	for _, fl := range ctx.Out {
		pt, err := p.Parse(fl, false)
		if err != nil {
			return err
		}
		protoArr = append(protoArr, pt)
	}

	for _, pt := range protoArr {
		val, ok := ctx.TypeMap[pt.Package.Package.Name]
		if !ok {
			fmt.Println("can`t parse ", pt.Name, " , skip value")
			continue
		}
		for _, srv := range pt.Service {
			for _, item := range srv.RPC {
				if item.StreamsReturns || item.StreamsRequest {
					fmt.Println("stream not support ", item.Name, " , skip rpc")
					continue
				}

				var Fn = func(mess string) string {
					if strings.Contains(mess, "google.protobuf") {
						if path, ok := ctx.TypeMap["types"]; ok {
							imports.AddStr(fmt.Sprintf(`	"%s"`, path))
						}
						values := strings.Split(mess, ".")
						return fmt.Sprintf("%s.%s", "types", values[len(values)-1])
					} else if strings.Contains(mess, ".") && !strings.HasPrefix(mess, pt.PbPackage+".") {
						pkgKey := strings.Split(mess, ".")[0]
						if path, ok := ctx.TypeMap[pkgKey]; ok {
							imports.AddStr(fmt.Sprintf(`	"%s"`, path))
							return fmt.Sprintf("%s", mess)
						} else {
							err = errors.New(fmt.Sprintf("request type %s must defined in flags type_map", pkgKey))
							return ""
						}
					} else if strings.HasPrefix(mess, pt.PbPackage+".") {
						imports.AddStr(fmt.Sprintf(`	"%s"`, val))
						return mess
					}
					imports.AddStr(fmt.Sprintf(`	"%s"`, val))
					return fmt.Sprintf("%s.%s", pt.PbPackage, parser.CamelCase(mess))
				}

				request := Fn(item.RequestType)
				response := Fn(item.ReturnsType)

				if request == "" || response == "" {
					fmt.Println("can`t parse: ", item.RequestType, " or ", item.ReturnsType, " - skip value")
					continue
				}

				paths.AddStr(fmt.Sprintf(PathTemplate,
					fmt.Sprintf("TLConstructor_%s_%s", pt.PbPackage, item.Name),
					fmt.Sprintf("/%s.%s/%s", pt.PbPackage, srv.Name, item.Name),
					request, response,
				))
				constructors.AddStr(fmt.Sprintf(ConstructorTemplate, fmt.Sprintf("%s_%s", pt.PbPackage, item.Name), constructor))
				constructor++
			}
		}
	}

	text, err := pathx.LoadTemplate("rpc", proxyTemplateFile, proxyTemplate)
	if err != nil {
		return err
	}
	importsSrt := imports.KeysStr()
	sort.Strings(importsSrt)
	pathsSrt := paths.KeysStr()
	sort.Strings(pathsSrt)
	if err = util.With("proxy").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
		"head":        head,
		"imports":     strings.Join(importsSrt, pathx.NL),
		"paths":       strings.Join(pathsSrt, pathx.NL),
		"packageName": filePackage,
	}, regFile, true); err != nil {
		return err
	}

	text, err = pathx.LoadTemplate("rpc", constructorTemplateFile, constructorTemplate)
	if err != nil {
		return err
	}
	constructorsSrt := constructors.KeysStr()
	sort.Strings(constructorsSrt)
	if err = util.With("proxy").GoFmt(false).Parse(text).SaveTo(map[string]interface{}{
		"head":         head,
		"constructors": strings.Join(constructorsSrt, pathx.NL),
		"packageName":  filePackage,
	}, conFile, true); err != nil {
		return err
	}

	return nil
}
