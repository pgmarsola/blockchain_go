package miner

import (
	"fmt"
	"runtime"

	"blockchain_go/structure/blockchain"
	"blockchain_go/structure/wallet"

	"github.com/robfig/cron/v3"
)

type Miner struct {
	address string
}

func Mine() {
	fmt.Println("Start Miner!")

	createMiner()

	c := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))

	c.AddFunc("@every 1m", task)

	go c.Start()

	runtime.Goexit()
}

func createMiner() *Miner {
	wallets, _ := wallet.CreateWallets()
	address := wallets.AddWallet()
	wallets.SaveFiles()

	value := Miner{address: address}

	return &value
}

func task() {
	fmt.Println("Minning...")
	miner := Miner{}.address
	chain := blockchain.ContinueBlockchain(miner)

	defer chain.Database.Close()

	tx := blockchain.NewTransaction(miner, miner, 0, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Println("Mine complete: %x\n", chain.LastHash)
}
