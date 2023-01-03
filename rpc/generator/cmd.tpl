package main

import (

	{{.imports}}
)

func main() {
	commands.Run(server.New())
}
