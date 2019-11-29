package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
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

//Hash method serializes the transaction and hashes it with the SHA-256 algorithm.
func (t *Transaction) Hash() []byte{
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
  return hash[:]
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
	t.ID = t.Hash()

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
	tx.ID = tx.Hash()

	return &tx
}

/*
Considering that transactions unlock previous outputs,
redistribute their values, and lock new outputs, the following data must be signed:

1 - Public key hashes stored in unlocked outputs. This identifies “sender” of a transaction.
2 - Public key hashes stored in new, locked, outputs. This identifies “recipient” of a transaction.
3 - Values of new outputs.
 */
func (tx *Transaction) Sign(privK ecdsa.PrivateKey, prevTxs map[string]Transaction) error {
	if tx.IsCoinbase() {
		//no real inputs in them.
		return nil
	}

	//include all inputs and outputs & exclude public key and signature.
	txCopy := tx.TrimmedCopy()

	for inID, vin := range txCopy.Vin {
		prevTx := prevTxs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKeyRaw = prevTx.Vout[vin.Vout].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[inID].PubKeyRaw = nil

		r, s, err := ecdsa.Sign(rand.Reader, &privK, txCopy.ID)
		if err != nil {
			return err
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Vin[inID].Signature = signature
	}

	return nil
}

//TrimmedCopy will include all the inputs and outputs, but TXInput.Signature and TXInput.PubKey are set to nil.
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	for _, vin := range tx.Vin {
		inputs = append(inputs, TXInput{vin.Txid, vin.Vout, nil, nil})
	}

	for _, vout := range tx.Vout {
		outputs = append(outputs, TXOutput{vout.Value, vout.PubKeyHash})
	}

	copy := Transaction{tx.ID, inputs, outputs}
	return copy
}

//Verify that a transaction is valid
func (t *Transaction) Verify(prevTxs map[string]Transaction) bool {
	copy := t.TrimmedCopy()
	curve := elliptic.P256()

	for inID, vin := range t.Vin {
		prevTx := prevTxs[hex.EncodeToString(vin.Txid)]
		copy.Vin[inID].Signature = nil
		copy.Vin[inID].PubKeyRaw = prevTx.Vout[vin.Vout].PubKeyHash
		copy.ID = copy.Hash()
		copy.Vin[inID].PubKeyRaw = nil

		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKeyRaw)
		x.SetBytes(vin.PubKeyRaw[:(keyLen / 2)])
		y.SetBytes(vin.PubKeyRaw[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, copy.ID, &r, &s) == false {
			return false
		}
	}
	return true
}