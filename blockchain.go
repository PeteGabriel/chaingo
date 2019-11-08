package main

import (
	"github.com/boltdb/bolt"
	"log"
	"os"
)

const dbFile = "some_file_name"

var blocksBucket string = "blocks"

//Blockchain is an ordered linked-set of blocks
type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte
	_ = bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})

	newBlock := NewBlock(data, lastHash)
	error := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			return err
		}
		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			return err
		}

		bc.tip = newBlock.Hash

		return nil
	})
	if error != nil {
		log.Fatal(error)
	}
}

//Iterator to iterate over blocks in a blockchain
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}
	return bci
}

//NewGenesisBlock creates first block of chain
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

/*
Open a DB file.
Check if there’s a blockchain stored in it.
If there’s a blockchain:
	Create a new Blockchain instance.
	Set the tip of the Blockchain instance to the last block hash stored in the DB.
If there’s no existing blockchain:
	Create the genesis block.
	Store in the DB.
	Save the genesis block’s hash as the last block hash.
	Create a new Blockchain instance with its tip pointing at the genesis block.
*/
func NewBlockchain() *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFile, os.ModeAppend, nil)
	if err != nil {
		log.Fatal(err)
	}
	_ = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Fatal(err)
			}
			_ = b.Put(genesis.Hash, genesis.Serialize())
			_ = b.Put([]byte("l"), genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}
		return nil
	})

	bc := Blockchain{tip, db}

	return &bc
}
