//Node implements the networking between blocks. It should be used by instantiating, and then loading in a wallet
package cnet

import (
	"wallet"
	"fmt"
	"net"
	"os"
	"reflect"
	"constants"
	"container/heap"
	"time"
	"math"
	"log"
)

type Node struct {
	peers   []string // The nodes this node gossips to
	Wallets []*wallet.Wallet // The wallets this node is using

	mineQ   * TxHeap
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
	nd.mineQ = new(TxHeap)
	heap.Init(nd.mineQ)
	return
}

func (nd Node) SendBlk (blk Block)(reterr error){
	for _,peer := range nd.peers{
		conn, err := net.Dial("tcp", peer)
		if err != nil {
			return err
		}
		_, blkData := blk.Dump()
		conn.Write(blkData)
		conn.Close()
	}
	return nil
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

func (nd * Node) StartMining(kill chan bool,txs [][constants.HASHSIZE]byte) {
	if len(txs) > 0{
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

		for i := 0; i  < math.MaxInt64; i += 1{
			log.Println("Trying ",i)
			select {
				case <- kill:
					return
				default:
					if blk.CheckNonce(uint32(i)){
						log.Println(i,"Success!")
						blk.Header.Noncetry = uint32(i)
						i = math.MaxInt64
					}
			}
		}

		blk.Header.Tstamp = uint64(time.Now().Unix())
		blk.SetHash(blk.Header.Noncetry)
		SaveBlk(*blk)
		nd.SendBlk(*blk)
	}
}
func (nd * Node) BlListener(){
	// Listen for incoming connections.
	l, err := net.Listen(constants.CONN_TYPE, constants.NETWORK_INT+":"+constants.BLOCKRXPORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	kill := make(chan bool)
	miningtxs := make([][constants.HASHSIZE]byte,0)
	for nd.mineQ.Len() > 0 {
		tx := heap.Pop(nd.mineQ)
		miningtxs = append(miningtxs,tx.(Transaction).Hash)
	}
	nd.StartMining(kill,miningtxs)
	fmt.Println("Blocks listening on " + constants.NETWORK_INT + ":" + constants.BLOCKRXPORT)
	blkbuffer := make([]byte,constants.MAXBLKNETSIZE)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		conn.Read(blkbuffer)
		blk:= LoadBlk(blkbuffer)
		if !reflect.DeepEqual(blk,getBlkFromHash(blk.Hash)){
			verified := verifyBlk(blk)
			if verified == nil {
				// Save the Block
				for _, miningtx := range miningtxs {
					for _, blktx := range blk.Txs {
						isKnown := reflect.DeepEqual(Transaction{},getTxFromHash(blktx))
						isMining := miningtx == blktx
						inToMine := nd.mineQ.Contains(blktx)
					}
				}
				SaveBlk(blk)
			} else {
				fmt.Println("Invalid Block Recieved")
				fmt.Println(verified)
			}
		}
	}
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
	fmt.Println("Transaction listening on " + constants.NETWORK_INT + ":" + constants.TRANSRXPORT)
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
