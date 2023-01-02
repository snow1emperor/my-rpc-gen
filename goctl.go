package main

import (
	"github.com/zeromicro/go-zero/core/load"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"my-rpc-gen/cmd"
)

func main() {
	log.Println("Start")
	logx.Disable()
	load.Disable()
	cmd.Execute()
	log.Println("Exit")
}
