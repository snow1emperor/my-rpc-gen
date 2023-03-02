package cli

import (
	"github.com/snow1emperor/my-rpc-gen/rpc/genproxy"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var (
	// VarStringProxyDst ...
	VarStringProxyDst string
	// VarStringProxyTypes ...
	VarStringProxyTypes string
	// VarStringProxyOut ...
	VarStringProxyOut []string
	// VarStringProxyTypeMap ...
	VarStringProxyTypeMap []string
)

func Proxy(_ *cobra.Command, args []string) error {

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	var dst = VarStringProxyDst

	if !filepath.IsAbs(dst) {
		dst = filepath.Join(pwd, dst)
	}
	typeMap := genproxy.ParseMap(VarStringProxyTypeMap)
	if VarStringProxyTypes != "" {
		typeMap["types"] = VarStringProxyTypes
	}
	var ctx = &genproxy.ProxyContext{
		Dst:     dst,
		Out:     VarStringProxyOut,
		TypeMap: typeMap,
	}

	g := genproxy.NewProxyGenerator()
	if err := g.Generate(ctx); err != nil {
		return err
	}
	return nil
}
