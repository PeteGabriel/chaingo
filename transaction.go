package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)


const subsidy = 10

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

type TXOutput struct {
	Value        int
	PubKeyHash []byte
}

type TXInput struct {
	Txid      []byte //the ID of such transaction
	Vout      int   //n index of an output in the transaction
	Signature []byte
	PubKeyRaw []byte
}

func (t *Transaction) setID(){
  var encoded bytes.Buffer
  var hash [32]byte

  //Encode a transaction to calculate the hash later
  enc := gob.NewEncoder(&encoded)
  err := enc.Encode(t)
  if err != nil {
  	log.Panic(err)
  }

  hash = sha256.Sum256(encoded.Bytes())
  //Since each hash is unique, it helps applying it to the ID field.
  t.ID = hash[:]
}

func (in *TXInput) UsesKey(publicKeyHash []byte) bool {
	lockingKey := HashPubKey(in.PubKeyRaw)
	return bytes.Compare(lockingKey, publicKeyHash) == 0
}

func (out *TXOutput) Lock(addr []byte) {
	pubKeyHash := PublicKeyHash(string(addr))
	out.PubKeyHash = pubKeyHash
}

func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

func (t *Transaction) IsCoinbase() bool {
	return len(t.Vin) == 1 && len(t.Vin[0].Txid) == 0 && t.Vin[0].Vout == -1
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

	if acc < amount {
		fmt.Println("ERROR: Not enough funds")
		os.Exit(1)
	}

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