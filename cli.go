package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type CLI struct {
	bc *Blockchain
}

//Run is the entrypoint of cli.
func (cli *CLI) Run() {
	cli.validateArgs()

	//possible commands
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	createBlockchainData := createBlockchainCmd.String("address", "",  "Wallet address")

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	addBlockData := addBlockCmd.String("data", "", "Block data")

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	getBalanceData := addBlockCmd.String("address", "", "Wallet address")

	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	switch os.Args[1] {
		case "addblock":
			_ = addBlockCmd.Parse(os.Args[2:])
		case "printchain":
			_ = printChainCmd.Parse(os.Args[2:])
	    case "createblockchain":
	    	_ = createBlockchainCmd.Parse(os.Args[2:])
	    case "getbalance":
		    _ = getBalanceCmd.Parse(os.Args[2:])
	    default:
			cli.printUsage()
			os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage() //print instructions
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if createBlockchainCmd.Parsed() {
		cli.createBlockchain(*createBlockchainData)
	}

	if getBalanceCmd.Parsed() {
		cli.getBalance(*getBalanceData)
	}
}

func (cli *CLI) addBlock(data string) {
	//cli.bc.AddBlock(data)
	fmt.Println("Success!")
}

func (cli *CLI) createBlockchain(addr string) {
	cli.bc = NewBlockchain(addr)
	fmt.Println("Success!")
}

func (cli *CLI) printChain() {
	bci := cli.bc.Iterator()
	for {
		block := bci.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Transactions: %s\n", block.Transactions)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("Proof of Work: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

func (cli *CLI) printUsage() {

}

func (cli *CLI) validateArgs() {

}


func (cli *CLI) getBalance(addr string) {
	bc := NewBlockchain(addr)
	defer bc.db.Close()

	balance := 0
	unspentOutputs := bc.FindUnspentOutputs(addr)

	for _, output := range unspentOutputs {
		balance += output.Value
	}
	fmt.Printf("Balance of '%s': %d\n", addr, balance)
}
