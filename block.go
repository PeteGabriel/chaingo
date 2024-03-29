package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

//Block represents the stored valuable info
type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{
		time.Now().Unix(),
		transactions,
		prevBlockHash,
		[]byte{},
		0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Fatal(err)
	}
	return result.Bytes()
}

// HashTransactions hashes of each transaction,
// concatenate them,
// and get a hash of the concatenated combination.
func (b *Block) HashTransactions() []byte {
	var hashes []byte
	for _, t := range b.Transactions  {
		for h := range t.ID {
			hashes = append(hashes, byte(h))
		}
	}
	hash := sha256.Sum256(hashes)
	return hash[:]
}


func DeserializeBlock(d []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Fatal(err)
	}
	return &block
}
