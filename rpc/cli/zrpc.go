package cli

import (
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"my-rpc-gen/rpc/generator"
)

var (
	errInvalidGrpcOutput = errors.New("ZRPC: missing --go-grpc_out")
	errInvalidGoOutput   = errors.New("ZRPC: missing --go_out")
	errInvalidZrpcOutput = errors.New("ZRPC: missing zrpc output, please use --zrpc_out to specify the output")
)

// ZRPC generates grpc code directly by protoc and generates
// zrpc code by goctl.
func ZRPC(_ *cobra.Command, args []string) error {
	if err, ok := gogoGen(args); ok {
		return err
	}
	protocArgs := wrapProtocCmd("protoc", args)
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	source := args[0]
	grpcOutList := VarStringSliceGoGRPCOut
	goOutList := VarStringSliceGoOut
	zrpcOut := VarStringZRPCOut
	style := VarStringStyle
	home := VarStringHome
	remote := VarStringRemote
	branch := VarStringBranch
	verbose := VarBoolVerbose

	if len(grpcOutList) == 0 {
		return errInvalidGrpcOutput
	}
	if len(goOutList) == 0 {
		return errInvalidGoOutput
	}
	goOut := goOutList[len(goOutList)-1]
	grpcOut := grpcOutList[len(grpcOutList)-1]
	if len(goOut) == 0 {
		return errInvalidGrpcOutput
	}
	if len(zrpcOut) == 0 {
		return errInvalidZrpcOutput
	}
	goOutAbs, err := filepath.Abs(goOut)
	if err != nil {
		return err
	}
	grpcOutAbs, err := filepath.Abs(grpcOut)
	if err != nil {
		return err
	}
	err = pathx.MkdirIfNotExist(goOutAbs)
	if err != nil {
		return err
	}
	err = pathx.MkdirIfNotExist(grpcOutAbs)
	if err != nil {
		return err
	}
	if len(remote) > 0 {
		repo, _ := util.CloneIntoGitHome(remote, branch)
		if len(repo) > 0 {
			home = repo
		}
	}

	if len(home) > 0 {
		pathx.RegisterGoctlHome(home)
	}
	if !filepath.IsAbs(zrpcOut) {
		zrpcOut = filepath.Join(pwd, zrpcOut)
	}

	isGooglePlugin := len(grpcOut) > 0
	goOut, err = filepath.Abs(goOut)
	if err != nil {
		return err
	}
	grpcOut, err = filepath.Abs(grpcOut)
	if err != nil {
		return err
	}
	zrpcOut, err = filepath.Abs(zrpcOut)
	if err != nil {
		return err
	}

	var ctx generator.ZRpcContext
	ctx.Multiple = VarBoolMultiple
	ctx.Src = source
	ctx.GoOutput = goOut
	ctx.GrpcOutput = grpcOut
	ctx.IsGooglePlugin = isGooglePlugin

	ctx.Output = zrpcOut
	ctx.ProtocCmd = strings.Join(protocArgs, " ")
	g := generator.NewGenerator(style, verbose)
	if err := g.Generate(&ctx); err != nil {
		return err
	}
	return nil
}

func gogoGen(args []string) (error, bool) {
	if len(VarStringSliceGogoOut) > 0 {
		protocArgs := wrapProtocCmd("protoc", args)
		source := args[0]
		zrpcOut := VarStringZRPCOut
		style := VarStringStyle
		verbose := VarBoolVerbose
		var ctx generator.ZRpcContext
		ctx.Multiple = VarBoolMultiple
		ctx.Src = source

		ctx.GoOutput = VarStringGogoDst
		ctx.GrpcOutput = VarStringGogoDst
		ctx.ProtoGenGoDir = VarStringGogoDst
		ctx.ProtoGenGrpcDir = VarStringGogoDst
		ctx.VarStringCommandsPkg = VarStringCommandsPkg
		ctx.IsGogoPlugin = true
		ctx.Output = zrpcOut
		ctx.ProtocCmd = strings.Join(protocArgs, " ")
		g := generator.NewGenerator(style, verbose)
		//log.Printf("zrpc: %+v\n", ctx)
		if err := g.Generate(&ctx); err != nil {
			return err, true
		}
		return nil, true
	}
	return nil, false
}

func wrapProtocCmd(name string, args []string) []string {
	ret := append([]string{name}, args...)
	for _, protoPath := range VarStringSliceProtoPath {
		ret = append(ret, "--proto_path", protoPath)
	}
	for _, goOpt := range VarStringSliceGoOpt {
		ret = append(ret, "--go_opt", goOpt)
	}
	for _, goGRPCOpt := range VarStringSliceGoGRPCOpt {
		ret = append(ret, "--go-grpc_opt", goGRPCOpt)
	}
	for _, goOut := range VarStringSliceGoOut {
		ret = append(ret, "--go_out", goOut)
	}
	if len(VarStringSliceGogoOut) > 0 {
		ret = append(ret, fmt.Sprintf("--gogo_out=%s,:%s", strings.Join(VarStringSliceGogoOut, ","), VarStringGogoDst))
	}
	for _, goGRPCOut := range VarStringSliceGoGRPCOut {
		ret = append(ret, "--go-grpc_out", goGRPCOut)
	}
	for _, plugin := range VarStringSlicePlugin {
		ret = append(ret, "--plugin="+plugin)
	}
	return ret
}
