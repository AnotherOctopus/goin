package cnet

import (
	"encoding/binary"
	"crypto/sha256"
	"time"
)

type Block struct{
	blocksize uint64
	header struct{
		prevBlockHash [32]byte
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

func (bl Block) Hash() ([32] byte){
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
func CreateGenesisBlock()(bl Block){
	bl.header.prevBlockHash = [32]byte{0}
	bl.header.tStamp = uint64(time.Now().Unix())
	bl.header.target = 250
	bl.header.noncetry = 0
	tx := new(Transaction)
	//DEFINE TRANSACTION HERE
	bl.txs[0] = *tx
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