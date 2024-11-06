package validator

import (
	"blockchain_go/structure/utils"
	"encoding/json"
	"log"
	"os"
)

const validatorFile = "./tmp/validators.data"

type Validators struct {
	Validators map[string]*utils.GenericWallet
}

func CreateValidators() (*Validators, error) {
	validators := Validators{}
	validators.Validators = make(map[string]*utils.GenericWallet)

	err := validators.LoadFile()

	return &validators, err
}

func (vs *Validators) AddValidator() string {
	validator := utils.MakeGenericWallet()
	address := string(validator.Address())
	vs.Validators[address] = validator

	return address
}

func (vs *Validators) GetAllAdresses() []string {
	var addresses []string

	for address := range vs.Validators {
		addresses = append(addresses, address)
	}

	return addresses
}

func (vs *Validators) LoadFile() error {
	if _, err := os.Stat(validatorFile); os.IsNotExist(err) {
		return err
	}

	fileContent, err := os.ReadFile(validatorFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(fileContent, vs)

	if err != nil {
		return err
	}

	return nil
}

func (vs *Validators) SaveFiles() {
	jsonData, err := json.Marshal(vs)
	if err != nil {
		log.Panic(err)
	}

	err = os.WriteFile(validatorFile, jsonData, 0666)
	if err != nil {
		log.Panic(err)
	}
}
