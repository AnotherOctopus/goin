// Transaction defines the transaction struct and the useful functions that are relevant to transactions
package cnet

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/AnotherOctopus/goin/constants"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	//"gopkg.in/mgo.v2/bson"
	"reflect"

	"github.com/AnotherOctopus/goin/wallet"
)

type Input struct {
	PrevTransHash [constants.HASHSIZE]byte //Input transaction that is being spent
	OutIdx        uint32                   //Index of the particular transaction
}

type Output struct {
	Addr      [constants.ADDRESSSIZE]byte //Address to send the money to
	Amount    uint32                      //Amount sending
	Signature []byte                      //The hash of the Output encrypted with the payers private key
}

type AnonTransaction struct {
	Inputs  []Input  `json:"Inputs"` // Inputs?
	Outputs []Output `json:"Output"` // Output?
}

type Transaction struct {
	Meta struct {
		TimePrepared int64                       //Time of the transaction
		Pubkey       rsa.PublicKey               //Payers public key
		Address      [constants.ADDRESSSIZE]byte //Payers address
	}
	Inputs  []Input                  `json:"Inputs"` // Inputs?
	Outputs []Output                 `json:"Output"` // Output?
	Hash    [constants.HASHSIZE]byte //Hash of the whole transaction
}

//Error associated with transaction
func (tx Transaction) Error() string {
	ret := "--------------------------------------------------------------\n"
	ret += "TRANSACTION Not Valid:\n"
	ret += tx.String()
	ret += "--------------------------------------------------------------\n"
	return ret
}

//transaction checkerror function
func checkerror(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

//Defines print
func (tx Transaction) String() string {
	retstring := ""
	retstring += "	TX Hash: " + base64.StdEncoding.EncodeToString(tx.Hash[:]) + "\n"
	retstring += "	Payer: " + base64.StdEncoding.EncodeToString(tx.Meta.Address[:]) + "\n"
	retstring += "	Transaction Prepared at Time: " + strconv.Itoa(int(tx.Meta.TimePrepared)) + "\n"
	for idx, inp := range tx.Inputs {
		retstring += "	Input idx " + strconv.Itoa(idx) + ": "
		retstring += base64.StdEncoding.EncodeToString(inp.PrevTransHash[:]) + ", idx " + strconv.Itoa(int(inp.OutIdx)) + "\n"
	}

	for idx, oup := range tx.Outputs {
		retstring += "	Output idx" + strconv.Itoa(idx) + ": "
		retstring += strconv.Itoa(int(oup.Amount)) + " To " + base64.StdEncoding.EncodeToString(oup.Addr[:]) + "\n"
	}
	return retstring
}

// Takes the transaction and sets the has field
func (tx *Transaction) SetHash() (err error) {
	tx.Hash = [constants.HASHSIZE]byte{}
	ret, err := json.Marshal(tx)
	if err != nil {
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
func LoadTX(b []byte) *Transaction {
	tx := new(Transaction)
	err := json.Unmarshal(b, tx)
	checkerror(err)
	return tx
}

// Load a transaction from a file
func LoadFTX(filename string) *Transaction {
	rettx := new(Transaction)
	raw, err := ioutil.ReadFile(filename)
	checkerror(err)
	err = json.Unmarshal(raw, &rettx)
	checkerror(err)
	return rettx
}

//Calculates the merkle root of a list of transactions
func Merkleify(txs [][constants.HASHSIZE]byte) []byte {
	hashlist := make([][]byte, len(txs))
	for i, tx := range txs {
		hashlist[i] = tx[:]
	}
	for len(hashlist) > 1 {
		parseIdx := 0
		hashIdx := 0
		for parseIdx < len(hashlist) {
			if parseIdx+1 != len(hashlist) {
				hashlist[hashIdx] = append(hashlist[parseIdx], hashlist[parseIdx+1]...)
				hashIdx += 1
				parseIdx += 2
			} else {
				h := sha256.New()
				copy(hashlist[hashIdx], h.Sum(hashlist[parseIdx])[:])
				hashIdx += 1
				parseIdx += 1
			}
		}
	}
	return hashlist[0]
}

// Generates the signature associated with an Output
func (o Output) GenSignature(key *rsa.PrivateKey) []byte {
	amountBytes := make([]byte, 4) // This will be a concatonation of the amount and the recipient
	binary.LittleEndian.PutUint32(amountBytes, o.Amount)
	toSign := append(amountBytes, o.Addr[:]...)
	toSignHash := sha256.Sum256(toSign)
	sig, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, toSignHash[:])
	checkerror(err)
	return sig
}

// Verifies the signature on an Output
func (o Output) VerifySignature(key *rsa.PublicKey) error {
	amountBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(amountBytes, o.Amount)
	toVerify := append(amountBytes, o.Addr[:]...)
	toVerifyHash := sha256.Sum256(toVerify)
	return rsa.VerifyPKCS1v15(key, crypto.SHA256, toVerifyHash[:], o.Signature)
}

// Returns the transaction object of an associated hash
func GetTxFromHash(hash [constants.HASHSIZE]byte) Transaction {
	var retTX Transaction
	sess, err := mgo.Dial("localhost")
	checkerror(err)
	defer sess.Close()
	handle := sess.DB("Goin").C("Transactions")
	TXQuery := handle.Find(bson.M{"hash": hash})
	err = TXQuery.One(&retTX)
	if err != nil {
		return Transaction{}
	}
	return retTX
}

// Verifies that the transaction object is valid
func verifyTx(tx Transaction) (err error) {
	// Check if the the Transaction is valid

	if reflect.DeepEqual(tx, GetTxFromHash(GenesisBlock().Txs[0])) {
		return nil
	}
	// Check whole Hash
	origHash := tx.Hash
	tx.SetHash()
	if tx.Hash != origHash {
		log.Println("There is an issue with the Hash")
		log.Println("Hash of Input", base64.StdEncoding.EncodeToString(origHash[:]))
		log.Println("Expeected Hash", base64.StdEncoding.EncodeToString(tx.Hash[:]))
		tx.Hash = origHash
		return tx
	}

	// Check address
	if wallet.Pkeytoaddress(tx.Meta.Pubkey) != tx.Meta.Address {
		log.Println("There is an issue with the Address")
		return tx
	}

	// Saving the total Input and Output
	totalOut := 0
	totalIn := 0

	// Verify the signature of each Output
	for i, outp := range tx.Outputs {
		//Check Signature of Outputs
		err = outp.VerifySignature(&tx.Meta.Pubkey)
		if err != nil {
			log.Println("There is an issue with output" + string(i))
			return err
		}
		totalOut += int(outp.Amount)
	}

	//Verify Inputs
	for i, inp := range tx.Inputs {
		prevTx := GetTxFromHash(inp.PrevTransHash)
		totalIn += int(prevTx.Outputs[inp.OutIdx].Amount)
		//Verify that the previous Output was directed to this address
		if prevTx.Outputs[inp.OutIdx].Addr != tx.Meta.Address {
			log.Println("There is an issue with input" + string(i))
			return tx
		}
		// Verify that the previous transactions that it referances are valid
		err = verifyTx(prevTx)
		if err != nil {
			log.Println("There is an issue with Transaction Input")
			return err
		}
	}

	// Verify that the total out is the total in
	if totalIn != totalOut {
		log.Println("There is an issue with Transaction Value")
		tx.Hash = origHash
		return tx
	}

	return nil
}

//Saves a transaction
func SaveTx(tx Transaction) (err error) {
	sess, err := mgo.Dial("localhost")
	checkerror(err)
	defer sess.Close()
	handle := sess.DB("Goin").C("Transactions")
	handle.Insert(tx)
	return nil
}
