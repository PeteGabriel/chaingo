package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"strconv"

)

const (
	targetBits = 24
	maxNonce = math.MaxInt64
)

type ProofOfWork struct {
	block  *Block
	target *big.Int //because of the way we use to compare hashes
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	return &ProofOfWork{b, target}
}

/*
	Prepare data.
	Hash it with SHA-256.
	Convert the hash to a big integer.
	Compare the integer with the target.
*/
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	nonce := 0
	fmt.Printf("Mining a new block")
	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash := sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}

		return nonce, hash[:]
	}
	fmt.Print("\n\n")
	return nonce, []byte{}
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	return hashInt.Cmp(pow.target) == -1
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.HashTransactions(),
			intToHex(pow.block.Timestamp),
			intToHex(int64(targetBits)),
			intToHex(int64(nonce)),
		},
		[]byte{},
	)
	return data
}

func intToHex(toConvert int64) []byte {
	return []byte(strconv.FormatInt(toConvert, 16))
}
