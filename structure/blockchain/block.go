package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Hash       []byte
	MerkleRoot []byte
	PrevHash   []byte
	Nonce      int
	Timestamp  int64
}

func (b *Block) HashTransactions() []byte {
	// Merkle root is computed at block creation time and stored in the block.
	// If available, return it. Otherwise return empty slice.
	if len(b.MerkleRoot) > 0 {
		return b.MerkleRoot
	}

	return []byte{}
}

func GetUnixTimestamp() int64 {
	return time.Now().Unix()
}

func CreateBlock(txs []*Transaction, prevHash []byte) *Block {
	// create block metadata (do not store transactions inside the block)
	block := &Block{Hash: []byte{}, MerkleRoot: nil, PrevHash: prevHash, Nonce: 0, Timestamp: 0}

	// compute Merkle root for the provided transactions and store only the root in the block
	var txHashes [][]byte
	for _, tx := range txs {
		txHashes = append(txHashes, tx.ID)
	}
	block.MerkleRoot = BuildMerkleRoot(txHashes)
	pow := NewProof(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	block.Timestamp = GetUnixTimestamp()

	return block
}

func Genesis(coinBase *Transaction) *Block {
	return CreateBlock([]*Transaction{coinBase}, []byte{})
}

func (b *Block) Serialize() []byte {
	// Serialize only block metadata (without transactions)
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	type blockData struct {
		Hash       []byte
		MerkleRoot []byte
		PrevHash   []byte
		Nonce      int
		Timestamp  int64
	}

	data := blockData{
		Hash:       b.Hash,
		MerkleRoot: b.MerkleRoot,
		PrevHash:   b.PrevHash,
		Nonce:      b.Nonce,
		Timestamp:  b.Timestamp,
	}

	err := encoder.Encode(&data)
	Handle(err)

	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	type blockData struct {
		Hash       []byte
		MerkleRoot []byte
		PrevHash   []byte
		Nonce      int
		Timestamp  int64
	}

	var bd blockData
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&bd)
	Handle(err)

	block := Block{
		Hash:       bd.Hash,
		MerkleRoot: bd.MerkleRoot,
		PrevHash:   bd.PrevHash,
		Nonce:      bd.Nonce,
		Timestamp:  bd.Timestamp,
	}

	return &block
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// SerializeTransactions serializes a slice of transactions for storing
func SerializeTransactions(txs []*Transaction) []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(txs)
	Handle(err)

	return res.Bytes()
}

// DeserializeTransactions decodes a slice of transactions
func DeserializeTransactions(data []byte) []*Transaction {
	var txs []*Transaction
	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&txs)
	Handle(err)

	return txs
}
