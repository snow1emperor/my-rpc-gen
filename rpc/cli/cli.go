package cli

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/snow1emperor/my-rpc-gen/rpc/generator"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var (
	// VarStringOutput describes the output.
	VarStringOutput string
	// VarStringHome describes the goctl home.
	VarStringHome string
	// VarStringRemote describes the remote git repository.
	VarStringRemote string
	// VarStringBranch describes the git branch.
	VarStringBranch string
	// VarStringGogoDst describes the git branch.
	VarStringGogoDst string
	// VarStringCommandsPkg describes the git branch.
	VarStringCommandsPkg string
	// VarStringTypes describes the git branch.
	VarStringTypes string
	// VarStringTypeMap describes the git branch.
	VarStringTypeMap string
	// VarStringSliceGoOut describes the go output.
	VarStringSliceGoOut []string
	// VarStringSliceGogoOut describes the go output.
	VarStringSliceGogoOut []string
	// VarStringSliceGoGRPCOut describes the grpc output.
	VarStringSliceGoGRPCOut []string
	// VarStringSlicePlugin describes the protoc plugin.
	VarStringSlicePlugin []string
	// VarStringSliceProtoPath describes the proto path.
	VarStringSliceProtoPath []string
	// VarStringSliceGoOpt describes the go options.
	VarStringSliceGoOpt []string
	// VarStringSliceGoGRPCOpt describes the grpc options.
	VarStringSliceGoGRPCOpt []string
	// VarStringStyle describes the style of output files.
	VarStringStyle string
	// VarStringZRPCOut describes the zRPC output.
	VarStringZRPCOut string
	// VarBoolIdea describes whether idea or not
	VarBoolIdea bool
	// VarBoolVerbose describes whether verbose.
	VarBoolVerbose bool
	// VarBoolMultiple describes whether support generating multiple rpc services or not.
	VarBoolMultiple bool
)

// RPCNew is to generate rpc greet service, this greet service can speed
// up your understanding of the zrpc service structure
func RPCNew(_ *cobra.Command, args []string) error {
	rpcname := args[0]
	ext := filepath.Ext(rpcname)
	if len(ext) > 0 {
		return fmt.Errorf("unexpected ext: %s", ext)
	}
	style := VarStringStyle
	home := VarStringHome
	remote := VarStringRemote
	branch := VarStringBranch
	verbose := VarBoolVerbose
	if len(remote) > 0 {
		repo, _ := util.CloneIntoGitHome(remote, branch)
		if len(repo) > 0 {
			home = repo
		}
	}
	if len(home) > 0 {
		pathx.RegisterGoctlHome(home)
	}

	protoName := rpcname + ".proto"
	filename := filepath.Join(".", rpcname, protoName)
	src, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	err = generator.ProtoTmpl(src)
	if err != nil {
		return err
	}

	var ctx generator.ZRpcContext
	ctx.Src = src
	ctx.VarStringCommandsPkg = VarStringCommandsPkg
	ctx.VarStringTypeMap = make(map[string]string)
	if len(VarStringTypes) != 0 {
		ctx.VarStringTypeMap["types"] = VarStringTypes
	}
	if len(VarStringTypeMap) != 0 {
		valueList := strings.Split(VarStringTypeMap, ",")
		for _, val := range valueList {
			if len(val) != 0 && strings.Contains(val, "=") {
				kv := strings.Split(val, "=")
				key, val := kv[0], kv[1]
				if len(key) != 0 && len(val) != 0 {
					ctx.VarStringTypeMap[key] = val
				}
			}
		}
	}
	ctx.GoOutput = filepath.Dir(src)
	ctx.GrpcOutput = filepath.Dir(src)
	ctx.IsGooglePlugin = true
	ctx.Output = filepath.Dir(src)
	ctx.ProtocCmd = fmt.Sprintf("protoc -I=%s %s --go_out=%s --go-grpc_out=%s", filepath.Dir(src), filepath.Base(src), filepath.Dir(src), filepath.Dir(src))

	grpcOptList := VarStringSliceGoGRPCOpt
	if len(grpcOptList) > 0 {
		ctx.ProtocCmd += " --go-grpc_opt=" + strings.Join(grpcOptList, ",")
	}

	goOptList := VarStringSliceGoOpt
	if len(goOptList) > 0 {
		ctx.ProtocCmd += " --go_opt=" + strings.Join(goOptList, ",")
	}

	g := generator.NewGenerator(style, verbose)
	//log.Printf("cli: %+v\n", ctx)
	return g.Generate(&ctx)
}

// RPCTemplate is the entry for generate rpc template
func RPCTemplate(latest bool) error {
	if !latest {
		console.Warning("deprecated: goctl rpc template -o is deprecated and will be removed in the future, use goctl rpc -o instead")
	}
	protoFile := VarStringOutput
	home := VarStringHome
	remote := VarStringRemote
	branch := VarStringBranch
	if len(remote) > 0 {
		repo, _ := util.CloneIntoGitHome(remote, branch)
		if len(repo) > 0 {
			home = repo
		}
	}
	if len(home) > 0 {
		pathx.RegisterGoctlHome(home)
	}

	if len(protoFile) == 0 {
		return errors.New("missing -o")
	}

	return generator.ProtoTmpl(protoFile)
}
