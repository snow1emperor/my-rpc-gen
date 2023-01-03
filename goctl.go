package main

import (
	"github.com/snow1emperor/my-rpc-gen/cmd"
	"github.com/zeromicro/go-zero/core/load"
	"github.com/zeromicro/go-zero/core/logx"
)

func main() {
	//log.Println(" Start")
	logx.Disable()
	load.Disable()
	cmd.Execute()
	//log.Println(" Exit")
}
