package miner

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"blockchain_go/structure/blockchain"
	"blockchain_go/structure/validator"

	"github.com/robfig/cron/v3"
)

type Miner struct {
	address string
}

func Mine() {
	// ensure we only start the miner once per process
	mineOnce.Do(func() {
		fmt.Println("Start Miner!")

		c := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))

		c.AddFunc("@every 15s", task)

		go c.Start()

		// keep the main goroutine alive for the miner
		runtime.Goexit()
	})
}

func randomMiner() *Miner {
	validators, _ := validator.CreateValidators()
	addresses := validators.GetAllAdresses()

	n := rand.Int() % len(addresses)

	value := Miner{address: addresses[n]}

	return &value
}

func task() {
	// Miner creates a transaction and enqueues it. When queue reaches 10,
	// publish only the oldest (FIFO) block to the chain.

	fmt.Println("Minning...")
	miner := randomMiner()

	// Need a chain handle to build the transaction (to compute UTXOs).
	chain := blockchain.ContinueBlockchain(miner.address)
	tx := blockchain.NewTransaction(miner.address, miner.address, 0, chain)
	// Close DB now; we'll reopen when publishing.
	chain.Database.Close()

	// enqueue mined transaction into pending buffer
	queueMutex.Lock()
	// set first tx time if this is the first pending tx
	if len(pendingTxs) == 0 {
		firstTxTime = time.Now()
	}
	pendingTxs = append(pendingTxs, tx)
	pendingLen := len(pendingTxs)
	fmt.Printf("Transaction added to pending buffer. Pending size: %d\n", pendingLen)

	// check publish conditions: reached threshold OR timed out since first tx
	var toPublish []*blockchain.Transaction
	now := time.Now()
	if pendingLen >= publishThreshold || (pendingLen > 0 && now.Sub(firstTxTime) >= miningTimeout) {
		if pendingLen >= publishThreshold {
			toPublish = pendingTxs[:publishThreshold]
			if len(pendingTxs) == publishThreshold {
				pendingTxs = nil
				firstTxTime = time.Time{}
			} else {
				pendingTxs = pendingTxs[publishThreshold:]
				// reset firstTxTime to the time of the new first tx
				firstTxTime = now
			}
		} else {
			// timeout case: publish whatever we have
			toPublish = pendingTxs
			pendingTxs = nil
			firstTxTime = time.Time{}
		}
	}
	queueMutex.Unlock()

	if toPublish != nil {
		blockchain.ChainMutex.Lock()
		publishChain := blockchain.ContinueBlockchain(miner.address)
		publishChain.AddBlock(toPublish)
		publishChain.Database.Close()
		fmt.Println("Published pending transactions as a block to chain")
		blockchain.ChainMutex.Unlock()
	}
}

const (
	publishThreshold = 10
	miningTimeout    = 15 * time.Second
)

var (
	pendingTxs  []*blockchain.Transaction
	firstTxTime time.Time
	queueMutex  sync.Mutex
	mineOnce    sync.Once
)
