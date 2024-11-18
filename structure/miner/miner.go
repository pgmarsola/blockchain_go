package miner

import (
	"fmt"
	"math/rand"
	"runtime"

	"blockchain_go/structure/blockchain"
	"blockchain_go/structure/validator"

	"github.com/robfig/cron/v3"
)

type Miner struct {
	address string
}

func Mine() {
	fmt.Println("Start Miner!")

	c := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))

	c.AddFunc("@every 30s", task)

	go c.Start()

	runtime.Goexit()
}

func randomMiner() *Miner {
	validators, _ := validator.CreateValidators()
	addresses := validators.GetAllAdresses()

	n := rand.Int() % len(addresses)

	value := Miner{address: addresses[n]}

	return &value
}

func task() {
	fmt.Println("Minning...")
	randomMiner()
	miner := Miner{}.address

	chain := blockchain.ContinueBlockchain(miner)

	defer chain.Database.Close()

	tx := blockchain.NewTransaction(miner, miner, 0, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Println("Mine complete: %x\n", chain.LastHash)
}
