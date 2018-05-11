package cnet

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	MAXTRANSNETSIZE = 10000
)
type input struct {
	PrevTransHash [32]byte //Input transaction that is being spent
	OutIdx        uint32 //Index of the particular transaction
}
type output struct {
	Addr            [32]byte //Address to send the money to
	Amount          uint32 //Amount sending
	Signature       []byte //The hash of the output encrypted with the payers private key
}
type Transaction struct {
	Meta struct {
		TotalTransAmt uint32  //Total amount moving in transaction
		TimePrepared  uint64  //Time of the transaction
		Pubkey     []byte  //Payers public key
		Address    [32]byte  //Payers address
	}
	Inputs [] input `json:"Inputs"`
	Outputs [] output `json:"Outputs"`
	Hash [32]byte //Hash of the whole transaction
}

func checkerror(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
func (tx Transaction) Dump() (ret []byte) {
	tx.Hash = [32]byte{}
	ret, err := json.Marshal(tx)
	checkerror(err)
	tx.Hash = sha256.Sum256(ret)
	ret, err = json.Marshal(tx)
	checkerror(err)
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
