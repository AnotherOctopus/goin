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
	Txs [][constants.HASHSIZE]byte // The actual transactions
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
	for _,tx := range bl.Txs{
		 	copy(bBlock[indx:indx+len(tx)],tx[:])
		 	indx += len(tx)
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
	for _,tx := range bl.Txs {
		retstring += "TX 1: \n"
		retstring += getTxFromHash(tx).String()
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
	tx.Inputs = make([]Input,1)
	tx.Outputs = make([]Output,1)
	tx.Inputs[0].OutIdx = 0
	copy(tx.Inputs[0].PrevTransHash[:], []byte("Tutturu!"))
	tx.Outputs[0].Amount = 100
	tx.Outputs[0].Addr = creator.Address[0]
	tx.Outputs[0].Signature = tx.Outputs[0].GenSignature(creator.Keys[0])
	tx.SetHash()
	SaveTx(tx)

	bl.Txs = make([][constants.HASHSIZE]byte,1)
	bl.Txs[0] = tx.Hash
	totalTxSize := 0
	for _, tx := range bl.Txs {
		txSize:= len(tx)
		totalTxSize += txSize
	}

	totalTxSize += bl.HeaderSize()
	fulltxs := make ([]Transaction,1)
	fulltxs[0] = getTxFromHash(tx.Hash)
	bl.Header.TransHash = Merkleify(fulltxs)
	bl.blocksize = uint64(totalTxSize)
	bl.transCnt = uint32(len(bl.Txs))
	log.Println(bl)
	bl.SetHash(mine(creator,bl))

	return
}
func GenesisBlock()(Block){
	sess, err := mgo.Dial("localhost")
	checkerror(err)
	defer sess.Close()
	handle := sess.DB("Goin").C("Blocks")
	genblockHash,err := base64.StdEncoding.DecodeString("AABi8dOFlRxS3FcczaW372CP6/13Hnpt145fh2FmHVo=")
	checkerror(err)
	genBlkQuery := handle.Find(bson.M{"hash":genblockHash})
	var genblk Block
	err = genBlkQuery.One(&genblk)
	checkerror(err)
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
func mine(w * wallet.Wallet,bl Block)(validNonce uint32){
	for i := 0; i  < math.MaxInt64; i += 1{
		log.Println("Trying ",i)
		if bl.CheckNonce(uint32(i)){
			log.Println(i," Success!")
			return uint32(i)
		}
	}
	return 0
}