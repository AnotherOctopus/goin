package cnet

import "constants"

type TxHeap [] Transaction

func (txh TxHeap) Contains(hash [constants.HASHSIZE]byte) bool{
	for _, tx := range txh{
		if tx.Hash == hash {
			return true
		}
	}
	return false
}
func (txh TxHeap) Len() int {
	return len(txh)
}

func (txh TxHeap) Less (i,j int) bool {
	return txh[i].Meta.TimePrepared < txh[j].Meta.TimePrepared
}

func (txh TxHeap) Swap (i,j int){
	txh[i],txh[j] = txh[j],txh[i]
}

func (txh * TxHeap) Push( x interface{}){
	*txh = append(*txh,x.(Transaction))
}


func (txh * TxHeap) Pop() (x interface{}) {
	old := *txh
	n := len(old)
	x = old[n-1]
	*txh = old[0:n-1]
	return x
}
