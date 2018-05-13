package cnet

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"constants"
)


type input struct {
	PrevTransHash [constants.HASHSIZE]byte //Input transaction that is being spent
	OutIdx        uint32 //Index of the particular transaction
}
type output struct {
	Addr            [constants.ADDRESSSIZE]byte //Address to send the money to
	Amount          uint32 //Amount sending
	Signature       []byte //The hash of the output encrypted with the payers private key
}
type Transaction struct {
	Meta struct {
		TotalTransAmt uint32  //Total amount moving in transaction
		TimePrepared  uint64  //Time of the transaction
		Pubkey     []byte  //Payers public key
		Address    [constants.ADDRESSSIZE]byte  //Payers address
	}
	Inputs [] input `json:"Inputs"`
	Outputs [] output `json:"Outputs"`
	Hash [constants.HASHSIZE]byte //Hash of the whole transaction
}

func checkerror(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
func (tx Transaction) Dump() (size int, ret []byte) {
	tx.Hash = [constants.HASHSIZE]byte{}
	ret, err := json.Marshal(tx)
	checkerror(err)
	tx.Hash = sha256.Sum256(ret)
	ret, err = json.Marshal(tx)
	checkerror(err)
	size = len(ret)
	return
}

func LoadTX(b []byte) (tx Transaction) {
	err := json.Unmarshal(b, tx)
	checkerror(err)
	return
}

func LoadFTX(filename string)(rettx Transaction){
	var pretx Transaction
	raw,err := ioutil.ReadFile(filename)
	checkerror(err)
	fmt.Println(raw)
	err = json.Unmarshal(raw,&pretx)
	checkerror(err)
	fmt.Println(pretx)
	return
}

func Merkleify(txs [] Transaction)([]byte){
	hashlist := make([][]byte,len(txs))
	for i, tx := range txs {
		hashlist[i] = tx.Hash[:]
	}
	for len(hashlist) > 1{
		parseIdx := 0
		hashIdx := 0
		for parseIdx < len(hashlist){
			if parseIdx + 1 != len(hashlist){
				hashlist[hashIdx] = append(hashlist[parseIdx],hashlist[parseIdx+1]...)
				hashIdx += 1
				parseIdx += 2
			}else {
				h := sha256.New()
				copy(hashlist[hashIdx],h.Sum(hashlist[parseIdx])[:])
				hashIdx += 1
				parseIdx += 1
			}
		}
	}
	return hashlist[0]
}
