package core

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/dappley/go-dappley/util"
	//"crypto/rand"
	//"crypto/elliptic"
	"crypto/rand"
	"crypto/elliptic"
	"github.com/dappley/go-dappley/crypto/hash"
	"github.com/dappley/go-dappley/crypto/keystore/secp256k1"
)

const version = byte(0x00)
const addressChecksumLen = 4

type KeyPair struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func NewKeyPair() *KeyPair {
	private, public := newKeyPair()
	return &KeyPair{private, public}
}

func (w KeyPair) GenerateAddress() Address {
	pubKeyHash := HashPubKey(w.PublicKey)

	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := util.Base58Encode(fullPayload)

	return NewAddress(fmt.Sprintf("%s", address))
}

func HashPubKey(pubKey []byte) []byte {

	sha := hash.Sha3256(pubKey)
	content := hash.Ripemd160(sha)
	return content
}


func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	private1, err1 := secp256k1.NewECDSAPrivateKey()
	if err1 != nil {
		log.Panic(err1)
	}
	fmt.Printf("\n Generating Private Key is \n")
	fmt.Print(private1)
	fmt.Printf("\n")

	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("\n Private Key is \n")
	fmt.Print(private)
	fmt.Printf("\n")
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}
