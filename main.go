package main

import (
	"blockchain_go/cli"
	"blockchain_go/jsonrpc/server"
)

func main() {
	go server.Run()

	cmd := cli.CommandLine{}
	cmd.Run()
}
