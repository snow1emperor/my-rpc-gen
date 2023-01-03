package main

import (
	"github.com/teamgram/marmota/pkg/commands"

	{{.imports}}
)

func main() {
	commands.Run(server.New())
}
