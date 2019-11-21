package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type CLI struct {
	bc *Blockchain
}

const (
	createCmd = "createblockchain"
	addBlockCmd = "addblock"
	getBalanceCmd = "getbalance"
	sendCmd = "send"
	printCmd = "printchain"
)

//Run is the entry point of cli.
func (cli *CLI) Run() {
	cli.validateArgs()

	//possible commands
	create := flag.NewFlagSet(createCmd, flag.ExitOnError)
	createData := create.String("address", "", "Wallet address")

	addBlock := flag.NewFlagSet(addBlockCmd, flag.ExitOnError)
	addBlockData := addBlock.String("data", "", "Block data")

	getBalance := flag.NewFlagSet(getBalanceCmd, flag.ExitOnError)
	getBalanceData := getBalance.String("address", "", "Wallet address")

	send := flag.NewFlagSet(sendCmd, flag.ExitOnError)
	from := send.String("from", "", "Sender address")
	to := send.String("to", "", "Receiver address")
	amount := send.Int("amount", 0, "Amount to send")

	printChain := flag.NewFlagSet(printCmd, flag.ExitOnError)


	if len(os.Args) >= 1 {
		switch os.Args[1] {
		case addBlockCmd:
			_ = addBlock.Parse(os.Args[2:])
		case printCmd:
			_ = printChain.Parse(os.Args[2:])
		case createCmd:
			_ = create.Parse(os.Args[2:])
		case getBalanceCmd:
			_ = getBalance.Parse(os.Args[2:])
		case sendCmd:
			_ = send.Parse(os.Args[2:])
		default:
			cli.printUsage()
			os.Exit(1)
		}
	} else {
		cli.printUsage()
		os.Exit(1)
		return
	}

	if addBlock.Parsed() {
		if *addBlockData == "" {
			addBlock.Usage() //print instructions
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}

	if printChain.Parsed() {
		cli.printChain()
	}

	if create.Parsed() {
		err := cli.createBlockchain(*createData)
		if err != nil {
			log.Fatal(err)
		}
	}

	if getBalance.Parsed() {
		if *getBalanceData == "" {
			getBalance.Usage()
			os.Exit(1)
		}
		err := cli.getBalance(*getBalanceData)
		if err != nil {
			log.Fatal(err)
		}
	}

	if send.Parsed() {
		if *from == "" || *to == "" || *amount <= 0 {
			send.Usage()
			os.Exit(1)
		}
		err := cli.send(*from, *to, *amount)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (cli *CLI) addBlock(data string) {
	//cli.bc.AddBlock(data)
	fmt.Println("Success!")
}

func (cli *CLI) createBlockchain(addr string)  error{
	cli.bc = CreateBlockchain(addr)
	fmt.Println("Done!")
	return cli.bc.db.Close()
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
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
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
	fmt.Printf("\t--%s\t\t\tAdd a new block to the chain\n", addBlockCmd)
	fmt.Printf("\t--%s\t\tPrint the whole chain\n", printCmd)
	fmt.Printf("\t--%s\tCreate a new blockchain\n", createCmd)
	fmt.Printf("\t--%s\t\tGet balance\n", getBalanceCmd)
}

func (cli *CLI) validateArgs() {

}

func (cli *CLI) getBalance(addr string) error {
	bc := NewBlockchain(addr)

	balance := 0
	unspentOutputs := bc.FindUnspentOutputs(addr)

	for _, output := range unspentOutputs {
		balance += output.Value
	}
	fmt.Printf("Balance of '%s': %d\n", addr, balance)

	return bc.db.Close()
}


func (cli *CLI) send(from, to string, amount int) error {
	bc := NewBlockchain(from)
	t := NewTransaction(from, to, amount, bc)
	bc.MineBlock([]*Transaction{t})
	fmt.Println("Success!")
	return bc.db.Close()
}