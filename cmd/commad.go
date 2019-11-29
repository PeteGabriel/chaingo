package cmd

import (
	"flag"
)

const (
	CreateWalletCmdId = "createwallet"
	CreateChainCmdId  = "createblockchain"
	BalanceCmdId      = "getbalance"
	SendCmdId         = "send"
	PrintCmdId        = "printchain"
)


func ConfigureCreateWalletCmd() *flag.FlagSet {
	fs := flag.NewFlagSet(CreateWalletCmdId, flag.ExitOnError)
	return fs
}

func ConfigureCreateChainCmd() (*flag.FlagSet, *string) {
	fs := flag.NewFlagSet(CreateChainCmdId, flag.ExitOnError)
	data := fs.String("address", "", "Wallet address")
	return fs, data
}

func ConfigureBalanceCmd() (*flag.FlagSet, *string) {
	fs := flag.NewFlagSet(BalanceCmdId, flag.ExitOnError)
	data := fs.String("address", "", "Wallet address")
	return fs, data
}

func ConfigureSendCmd() (*flag.FlagSet, *string, *string, *int) {
	fs := flag.NewFlagSet(SendCmdId, flag.ExitOnError)
	from := fs.String("from", "", "Sender address")
	to := fs.String("to", "", "Receiver address")
	amount := fs.Int("amount", 0, "Amount to fs")
	return fs, from, to, amount
}

func ConfigurePrintCmd() *flag.FlagSet {
	return flag.NewFlagSet(PrintCmdId, flag.ExitOnError)
}