package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"
    "golang.org/x/crypto/ripemd160"
	"github.com/btcsuite/btcutil/base58"
)

//Lenght of checksum (4 bytes)
const addressChecksumLen = 4

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

type Wallets struct {
	Wallets map[string]*Wallet
}


func NewWallet() *Wallet {
	private, public := newKeyPair()
	w := Wallet{private, public}
	return &w
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}

const version = "version"

func (w *Wallet) GetAddress() []byte {
	// 3 parts of getting an address:
	pubKeyHash := hashPubKey(w.PublicKey)
	version := append([]byte(version), pubKeyHash...)
	checksum := checksum(version)

	payload := append(checksum, version...)
	addr := encode(payload)
	return addr
}

func hashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)

	hasher := ripemd160.New()
	_, err := hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	publicRIPEMD160 := hasher.Sum(nil)
	return publicRIPEMD160
}

//At this point, the payload contains only two parts:
//The version and the hash of the public key.
//The checksum is the first four bytes of the resulted hash.
func checksum(payload []byte) []byte {
	firstStep := sha256.Sum256(payload)
	secondStep := sha256.Sum256(firstStep[:])
	//Slice is exclusive
	return secondStep[:addressChecksumLen + 1]
}

func encode(payload []byte) []byte {
	return []byte(base58.Encode(payload))
}