package main

import (
	"github.com/boltdb/bolt"
	"log"
)

const blockBucket = "blocks"
/*
BlockchainIterator exists because we don’t want to load all the blocks into memory
so we’ll read them one by one.
For this purpose, we’ll need a blockchain iterator:
*/
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (bci *BlockchainIterator) Next() *Block {
	var block *Block
	err := bci.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockBucket))
		encodedBlock := b.Get(bci.currentHash)
		block = DeserializeBlock(encodedBlock)
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	bci.currentHash = block.PrevBlockHash
	return block
}
