package apis

import (
	"blockchain_go/structure/wallet"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const baseUrl = "http://localhost:8080"

func allwallets() {
	response, err := http.Get(baseUrl + "/allwallets")
	if err != nil {
		fmt.Println("Erro ao listar carteiras:", err)
		return
	}
	defer response.Body.Close()

	var wallets []wallet.Wallet
	if err := json.NewDecoder(response.Body).Decode(&wallets); err != nil {
		fmt.Println("Erro ao decodificar resposta:", err)
		return
	}

	for _, wallet := range wallets {
		fmt.Println(wallet)
	}
}

func createwallet() {
	walletData := map[string]string{
		"name": "",
	}

	jsonData, err := json.Marshal(walletData)
	if err != nil {
		fmt.Println("Error encoding json:", err)
		return
	}

	response, err := http.Post(baseUrl+"/createwallet", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating wallet:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusCreated {
		fmt.Println("Wallet created successfully.")
	} else {
		fmt.Println("Failed to create wallet. Status code:", response.StatusCode)
	}
}
