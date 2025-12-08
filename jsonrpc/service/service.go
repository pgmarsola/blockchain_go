package service

import (
	"blockchain_go/structure/blockchain"
	"blockchain_go/structure/wallet"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Args struct {
    Address string `json:"address"`
	To string `json:"to"`
	Amount int `json:"amount"`
}

type TransferReply struct {
	From string `json:"address"`
	To string `json:"to"`
	Value int `json:"amount"`
}

type MintReply struct {
	To string `json:"to"`
	Amount int `json:"amount"`
	Message string `json:"message"`
}

type Server string

func (t *Server) BlockNumber(r *http.Request, args *Args, reply *[]byte) (err error) {
	chain := blockchain.ContinueBlockchain("")
	defer chain.Database.Close()

	*reply = chain.LastHash
	json.NewEncoder(os.Stdout).Encode(*reply)
	return nil
}

func (t *Server) CreateWallet(r *http.Request, args *Args, reply *string) (err error) {
	wallets, _ := wallet.CreateWallets()
	address := wallets.AddWallet()
	wallets.SaveFiles()

	*reply = address
	json.NewEncoder(os.Stdout).Encode(*reply)
	return nil
}

func (t *Server) AllWallets(r *http.Request, args *Args, reply *[]string) (err error) {
	wallets, _ := wallet.CreateWallets()
	addresses := wallets.GetAllAdresses()

	*reply = addresses
	json.NewEncoder(os.Stdout).Encode(*reply)
	return nil
}

func (t *Server) GetBalance(r *http.Request, args *Args, reply *int) (err error) {
    chain := blockchain.ContinueBlockchain(args.Address)
    defer chain.Database.Close()

    fmt.Println(args.Address)

    balance := 0
    UTXOs := chain.FindUTXO(args.Address)
    for _, out := range UTXOs {
        balance += out.Value
    }

    *reply = balance
    json.NewEncoder(os.Stdout).Encode(*reply)
    return nil
}

func (t *Server) Transfer(r *http.Request, args *Args, reply *TransferReply) (err error) {
	blockchain.ChainMutex.Lock()
    defer blockchain.ChainMutex.Unlock()
	
	chain := blockchain.ContinueBlockchain(args.Address)
	defer chain.Database.Close()

	tx := blockchain.NewTransaction(args.Address, args.To, args.Amount, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})
	
	*reply = TransferReply{From: args.Address, To: args.To, Value: args.Amount}
	json.NewEncoder(os.Stdout).Encode(*reply)
	return nil
}

func (t *Server) Mint(r *http.Request, args *Args, reply *MintReply) (err error) {
	blockchain.ChainMutex.Lock()
	defer blockchain.ChainMutex.Unlock()

	chain := blockchain.ContinueBlockchain(args.To)
	defer chain.Database.Close()

	mintCoinTx := blockchain.MintCoinTx(args.To, args.Amount)
	chain.AddBlock([]*blockchain.Transaction{mintCoinTx})

	*reply = MintReply{To: args.To, Amount: args.Amount, Message: "Minting Successful"}
	json.NewEncoder(os.Stdout).Encode(*reply)
	return nil
}
