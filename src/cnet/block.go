package cnet

type Block struct{
	blocksize uint64
	header struct{
		prevBlockHash [32]byte
		transHash [32]byte
		tStamp uint32
	}
	transCnt uint32
	txs []Transaction
}