package main

import "github.com/boltdb/bolt"

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
	bci.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte{byte(BlocksBucket)})
		encodedBlock := b.Get(bci.currentHash)
		block = DeserializeBlock(encodedBlock)
		return nil
	})

	bci.currentHash = block.PrevBlockHash
	return block
}
