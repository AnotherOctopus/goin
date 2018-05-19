//Block implements the block and the block structure
package cnet

import (
	"encoding/binary"
	"crypto/sha256"
	"wallet"
	"strconv"
	"encoding/hex"
	"constants"
	"math/big"
	"log"
	"math"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"encoding/base64"
	"fmt"
)

type Block struct{
	blocksize uint64 // How many bytes are in a block
	Header struct{
		PrevBlockHash [constants.HASHSIZE]byte // The hash of the previous block
		TransHash []byte // The Merkle root of all the transactions
		Tstamp uint64 // When this block was mined
		Target uint32 // The ease of this block
		Noncetry uint32 // Nonce for hash tests
	}
	Hash [constants.HASHSIZE]byte
	transCnt uint32 // How many transactions are in this block
	txs []Transaction // The actual transactions
}

// Serializes the block
func (bl Block) Dump() (int, []byte){
	bBlock := make([]byte,bl.blocksize)
	binary.LittleEndian.PutUint64(bBlock[0:8], bl.blocksize)
	binary.LittleEndian.PutUint64(bBlock[8:16], bl.Header.Tstamp)
	binary.LittleEndian.PutUint32(bBlock[16:20], bl.Header.Target)
	binary.LittleEndian.PutUint32(bBlock[20:24], bl.Header.Noncetry)
	binary.LittleEndian.PutUint32(bBlock[24:28], bl.transCnt)
	indx := 28
	for _,tx := range bl.txs{
			txSize,txBytes := tx.Dump()
		 	copy(bBlock[indx:indx+txSize],txBytes)
		 	indx += txSize
	}
	return indx, bBlock
}

// Generates the hash of the whole block
func (bl Block) HashBlock() ([constants.HASHSIZE] byte){
	_, blBytes := bl.Dump()
	blHash := sha256.Sum256(blBytes)
	return blHash
}

// Number of bytes in the Header
func (bl Block) HeaderSize()(int){
	size := 0
	size += 32 // PrevBlockHash
	size += len(bl.Header.TransHash) //TransHash
	size += 8 // Tstamp
	size += 4 // Target
	size += 4 // Noncetry
	return size
}

//For printing uses
func (bl Block) String()(string){
	retstring := ""
	retstring += "Block Made At: " + strconv.Itoa(int(bl.Header.Tstamp)) + "\n"
	retstring += "Previous Block: "  + hex.EncodeToString(bl.Header.PrevBlockHash[:]) + "\n"
	retstring += "Ease: " + strconv.Itoa(int(bl.Header.Target)) + "\n"
	retstring += "Transactions: \n"
	for _,tx := range bl.txs {
		retstring += "TX 1: \n"
		retstring += tx.String()
	}
	retstring += "\n"
	return retstring
}

// How the genesis block is defined
func CreateGenesisBlock(creator * wallet.Wallet)(bl Block){
	bl.Header.PrevBlockHash = [constants.HASHSIZE]byte{0}
	bl.Header.Tstamp = uint64(100)//time.Now().Unix())
	bl.Header.Target = 243
	bl.Header.Noncetry = 0

	var tx Transaction
	tx.Meta.TimePrepared = int64(100)//time.Now().Unix()
	tx.Meta.Pubkey = creator.Keys[0].PublicKey
	tx.Meta.Address = creator.Address[0]
	tx.Inputs = make([]input,1)
	tx.Outputs = make([]output,1)
	tx.Inputs[0].OutIdx = 0
	copy(tx.Inputs[0].PrevTransHash[:], []byte("Tutturu!"))
	tx.Outputs[0].Amount = 100
	tx.Outputs[0].Addr = creator.Address[0]
	tx.Outputs[0].Signature = tx.Outputs[0].GenSignature(creator.Keys[0])
	tx.SetHash()

	bl.txs = make([]Transaction,1)
	bl.txs[0] = tx
	totalTxSize := 0
	for _, tx := range bl.txs {
		txSize, _ := tx.Dump()
		totalTxSize += txSize
	}

	totalTxSize += bl.HeaderSize()
	bl.Header.TransHash = Merkleify(bl.txs)
	bl.blocksize = uint64(totalTxSize)
	bl.transCnt = uint32(len(bl.txs))
	return
}
func GenesisBlock()(Block){
	sess, err := mgo.Dial("localhost")
	checkerror(err)
	defer sess.Close()
	handle := sess.DB("Goin").C("Blocks")
	genblockHash,err := base64.StdEncoding.DecodeString("AAXK33DemUW0nQGpfu3SRuOBkIfp1hdsCk3Qgq8mzl0=")
	checkerror(err)
	genBlkQuery := handle.Find(bson.M{"hash":genblockHash})
	var genblk Block
	err = genBlkQuery.One(&genblk)
	checkerror(err)
	fmt.Print(genblk)
	return genblk
}
// Checks the block for a working nonce
func (bl Block) CheckNonce(nonce uint32) (bool){
	bl.Header.Noncetry = nonce
	hash := bl.HashBlock()
	hashval := big.NewInt(0).SetBytes(hash[:])
	maxHash := big.NewInt(1)
	maxHash = maxHash.Lsh(maxHash,uint(bl.Header.Target))
	if hashval.Cmp(maxHash) < 0 {
		return true
	}else {
		return false
	}
}

func SaveBlock(bl Block)(error){
	sess, err := mgo.Dial("localhost")
	checkerror(err)
	defer sess.Close()
	handle := sess.DB("Goin").C("Blocks")
	handle.Insert(bl)
	return nil
}

func (bl * Block) SetHash(nonce uint32) {
	bl.Header.Noncetry = nonce
	bl.Hash = bl.HashBlock()
	return
}
// Runs the check nonce over and over
func mine(w * wallet.Wallet){
	for i := 0; i  < math.MaxInt64; i += 1{
		log.Println("Trying ",i)
		if CreateGenesisBlock(w).CheckNonce(uint32(i)){
			log.Println(i," Success!")
			break
		}
	}
}