package miner

import (
	"fmt"
	"runtime"

	"blockchain_go/blockchain"
	"blockchain_go/wallet"

	"github.com/carlescere/scheduler"
)

type Miner struct {
}

func Mine() {
	fmt.Println("Start Miner!")

	wallets, _ := wallet.CreateWallets()
	address := wallets.AddWallet()
	wallets.SaveFiles()

	job := func() {
		fmt.Println("Minning...")
		miner := address
		chain := blockchain.ContinueBlockchain(miner)

		defer chain.Database.Close()

		tx := blockchain.NewTransaction(miner, miner, 0, chain)
		chain.AddBlock([]*blockchain.Transaction{tx})
		fmt.Println("Mine complete: %s", chain.LastHash)
	}

	scheduler.Every(60).Seconds().NotImmediately().Run(job)

	runtime.Goexit()
}
