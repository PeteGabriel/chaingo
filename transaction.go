package main

import "fmt"


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