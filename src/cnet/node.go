//Node implements the networking between blocks. It should be used by instantiating, and then loading in a wallet
package cnet

import (
	"wallet"
	"fmt"
	"net"
	"os"
	"reflect"
	"constants"
	"txheap"
	"container/heap"
	"time"
)

type Node struct {
	peers   []string // The nodes this node gossips to
	Wallets []*wallet.Wallet // The wallets this node is using

	mineQ   * txheap.TxHeap
	chainHeads []Block
	isMiner bool // Whether this node should act as a miner
}

//New make a node that has no wallets it is using
func New(peerips []string)(nd Node){


	nd.peers =	make([]string,len(peerips))
	for i, p := range peerips {
		nd.peers[i] = p + ":" + constants.TRANSRXPORT
	}
	nd.Wallets = make([]*wallet.Wallet,0)
	nd.mineQ = new(txheap.TxHeap)
	heap.Init(nd.mineQ)
	return
}

//Send a transaction from this node
func (nd Node) SendTx (tx Transaction)(reterr error){
	for _,peer := range nd.peers{
		conn, err := net.Dial("tcp", peer)
		if err != nil {
			return err
		}
		_, txData := tx.Dump()
		conn.Write(txData)
		conn.Close()
	}
	return nil
}

//Take the received transaction and save and load it
func (nd *Node) handleTX(tx Transaction) {
	verified := verifyTx(tx)
	if verified == nil {
		// Save the transaction
		SaveTx(tx)
		heap.Push(nd.mineQ,tx)
		// If this transaction is relevant to us, save it
		for _, txOut := range tx.Outputs {
			for _, w := range nd.Wallets{
				canclaim := false
				for _, addr := range w.Address{
					if reflect.DeepEqual(txOut.Addr, addr){
						canclaim = true
					}
				}
				if canclaim {
					w.ClaimedTxs = append(w.ClaimedTxs, tx.Hash)
				}
			}
		}
	}else {
		fmt.Println("Invalid transaction Recieved")
		fmt.Println(verified)
	}
}

func (nd * Node) BlListener(){
	tomine := make([][constants.HASHSIZE]byte,0)
	for nd.mineQ.Len() > 0 {
		tx := heap.Pop(nd.mineQ)
		tomine = append(tomine,tx.(Transaction).Hash)
	}
	go nd.Mine(tomine)
}

func (nd * Node) MostTrustedBlock()(Block){
	var mostTrusted Block
	mostTrusted.blockchainlength = 0
	for _,bl := range nd.chainHeads {
		if bl.blockchainlength > mostTrusted.blockchainlength{
			mostTrusted = bl
		}
	}
	return mostTrusted
}
func (nd * Node) Mine (txs [][constants.HASHSIZE]byte) (*Block){
	blk := new(Block)
	prevBlock := nd.MostTrustedBlock()
	blk.blockchainlength = prevBlock.blockchainlength + 1
	blk.Header.PrevBlockHash = prevBlock.Hash
	blk.Header.Target = prevBlock.Header.Target
	blk.Header.TransHash = Merkleify(txs)
	blk.transCnt = uint32(len(txs))
	blk.Txs = txs

	totalTxSize := 0
	for _, tx := range blk.Txs {
		txSize:= len(tx)
		totalTxSize += txSize
	}

	totalTxSize += blk.HeaderSize()
	blk.blocksize = uint64(totalTxSize)

	blk.Header.Noncetry = mine(nd.Wallets[0],*blk)
	blk.Header.Tstamp = uint64(time.Now().Unix())
	blk.SetHash(blk.Header.Noncetry)
	return blk
}

// Listener function for the transactions
func (nd * Node) TxListener() {
	// Listen for incoming connections.
	l, err := net.Listen(constants.CONN_TYPE, constants.NETWORK_INT+":"+constants.TRANSRXPORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + constants.NETWORK_INT + ":" + constants.TRANSRXPORT)
	txbuffer := make([]byte,constants.MAXTRANSNETSIZE)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		conn.Read(txbuffer)
		tx:= LoadTX(txbuffer)
		// Handle connections in a new goroutine.
		if !reflect.DeepEqual(getTxFromHash(tx.Hash),tx){
			go nd.handleTX(*tx)
		}
	}
}
