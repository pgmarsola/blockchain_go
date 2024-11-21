package service

import (
	"blockchain_go/structure/blockchain"
	"blockchain_go/structure/wallet"
	"encoding/json"
	"net/http"
	"os"
)

type Args struct{}

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
