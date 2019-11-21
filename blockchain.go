package main

import (
	"encoding/hex"
	"github.com/boltdb/bolt"
	"fmt"
	"log"
	"os"
)

const (
	dbFile = "some_file_name"
	blocksBucket = "blocks"
	genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
	bucketName = "l"
)
//Blockchain is an ordered linked-set of blocks
type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

func (bc *Blockchain) MineBlock(transactions  []*Transaction) {
	var lastHash []byte
	_ = bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte(bucketName))
		return nil
	})

	newBlock := NewBlock(transactions, lastHash)
	err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			return err
		}
		err = b.Put([]byte(bucketName), newBlock.Hash)
		if err != nil {
			return err
		}

		bc.tip = newBlock.Hash

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

//Iterator to iterate over blocks in a blockchain
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}
	return bci
}

//NewGenesisBlock creates first block of chain
func NewGenesisBlock(cb *Transaction) *Block {
	return NewBlock([]*Transaction{cb}, []byte{})
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
func NewBlockchain(addr string) *Blockchain {

	if dbExists() == false {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte(bucketName))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc
}


func CreateBlockchain(addr string) *Blockchain {
	if dbExists() {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		cbtx := NewCoinbase(addr, genesisCoinbaseData)
		genesis := NewGenesisBlock(cbtx)

		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}

		err = b.Put(genesis.Hash, genesis.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte(bucketName), genesis.Hash)
		if err != nil {
			log.Panic(err)
		}

		tip = genesis.Hash

		return nil
	})
}


func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}


func (bc *Blockchain) findUnspentTransactions(addr string) []Transaction {
	var unspent []Transaction
	spent := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()
		for _, t := range block.Transactions{
			id := hex.EncodeToString(t.ID)
		Outputs:
			for outIdx, out := range t.Vout  {
				//check if the output was already spent
				if spent[id] != nil {
					for _, spentOut := range spent[id] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.CanBeUnlockedWith(addr) {
					unspent = append(unspent, *t)
				}
			}

			if t.IsCoinbase() == false {
				for _, in := range t.Vin {
					if in.CanUnlockOutputWith(addr) {
						inTxID := hex.EncodeToString(in.Txid)
						spent[inTxID] = append(spent[inTxID], in.Vout)
					}
				}
			}

			if len(block.PrevBlockHash) == 0 {
				break
			}
		}
		return unspent
	}
}

func (bc *Blockchain) FindUnspentOutputs(addr string) []TXOutput {
	var unspentOutputs []TXOutput
	unspentTransactions := bc.findUnspentTransactions(addr)
	for _, t := range unspentTransactions {
		for _, output := range t.Vout {
			if output.CanBeUnlockedWith(addr) {
				unspentOutputs = append(unspentOutputs, output)
			}
		}
	}
	return unspentOutputs
}

//FindSpendableOutputs find all unspent outputs. It groups by transaction IDs.
func (bc *Blockchain) FindSpendableOutputs(from string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTransactions := bc.findUnspentTransactions(from)
	acc := 0

	Work:
		for _, t := range unspentTransactions {
			for idx,output := range t.Vout {
				if output.CanBeUnlockedWith(from) && acc < amount {
					acc += output.Value
					txID := hex.EncodeToString(t.ID)

					unspentOutputs[txID] = append(unspentOutputs[txID], idx)

					if acc > amount {
						break Work
					}
				}
			}
		}

	return acc, unspentOutputs
}