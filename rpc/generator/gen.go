package generator

import (
	"path/filepath"

	"github.com/snow1emperor/my-rpc-gen/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/ctx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

type ZRpcContext struct {
	// Sre is the source file of the proto.
	Src string
	// ProtoCmd is the command to generate proto files.
	ProtocCmd string
	// ProtoGenGrpcDir is the directory to store the generated proto files.
	ProtoGenGrpcDir string
	// ProtoGenGoDir is the directory to store the generated go files.
	ProtoGenGoDir string
	// IsGooglePlugin is the flag to indicate whether the proto file is generated by google plugin.
	IsGooglePlugin bool
	// IsGooglePlugin is the flag to indicate whether the proto file is generated by google plugin.
	IsGogoPlugin bool
	// GoOutput is the output directory of the generated go files.
	GoOutput string
	// GrpcOutput is the output directory of the generated grpc files.
	GrpcOutput string
	// Output is the output directory of the generated files.
	Output string
	// Multiple is the flag to indicate whether the proto file is generated in multiple mode.
	Multiple bool
	// VarStringCommandsPkg
	VarStringCommandsPkg string
	// VarStringMTProtPkg
	VarStringMTProtPkg string
}

// Generate generates a rpc service, through the proto file,
// code storage directory, and proto import parameters to control
// the source file and target location of the rpc service that needs to be generated
func (g *Generator) Generate(zctx *ZRpcContext) error {
	// log.Printf("%+v\n", zctx)
	abs, err := filepath.Abs(zctx.Output)
	if err != nil {
		return err
	}

	err = pathx.MkdirIfNotExist(abs)
	if err != nil {
		return err
	}

	err = g.Prepare()
	if err != nil {
		return err
	}

	projectCtx, err := ctx.Prepare(abs)
	if err != nil {
		return err
	}

	p := parser.NewDefaultProtoParser()
	proto, err := p.Parse(zctx.Src, zctx.Multiple)
	if err != nil {
		return err
	}

	dirCtx, err := mkdir(projectCtx, proto, g.cfg, zctx)
	if err != nil {
		return err
	}

	err = g.GenEtc(dirCtx, proto, g.cfg)
	if err != nil {
		return err
	}

	err = g.GenPb(dirCtx, zctx)
	if err != nil {
		return err
	}

	err = g.GenConfig(dirCtx, proto, g.cfg)
	if err != nil {
		return err
	}

	err = g.GenSvc(dirCtx, proto, g.cfg)
	if err != nil {
		return err
	}

	err = g.GenLogic(dirCtx, proto, g.cfg, zctx)
	if err != nil {
		return err
	}

	err = g.GenServer(dirCtx, proto, g.cfg, zctx)
	if err != nil {
		return err
	}

	err = g.GenMain(dirCtx, proto, g.cfg, zctx)
	if err != nil {
		return err
	}

	err = g.GenCall(dirCtx, proto, g.cfg, zctx)

	console.NewColorConsole().MarkDone()

	return err
}
