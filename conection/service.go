package conection

import (
	"blockchain_go/structure/wallet"
	"encoding/json"
	"log"
	"net/http"
)

func Client() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/wallet", handleWallet)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Blockchain Local Go"))
}

func handleWallet(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listWallets(w, r)
	case http.MethodPost:
		createWallet(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
	}
}

func listWallets(w http.ResponseWriter, r *http.Request) {
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
