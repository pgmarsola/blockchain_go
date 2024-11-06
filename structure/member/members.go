package member

import (
	"blockchain_go/structure/utils"
	"encoding/json"
	"log"
	"os"
)

const memberFile = "./tmp/members.data"

type Members struct {
	Members map[string]*utils.GenericWallet
}

func CreateMembers() (*Members, error) {
	members := Members{}
	members.Members = make(map[string]*utils.GenericWallet)

	err := members.LoadFile()

	return &members, err
}

func (ms *Members) AddMember() string {
	member := utils.MakeGenericWallet()
	address := string(member.Address())
	ms.Members[address] = member

	return address
}

func (ms *Members) GetAllAdresses() []string {
	var addresses []string

	for address := range ms.Members {
		addresses = append(addresses, address)
	}

	return addresses
}

func (ms *Members) LoadFile() error {
	if _, err := os.Stat(memberFile); os.IsNotExist(err) {
		return err
	}

	fileContent, err := os.ReadFile(memberFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(fileContent, ms)

	if err != nil {
		return err
	}

	return nil
}

func (ms *Members) SaveFiles() {
	jsonData, err := json.Marshal(ms)
	if err != nil {
		log.Panic(err)
	}

	err = os.WriteFile(memberFile, jsonData, 0666)
	if err != nil {
		log.Panic(err)
	}
}
