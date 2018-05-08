package node

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"cnet"
)

type Wallet struct {
	Keys       *rsa.PrivateKey
	Address    []byte
	ClaimedTxs []cnet.Transaction
}
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
func NewWallet() (* Wallet) {
	w := new(Wallet)
	skey, err := rsa.GenerateKey(rand.Reader, 1000)
	CheckError(err)
	w.Keys = new(rsa.PrivateKey)
	w.Keys = skey
	return w
}
