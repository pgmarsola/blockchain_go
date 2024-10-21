package conection

import (
	"blockchain_go/structure/wallet"
	"encoding/json"
	"fmt"
	"net/http"
)

const baseUrl = "http://localhost:8080"

func wallets() {
	response, err := http.Get(baseUrl + "/wallet")
	if err != nil {
		fmt.Println("Erro ao listar tarefas:", err)
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
