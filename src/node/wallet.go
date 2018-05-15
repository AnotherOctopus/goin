package node

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"constants"
	"crypto/x509"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"strconv"
	"encoding/hex"
)

type Wallet struct {
	Keys       []*rsa.PrivateKey
	Address    [][constants.ADDRESSSIZE]byte
	ClaimedTxs [][constants.HASHSIZE]byte
}
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
func NewWallet(numkeys int) (* Wallet) {
	w := new(Wallet)
	w.Keys = make([]*rsa.PrivateKey,numkeys)
	w.Address = make([][constants.ADDRESSSIZE]byte,numkeys)
	for i,_ := range w.Keys{
		skey, err := rsa.GenerateKey(rand.Reader, 2046)
		CheckError(err)
		w.Keys[i]	= skey
		w.Address[i] = pkeytoaddress(skey.PublicKey)
	}
	return w
}
func (w Wallet) Dump()([]byte){
	ret := make([]byte,constants.PRIVKEYSIZE*len(w.Keys)+
		               constants.ADDRESSSIZE*len(w.Address)+
		               constants.HASHSIZE*len(w.ClaimedTxs))
	keyidx := 0
	for _,key := range w.Keys{
		copy(ret[keyidx:keyidx+constants.PRIVKEYSIZE],x509.MarshalPKCS1PrivateKey(key))
		keyidx += constants.PRIVKEYSIZE
	}

	for _, addr := range w.Address{
		copy(ret[keyidx:keyidx+constants.ADDRESSSIZE],addr[:])
		keyidx += constants.ADDRESSSIZE
	}

	for _, ctx := range w.ClaimedTxs{
		copy(ret[keyidx:keyidx+constants.HASHSIZE],ctx[:])
		keyidx += constants.HASHSIZE
	}

	return ret
}

func LoadWallet(b []byte)(*Wallet){
	w := NewWallet(0)
	idx := 0
	key,err := x509.ParsePKCS1PrivateKey(bytes.Trim(b[idx:idx+constants.PRIVKEYSIZE],string(0x00)))
	CheckError(err)
	for err == nil {
		w.Keys = append(w.Keys, key)
		idx += constants.PRIVKEYSIZE
		if idx+constants.PRIVKEYSIZE >= len(b){
			break
		}
		key,err = x509.ParsePKCS1PrivateKey(bytes.Trim(b[idx:idx+constants.PRIVKEYSIZE],string(0x00)))
	}
	var addressbuffer [constants.ADDRESSSIZE]byte
	for addrIdx := 0 ; addrIdx < len(w.Keys) ; addrIdx += 1{
		copy(addressbuffer[:], b[idx:idx+constants.ADDRESSSIZE])
		w.Address = append(w.Address,addressbuffer)
		idx += constants.ADDRESSSIZE
	}
	var txbuffer [constants.HASHSIZE]byte
	for idx < len(b){
		copy(txbuffer[:], b[idx:idx+constants.HASHSIZE])
		w.ClaimedTxs = append(w.ClaimedTxs,txbuffer)
		idx += constants.HASHSIZE
	}
	return w
}
func pkeytoaddress( pkey rsa.PublicKey)([constants.ADDRESSSIZE]byte){
	EBytes := make([]byte, 8)
	NBytes := pkey.N.Bytes()
	var ret [constants.ADDRESSSIZE]byte
	h := sha256.New()
	binary.LittleEndian.PutUint64(EBytes,uint64(pkey.E))
	subhash := h.Sum(NBytes)
	for i := 0 ; i < constants.ADDRESSSIZE; i += 1 {
		subhash = append(subhash, EBytes[i])
	}
	copy(ret[:],h.Sum(subhash))

	return ret
}

func (w Wallet) String()(string){
	retstring := ""
	retstring += "Number of Addresses: "
	retstring += strconv.Itoa(len(w.Keys)) + "\n"
	for _,addr := range w.Address {
		retstring += hex.EncodeToString(addr[:]) + "\n"
	}
	return retstring

}