package apis

import (
	"blockchain_go/structure/wallet"
	"encoding/json"
	"net/http"
)

func handleCreateWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		createWallet(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
	}
}
func handleAllWallets(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		allWallets(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
	}
}

func allWallets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	wallets, _ := wallet.CreateWallets()
	addresses := wallets.GetAllAdresses()

	for _, address := range addresses {
		json.NewEncoder(w).Encode(address)
	}
}

func createWallet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	wallets, _ := wallet.CreateWallets()
	address := wallets.AddWallet()
	wallets.SaveFiles()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(address)
}
