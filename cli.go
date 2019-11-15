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
	createBlockchainData := createBlockchainCmd.String("address", "", "Wallet address")

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	addBlockData := addBlockCmd.String("data", "", "Block data")

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	getBalanceData := getBalanceCmd.String("address", "", "Wallet address")

	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	from := sendCmd.String("from", "", "Sender address")
	to := sendCmd.String("to", "", "Receiver address")
	amount := sendCmd.Int("amount", 0, "Amount to send")

	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)


	if len(os.Args) >= 1 {
		switch os.Args[1] {
		case "addblock":
			_ = addBlockCmd.Parse(os.Args[2:])
		case "printchain":
			_ = printChainCmd.Parse(os.Args[2:])
		case "createblockchain":
			_ = createBlockchainCmd.Parse(os.Args[2:])
		case "getbalance":
			_ = getBalanceCmd.Parse(os.Args[2:])
		case "send":
			_ = sendCmd.Parse(os.Args[2:])
		default:
			cli.printUsage()
			os.Exit(1)
		}
	} else {
		cli.printUsage()
		os.Exit(1)
		return
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

	if sendCmd.Parsed() {
		cli.send(*from, *to, *amount)
	}
}

func (cli *CLI) addBlock(data string) {
	//cli.bc.AddBlock(data)
	fmt.Println("Success!")
}

func (cli *CLI) createBlockchain(addr string) {
	cli.bc = NewBlockchain(addr)
	fmt.Println("Done!")
}

func (cli *CLI) printChain() {
	if cli.bc == nil {
		fmt.Println("No blockchain has been created yet.")
		return
	}
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
	fmt.Println("Usage:")
	fmt.Println("\tmain.exe [--option] [<arguments>]")
	fmt.Println("Options:")
	fmt.Println("\t--addblock\t\t\tAdd a new block to the chain")
	fmt.Println("\t--printchain\t\tPrint the whole chain")
	fmt.Println("\t--createblockchain\tCreate a new blockchain")
	fmt.Println("\t--getbalance\t\tGet balance")
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


func (cli *CLI) send(from, to string, amount int) {
	bc := NewBlockchain(from)
	defer bc.db.Close()
	t := NewTransaction(from, to, amount, bc)
	bc.AddBlock([]*Transaction{t})
	fmt.Println("Success!")
}