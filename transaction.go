package main

import (
	"encoding/hex"
	"fmt"
	"log"
)


const subsidy = 10

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

type TXOutput struct {
	Value        int
	//For simplicity: user defined wallet address
	ScriptPubKey string
}

type TXInput struct {
	Txid      []byte //the ID of such transaction
	Vout      int   //n index of an output in the transaction
	ScriptSig string
}

func (t *Transaction) setID(){
  t.ID = []byte{1} //TODO change this to something random
}

func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

func (t *Transaction) IsCoinbase() bool {
	return len(t.Vin) == 0
}

//NewCoinbase creates a  special type of transactions,
// which doesn’t require previously existing outputs
func NewCoinbase(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to %s", to)
	}
	in := TXInput{[]byte{}, -1, data}
	out :=  TXOutput{subsidy, to}
	t := Transaction{nil, []TXInput{in}, []TXOutput{out}}
	t.setID()

	return &t
}

func NewTransaction(from, to string, amount int, bc *Blockchain) *Transaction{
	var inputs []TXInput
	var outputs []TXOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)

	//build a list of inputs
	for idx, outputs := range validOutputs {
		idx, err := hex.DecodeString(idx)

		if err != nil {
			log.Fatal(err)
		}

		for _, output := range outputs {
			input := TXInput{idx, output, from}
			inputs = append(inputs, input)
		}
	}

	//build a list of outputs
	outputs = append(outputs, TXOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, from}) // a change
	}

	tx := Transaction{nil, inputs, outputs}
	tx.setID()

	return &tx
}