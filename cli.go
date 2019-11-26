package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	cmds "./cmd"
)

type CLI struct {
	bc *Blockchain
}


//Run is the entry point of cli.
func (cli *CLI) Run() {
	cli.validateArgs()

	//possible commands
	wallet := cmds.ConfigureCreateWalletCmd()
	create, createData := cmds.ConfigureCreateChainCmd()
	getBalance, getBalanceData := cmds.ConfigureBalanceCmd()
	send, from, to, amount := cmds.ConfigureSendCmd()
	printChain := cmds.ConfigurePrintCmd()

	if len(os.Args) >= 1 {
		switch os.Args[1] {
		case cmds.CreateWalletCmdId:
			_ = wallet.Parse(os.Args[2:])
		case cmds.PrintCmdId:
			_ = printChain.Parse(os.Args[2:])
		case cmds.CreateChainCmdId:
			_ = create.Parse(os.Args[2:])
		case cmds.BalanceCmdId:
			_ = getBalance.Parse(os.Args[2:])
		case cmds.SendCmdId:
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

	if wallet.Parsed() {
		cli.createWallet()
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

func (cli *CLI) createWallet() {

	fmt.Printf("Your new address: %x\n", )
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
	fmt.Printf("\t--%s\t\t\tCreate a new wallet address\n", cmds.CreateWalletCmdId)
	fmt.Printf("\t--%s\t\tPrint the whole chain\n", cmds.PrintCmdId)
	fmt.Printf("\t--%s\tCreate a new blockchain\n", cmds.CreateChainCmdId)
	fmt.Printf("\t--%s\t\tGet balance\n", cmds.BalanceCmdId)
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