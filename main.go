package main

import (
	"blockchain_go/jsonrpc/server"
	"blockchain_go/structure/blockchain"
	"blockchain_go/structure/miner"
)

func main() {
	go server.Run()

	chain := blockchain.InitBlockchain()
	defer chain.Database.Close()

	miner.Mine()
}
