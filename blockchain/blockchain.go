package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"runtime"

	"github.com/dgraph-io/badger"
)

const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"
	genesisData = "First Transaction from Genesis"
)

type Data struct {
	Name string
}

type Blockchain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockchainInterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func DBExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

func ContinueBlockchain(address string) *Blockchain {
	if DBExists() == false {
		fmt.Printf("No existing Blockchain found, please create one!")
		runtime.Goexit()
	}

	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)

	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		err1 := item.Value(func(val []byte) error {
			lastHash = append([]byte{}, val...)
			return nil
		})

		Handle(err1)
		return err
	})

	chain := Blockchain{lastHash, db}
	return &chain
}

func (chain *Blockchain) FindUnspentTransactions(address string) []Transaction {
	var unspentTxs []Transaction
	spentTxOs := make(map[string][]int)

	iter := chain.Interator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTxOs[txID] != nil {
					for _, spentOut := range spentTxOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				if out.CanBeUnlocked(address) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}

			if tx.IsCoinBase() == false {
				for _, in := range tx.Inputs {
					if in.CanUnlock(address) {
						inTxID := hex.EncodeToString(in.ID)
						spentTxOs[inTxID] = append(spentTxOs[inTxID])
					}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return unspentTxs
}

func (chain *Blockchain) FindUTXO(address string) []TxOutput {
	var UTXOs []TxOutput
	unspentTransactions := chain.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

func (chain *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := chain.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Outputs {
			if out.CanBeUnlocked(address) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOuts
}

func InitBlockchain() *Blockchain {
	var lastHash []byte

	if DBExists() == true {
		fmt.Printf("Blockchain already exists")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)

	Handle(err)

	SetName()

	err = db.Update(func(txn *badger.Txn) error {
		cbtx := CoinBaseTx(Data{}.Name, genesisData)
		genesis := Genesis(cbtx)

		fmt.Println("Genesis created")

		err = txn.Set(genesis.Hash, genesis.Serialize())
		Handle(err)

		err = txn.Set([]byte("lh"), genesis.Hash)
		lastHash = genesis.Hash

		return err
	})

	Handle(err)

	blockchain := Blockchain{lastHash, db}

	return &blockchain
}

func (chain *Blockchain) AddBlock(transactions []*Transaction) {
	var lastHash []byte
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		err1 := item.Value(func(val []byte) error {
			lastHash = append([]byte{}, val...)
			return nil
		})
		Handle(err1)
		return err
	})

	Handle(err)

	newBlock := CreateBlock(transactions, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash
		return err
	})

	Handle(err)
}

func (chain *Blockchain) Interator() *BlockchainInterator {
	iter := &BlockchainInterator{chain.LastHash, chain.Database}

	return iter
}

func (iter *BlockchainInterator) Next() *Block {
	var block *Block
	var blockData []byte

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)

		Handle(err)

		err1 := item.Value(func(val []byte) error {
			blockData = append([]byte{}, val...)
			return nil
		})

		Handle(err1)

		block = Deserialize(blockData)

		return err
	})

	Handle(err)

	iter.CurrentHash = block.PrevHash
	return block
}

func SetName() *Data {
	var encoded bytes.Buffer
	var hash [32]byte

	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))

	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(target)
	Handle(err)

	hash = sha256.Sum256(encoded.Bytes())
	value := Data{string(hash[:])}

	return &value
}
