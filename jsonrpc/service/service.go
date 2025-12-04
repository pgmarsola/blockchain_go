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
	TransferTo string `json:"transfer_to"`
	TransferValue int `json:"transfer_value"`
}

type TransferReply struct {
	From string `json:"address"`
	To string `json:"transfer_to"`
	Value int `json:"transfer_value`
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

	tx := blockchain.NewTransaction(args.Address, args.TransferTo, args.TransferValue, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})
	
	*reply = TransferReply{From: args.Address, To: args.TransferTo, Value: args.TransferValue}
	json.NewEncoder(os.Stdout).Encode(*reply)
	return nil
}
