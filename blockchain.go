package main

import "github.com/boltdb/bolt"

const dbFile = "some_file_name"
var BlocksBucket string = "blocks"

//Blockchain is an ordered linked-set of blocks
type Blockchain struct {
	tips []byte
	db   *bolt.DB
}

func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte
	bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))
		lastHash = b.Get([]byte("1"))
		return nil
	})

	newBlock := NewBlock(data, lastHash)
	bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))
		b.Put(newBlock.Hash, newBlock.Serialize())
		b.Put([]byte("l"), newBlock.Hash)
		bc.tips = newBlock.Hash

		return nil
	})
}

//Iterator to iterate over blocks in a blockchain
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tips, bc.db}
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
	db, _ := bolt.Open(dbFile, 0600, nil)
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))
		if b == nil {
			genesis := NewGenesisBlock()
			b, _ := tx.CreateBucket([]byte(BlocksBucket))
			b.Put(genesis.Hash, genesis.Serialize())
			b.Put([]byte("l"), genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}
		return nil
	})

	bc := Blockchain{tip, db}

	return &bc
}
