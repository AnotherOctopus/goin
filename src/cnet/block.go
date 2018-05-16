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
)

type Block struct{
	blocksize uint64
	header struct{
		prevBlockHash [constants.HASHSIZE]byte
		transHash []byte
		tStamp uint64
		target uint32
		noncetry uint32
	}
	transCnt uint32
	txs []Transaction
}

func (bl Block) Dump() (int, []byte){
	bBlock := make([]byte,bl.blocksize)
	binary.LittleEndian.PutUint64(bBlock[0:8], bl.blocksize)
	binary.LittleEndian.PutUint64(bBlock[8:16], bl.header.tStamp)
	binary.LittleEndian.PutUint32(bBlock[16:20], bl.header.target)
	binary.LittleEndian.PutUint32(bBlock[20:24], bl.header.noncetry)
	binary.LittleEndian.PutUint32(bBlock[24:28], bl.transCnt)
	indx := 28
	for _,tx := range bl.txs{
			txSize,txBytes := tx.Dump()
		 	copy(bBlock[indx:indx+txSize],txBytes)
		 	indx += txSize
	}
	return indx, bBlock
}

func (bl Block) Hash() ([constants.HASHSIZE] byte){
	_, blBytes := bl.Dump()
	blHash := sha256.Sum256(blBytes)
	return blHash
}

func (bl Block) HeaderSize()(int){
	size := 0
	size += 32 // prevBlockHash
	size += len(bl.header.transHash) //transHash
	size += 8 // tStamp
	size += 4 // target
	size += 4 // noncetry
	return size
}
func (bl Block) String()(string){
	retstring := ""
	retstring += "Block Made At: " + strconv.Itoa(int(bl.header.tStamp)) + "\n"
	retstring += "Previous Block: "  + hex.EncodeToString(bl.header.prevBlockHash[:]) + "\n"
	retstring += "Ease: " + strconv.Itoa(int(bl.header.target)) + "\n"
	retstring += "Transactions: \n"
	for _,tx := range bl.txs {
		retstring += "TX 1: \n"
		retstring += tx.String()
	}
	retstring += "\n"
	return retstring
}

func CreateGenesisBlock(creator * wallet.Wallet)(bl Block){
	bl.header.prevBlockHash = [constants.HASHSIZE]byte{0}
	bl.header.tStamp = uint64(100)//time.Now().Unix())
	bl.header.target = 243
	bl.header.noncetry = 0

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
	bl.header.transHash = Merkleify(bl.txs)
	bl.blocksize = uint64(totalTxSize)
	bl.transCnt = uint32(len(bl.txs))
	return
}
func (bl Block) CheckNonce(nonce uint32) (bool){
	bl.header.noncetry = nonce
	hash := bl.Hash()
	hashval := big.NewInt(0).SetBytes(hash[:])
	maxHash := big.NewInt(1)
	maxHash = maxHash.Lsh(maxHash,uint(bl.header.target))
	if hashval.Cmp(maxHash) < 0 {
		return true
	}else {
		return false
	}
}
func mine(w * wallet.Wallet){
	for i := 0; i  < math.MaxInt64; i += 1{
		log.Println("Trying ",i)
		if CreateGenesisBlock(w).CheckNonce(uint32(i)){
			log.Println(i," Success!")
			break
		}
	}
}