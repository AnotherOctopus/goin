//Block implements the block and the block structure
package cnet

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/AnotherOctopus/goin/constants"
	"github.com/AnotherOctopus/goin/wallet"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const hashfile string = "networkfiles/genhash"

type Block struct {
	blocksize        uint64 // How many bytes are in a block
	blockchainlength uint64
	Header           struct {
		PrevBlockHash [constants.HASHSIZE]byte // The hash of the previous block
		TransHash     []byte                   // The Merkle root of all the transactions
		Tstamp        uint64                   // When this block was mined
		Target        uint32                   // The ease of this block
		Noncetry      uint32                   // Nonce for hash tests
	}
	Hash     [constants.HASHSIZE]byte
	transCnt uint32                     // How many transactions are in this block
	Txs      [][constants.HASHSIZE]byte // The actual transactions
}

func getgenHash() string {
	raw, err := ioutil.ReadFile(hashfile)
	checkerror(err)
	return string(raw)
}
func LoadBlk(b []byte) (block Block) {
	block.blocksize = binary.LittleEndian.Uint64(b[0:8])
	block.Header.Tstamp = binary.LittleEndian.Uint64(b[8:16])
	block.Header.Target = binary.LittleEndian.Uint32(b[16:20])
	block.Header.Noncetry = binary.LittleEndian.Uint32(b[20:24])
	block.transCnt = binary.LittleEndian.Uint32(b[24:28])
	block.blockchainlength = binary.LittleEndian.Uint64(b[28:36])
	indx := 36
	var hashbuff [constants.HASHSIZE]byte
	block.Txs = make([][constants.HASHSIZE]byte, 0)
	for indx < int(block.blocksize) {
		copy(hashbuff[:], b[indx:indx+constants.HASHSIZE])
		block.Txs = append(block.Txs, hashbuff)
		indx += constants.HASHSIZE
	}
	return
}

// Serializes the block
func (bl Block) Dump() (int, []byte) {
	bBlock := make([]byte, bl.blocksize)
	binary.LittleEndian.PutUint64(bBlock[0:8], bl.blocksize)
	binary.LittleEndian.PutUint64(bBlock[8:16], bl.Header.Tstamp)
	binary.LittleEndian.PutUint32(bBlock[16:20], bl.Header.Target)
	binary.LittleEndian.PutUint32(bBlock[20:24], bl.Header.Noncetry)
	binary.LittleEndian.PutUint32(bBlock[24:28], bl.transCnt)
	binary.LittleEndian.PutUint64(bBlock[28:36], bl.blockchainlength)
	indx := 36
	for _, tx := range bl.Txs {
		copy(bBlock[indx:indx+len(tx)], tx[:])
		indx += len(tx)
	}
	return indx, bBlock
}

// Generates the hash of the whole block
func (bl Block) HashBlock() [constants.HASHSIZE]byte {
	_, blBytes := bl.Dump()
	blHash := sha256.Sum256(blBytes)
	return blHash
}

// Number of bytes in the Header
func (bl Block) HeaderSize() int {
	size := 0
	size += 32                       // PrevBlockHash
	size += len(bl.Header.TransHash) //TransHash
	size += 8                        // Tstamp
	size += 4                        // Target
	size += 4                        // Noncetry
	return size
}

//For printing uses
func (bl Block) String() string {
	retstring := "\n"
	retstring += "---------------------------------------------------------------------------\n"
	retstring += "Block Made At: " + strconv.Itoa(int(bl.Header.Tstamp)) + "\n"
	retstring += "Previous Block: " + hex.EncodeToString(bl.Header.PrevBlockHash[:]) + "\n"
	retstring += "Ease: " + strconv.Itoa(int(bl.Header.Target)) + "\n"
	retstring += "Transactions: \n"
	for _, tx := range bl.Txs {
		retstring += "TX 1: \n"
		retstring += GetTxFromHash(tx).String()
	}
	retstring += "\n"
	retstring += "---------------------------------------------------------------------------\n"
	return retstring
}

// How the genesis block is defined
func CreateGenesisBlock(creator *wallet.Wallet) (bl Block, tx Transaction) {
	bl.Header.PrevBlockHash = [constants.HASHSIZE]byte{0}
	bl.Header.Tstamp = uint64(time.Now().Unix())
	bl.Header.Target = 230
	bl.Header.Noncetry = 0
	bl.blockchainlength = 0

	tx.Meta.TimePrepared = int64(100) //time.Now().Unix()
	tx.Meta.Pubkey = creator.Keys[0].PublicKey
	tx.Meta.Address = creator.Address[0]
	tx.Inputs = make([]Input, 1)
	tx.Outputs = make([]Output, 1)
	tx.Inputs[0].OutIdx = 0
	copy(tx.Inputs[0].PrevTransHash[:], []byte("Tutturu!"))
	tx.Outputs[0].Amount = 100
	tx.Outputs[0].Addr = creator.Address[0]
	tx.Outputs[0].Signature = tx.Outputs[0].GenSignature(creator.Keys[0])
	tx.SetHash()

	bl.Txs = make([][constants.HASHSIZE]byte, 1)
	bl.Txs[0] = tx.Hash
	totalTxSize := 0
	for _, tx := range bl.Txs {
		txSize := len(tx)
		totalTxSize += txSize
	}

	totalTxSize += bl.HeaderSize()
	fulltxs := [][constants.HASHSIZE]byte{tx.Hash}
	bl.Header.TransHash = Merkleify(fulltxs)
	bl.blocksize = uint64(totalTxSize)
	bl.transCnt = uint32(len(bl.Txs))
	bl.SetHash(mine(creator, bl))
	creator.ClaimedTxs = append(creator.ClaimedTxs, tx.Hash)

	return
}
func GenesisBlock() Block {
	sess, err := mgo.Dial("localhost")
	if err != nil {
		log.Println("CANT FIND GENISIS BLOCK (can't dial)")
		os.Exit(1)
	}
	defer sess.Close()
	handle := sess.DB("Goin").C("Blocks")
	genblockHash, err := base64.StdEncoding.DecodeString(getgenHash())
	if err != nil {
		log.Println("CANT FIND GENISIS BLOCK(cant decode hash)")
		os.Exit(1)
	}
	genBlkQuery := handle.Find(bson.M{"hash": genblockHash})
	var genblk Block
	err = genBlkQuery.One(&genblk)
	if err != nil {
		log.Println("CANT FIND GENISIS BLOCK(cant find one)")
		os.Exit(1)
	}
	return genblk
}

// Checks the block for a working nonce
func (bl Block) CheckNonce(nonce uint32) bool {
	bl.Header.Noncetry = nonce
	hash := bl.HashBlock()
	hashval := big.NewInt(0).SetBytes(hash[:])
	maxHash := big.NewInt(1)
	maxHash = maxHash.Lsh(maxHash, uint(bl.Header.Target))
	if hashval.Cmp(maxHash) < 0 {
		return true
	} else {
		return false
	}
}

func SaveBlock(bl Block) error {
	sess, err := mgo.Dial("localhost")
	checkerror(err)
	defer sess.Close()
	handle := sess.DB("Goin").C("Blocks")
	handle.Insert(bl)
	return nil
}

func (bl *Block) SetHash(nonce uint32) {
	bl.Header.Noncetry = nonce
	bl.Hash = bl.HashBlock()
	return
}

// Returns the transaction object of an associated hash
func getBlkFromHash(hash [constants.HASHSIZE]byte) Block {
	var retBlk Block
	sess, err := mgo.Dial("localhost")
	checkerror(err)
	defer sess.Close()
	handle := sess.DB("Goin").C("Blocks")
	BlkQuery := handle.Find(bson.M{"hash": hash})
	err = BlkQuery.One(&retBlk)
	if err != nil {
		return Block{}
	}
	return retBlk
}

func verifyBlk(blk Block) (err error) {
	return nil
}

func SaveBlk(blk Block) (err error) {
	sess, err := mgo.Dial("localhost")
	checkerror(err)
	defer sess.Close()
	handle := sess.DB("Goin").C("Blocks")
	handle.Insert(blk)
	return nil
}

// Runs the check nonce over and over
func mine(w *wallet.Wallet, bl Block) (validNonce uint32) {
	for i := 0; i < math.MaxInt64; i += 1 {
		if bl.CheckNonce(uint32(i)) {
			log.Println(i, " Success!")
			return uint32(i)
		}
	}
	return 0
}
