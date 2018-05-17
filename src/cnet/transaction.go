// Transaction defines the transaction struct and the useful functions that are relevant to transactions
package cnet

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"constants"
	"crypto/rsa"
	"encoding/binary"
	"crypto"
	"strconv"
	"encoding/hex"
	"wallet"
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
		TimePrepared  int64  //Time of the transaction
		Pubkey     rsa.PublicKey  //Payers public key
		Address    [constants.ADDRESSSIZE]byte  //Payers address
	}
	Inputs [] input `json:"Inputs"` // Inputs?
	Outputs [] output `json:"Outputs"` // Outputs?
	Hash [constants.HASHSIZE]byte //Hash of the whole transaction
}

//Error associated with transaction
func (tx Transaction) Error() string {
	return "TRANSACTION Not Valid:" + hex.EncodeToString(tx.Hash[:])
}

//transaction checkerror function
func checkerror(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

//Defines print
func (tx Transaction) String() (string) {
	retstring := ""
	retstring += "	TX Hash: " + hex.EncodeToString(tx.Hash[:]) + "\n"
	retstring += "	Payer: " + hex.EncodeToString(tx.Meta.Address[:]) + "\n"
	retstring += "	Transaction Prepared at Time: " + strconv.Itoa(int(tx.Meta.TimePrepared)) + "\n"
	for idx, inp := range tx.Inputs{
		retstring += "	Input idx " + strconv.Itoa(idx) + ": "
		retstring += hex.EncodeToString(inp.PrevTransHash[:]) + ", idx " + strconv.Itoa(int(inp.OutIdx)) + "\n"
	}

	for idx, oup := range tx.Outputs{
		retstring += "	Output idx" + strconv.Itoa(idx) + ": "
		retstring += strconv.Itoa(int(oup.Amount)) + " To " + hex.EncodeToString(oup.Addr[:]) + "\n"
	}
	return retstring
}

// Takes the transaction and sets the has field
func (tx * Transaction) SetHash() (err error){
	tx.Hash = [constants.HASHSIZE]byte{}
	ret, err := json.Marshal(tx)
	if err != nil{
		return err
	}
	tx.Hash = sha256.Sum256(ret)
	return nil
}

// Serialize a transaction
func (tx Transaction) Dump() (size int, ret []byte) {
	tx.SetHash()
	ret, err := json.Marshal(tx)
	checkerror(err)
	size = len(ret)
	return
}

// Load a serialized transaction
func LoadTX(b []byte) (tx Transaction) {
	err := json.Unmarshal(b, tx)
	checkerror(err)
	return
}

// Load a transaction from a file
func LoadFTX(filename string)(rettx Transaction){
	var pretx Transaction
	raw,err := ioutil.ReadFile(filename)
	checkerror(err)
	err = json.Unmarshal(raw,&pretx)
	checkerror(err)
	return
}

//Calculates the merkle root of a list of transactions
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

// Generates the signature associated with an output
func (o output) GenSignature(key * rsa.PrivateKey)([]byte){
	amountBytes := make([]byte,4) // This will be a concatonation of the amount and the recipient
	binary.LittleEndian.PutUint32(amountBytes,o.Amount)
	toSign := append(amountBytes,o.Addr[:]...)
	sig,err := rsa.SignPKCS1v15(rand.Reader,key,crypto.SHA256,toSign)
	checkerror(err)
	return sig
}

// Verifies the signature on an output
func (o output) VerifySignature(key * rsa.PublicKey)(error){
	amountBytes := make([]byte,4)
	binary.LittleEndian.PutUint32(amountBytes,o.Amount)
	toVerify := append(amountBytes,o.Addr[:]...)
	return rsa.VerifyPKCS1v15(key,crypto.SHA256,toVerify,o.Signature)
}

// Returns the transaction object of an associated hash
func getTxFromHash([constants.HASHSIZE] byte)(tx Transaction) {
	return  tx
}

// Verifies that the transaction object is valid
func verifyTx(tx Transaction)(err error) {
	// Check if the the Transaction is valid

	// Check whole Hash
	origHash := tx.Hash
	tx.SetHash()
	if tx.Hash != origHash {
		return tx
	}

	// Check address
	if wallet.Pkeytoaddress(tx.Meta.Pubkey) != tx.Meta.Address {
		return tx
	}

	// Saving the total input and output
	totalOut := 0
	totalIn := 0

	// Verify the signature of each output
	for _, outp := range tx.Outputs {
		//Check Signature of outputs
		err = outp.VerifySignature(&tx.Meta.Pubkey)
		if err != nil {
			return err
		}
	}

	// Verify that the previous transactions that it referances are valid
	for _, inp := range tx.Inputs {
		prevTx := getTxFromHash(inp.PrevTransHash)
		totalIn += int(prevTx.Outputs[inp.OutIdx].Amount)
		err = verifyTx(prevTx)
		if err != nil {
			return err
		}
	}

	// Verify that the total out is the total in
	if totalIn != totalOut {
		return tx
	}

	return nil
}

//Saves a transaction
func saveTx(tx Transaction)(err error){

	return nil
}