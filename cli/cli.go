package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"blockchain_go/structure/blockchain"
	"blockchain_go/structure/miner"
	"blockchain_go/structure/wallet"
)

type CommandLine struct{}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage: ")
	fmt.Println("getbalance -address ADDRESS - get balance for address")
	fmt.Println("createblockchain - creates a blockchain")
	fmt.Println("printchain - Prints the blocks in the chain")
	fmt.Println("send -from FROM -to TO -amount AMOUNT - send amount from to")
	fmt.Println("createwallet - create a new Wallet")
	fmt.Println("listaddresses - list all adresses in our wallet file")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		{
			cli.printUsage()
			runtime.Goexit()
		}
	}
}

func (cli *CommandLine) listAddresses() {
	wallets, _ := wallet.CreateWallets()
	addresses := wallets.GetAllAdresses()

	for _, address := range addresses {
		fmt.Printf(address)
	}
}

func (cli *CommandLine) createwallet() {
	wallets, _ := wallet.CreateWallets()
	address := wallets.AddWallet()
	wallets.SaveFiles()

	fmt.Printf("New address is: %s\n", address)
}

func (cli *CommandLine) printChain() {
	chain := blockchain.ContinueBlockchain("")
	defer chain.Database.Close()

	iter := chain.Interator()

	for {
		block := iter.Next()
		fmt.Printf("Prev Hash: %x\n", block.PrevHash)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("PoW %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) createBlockchain() {
	chain := blockchain.InitBlockchain()
	defer chain.Database.Close()

	fmt.Println("Blockchain created successfully!")

	miner.Mine()
}

func (cli *CommandLine) getBalance(address string) {
	chain := blockchain.ContinueBlockchain(address)
	defer chain.Database.Close()

	balance := 0
	UTXOs := chain.FindUTXO(address)
	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s:  %d\n", address, balance)
}

func (cli *CommandLine) send(from string, to string, amount int) {
	chain := blockchain.ContinueBlockchain(from)
	defer chain.Database.Close()

	tx := blockchain.NewTransaction(from, to, amount, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Println("Success!")
}

func (cli *CommandLine) Run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "address of wallet")
	sendFrom := sendCmd.String("from", "", "address of wallet sender")
	sendTo := sendCmd.String("to", "", "address of wallet receiver")
	sendAmount := sendCmd.Int("amount", 0, "amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listaddresses":
		err := listAddressesCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getBalance(*getBalanceAddress)
	}
	if createBlockchainCmd.Parsed() {
		cli.createBlockchain()
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if createWalletCmd.Parsed() {
		cli.createwallet()
	}

	if listAddressesCmd.Parsed() {
		cli.listAddresses()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}
