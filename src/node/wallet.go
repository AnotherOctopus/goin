package node

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"cnet"
)

type Wallet struct {
	Keys       []*rsa.PrivateKey
	Address    [][]byte
	ClaimedTxs []cnet.Transaction
}
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
func NewWallet(numkeys int) (* Wallet) {
	w := new(Wallet)
	w.Keys = make([]*rsa.PrivateKey,numkeys)
	w.Address = make([][]byte,numkeys)
	for i,_ := range w.Keys{
		skey, err := rsa.GenerateKey(rand.Reader, 1000)
		CheckError(err)
		w.Keys[i]	= skey
		w.Address[i] = pkeytoaddress(skey.PublicKey)
	}
	return w
}
func pkeytoaddress( pkey rsa.PublicKey)([]byte){

	return []byte{0x00}
}