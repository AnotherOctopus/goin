//Node implements the networking between blocks. It should be used by instantiating, and then loading in a wallet
package cnet

import (
	"container/heap"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/AnotherOctopus/goin/constants"
	"github.com/AnotherOctopus/goin/wallet"
)

type Node struct {
	peers   []string         // The nodes this node gossips to
	Wallets []*wallet.Wallet // The wallets this node is using

	mineQ      *TxHeap //queue of transactions to mine
	chainHeads []Block
	isMiner    bool // Whether this node should act as a miner
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func (nd Node) GetPeers() []string {
	return nd.peers
}

//New make a node that has no wallets it is using
func New(peerips []string) (nd Node) {

	nd.peers = make([]string, len(peerips))
	for i, p := range peerips {
		nd.peers[i] = p
	}
	nd.Wallets = make([]*wallet.Wallet, 0)
	nd.mineQ = new(TxHeap)
	heap.Init(nd.mineQ)
	return
}

func (nd Node) SendBlk(blk Block) (reterr error) {
	for _, peer := range nd.peers {
		_, blkData := blk.Dump()
		raw := base64.StdEncoding.EncodeToString(blkData)
		reqstring := fmt.Sprintf("http://%v:%v/block", peer, constants.CMDRXPORT)
		_, err := http.NewRequest("PUT", reqstring, strings.NewReader(raw))
		if err != nil {
			return err
		}
	}
	return nil
}

//Send a transaction from this node
func (nd Node) SendTx(tx Transaction) (reterr error) {
	if v := verifyTx(tx); v != nil {
		log.Println(v)
		return errors.New("Transaction not valid")
	}
	SaveTx(tx)
	for _, peer := range nd.peers {
		_, txData := tx.Dump()
		raw := base64.StdEncoding.EncodeToString(txData)
		reqstring := fmt.Sprintf("http://%v:%v/transaction", peer, constants.CMDRXPORT)
		req, err := http.NewRequest("PUT", reqstring, strings.NewReader(raw))
		if err != nil {
			return err
		}
		client := &http.Client{}
		client.Do(req)
	}
	return nil
}

func (nd *Node) RequestToJoin(nodeip string, netip string, newNet bool) (err error) {
	if newNet {
		fmt.Println("Creating New Network!")
		return nil
	}
	nd.peers = append(nd.peers, netip)
	reqstring := fmt.Sprintf("http://%v:%v/join", netip, constants.CMDRXPORT)
	log.Println("DIALING: %v", reqstring, nodeip)
	client := &http.Client{}
	req, err := http.NewRequest("POST", reqstring, strings.NewReader(nodeip))
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	raw, err := ioutil.ReadAll(resp.Body)
	if raw[0] != byte('A') {
		fmt.Println("BAD RESOPONS")
		fmt.Println(string(raw), err)
		return nil
	}
	fmt.Println("JOINED NETWORK!")
	return nil
}

//Take the received transaction and save and load it
func (nd *Node) HandleTX(tx Transaction) {
	verified := verifyTx(tx)
	if verified == nil {
		// Save the transaction
		SaveTx(tx)
		heap.Push(nd.mineQ, tx)
		// If this transaction is relevant to us, save it
		for _, txOut := range tx.Outputs {
			for _, w := range nd.Wallets {
				canclaim := false
				for _, addr := range w.Address {
					if reflect.DeepEqual(txOut.Addr, addr) {
						canclaim = true
					}
				}
				if canclaim {
					w.ClaimedTxs = append(w.ClaimedTxs, tx.Hash)
				}
			}
		}
	} else {
		fmt.Println("Invalid transaction Recieved")
		fmt.Println(verified)
	}
}

//Starts mining a list of transactions. When you send a true to the kill channel, it forcefully stops the goroutine
func (nd *Node) StartMining(kill chan bool, txs [][constants.HASHSIZE]byte) {
	if len(txs) > 0 {
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
			txSize := len(tx)
			totalTxSize += txSize
		}

		totalTxSize += blk.HeaderSize()
		blk.blocksize = uint64(totalTxSize)

		for i := 0; i < math.MaxInt64; i += 1 {
			log.Println("Trying ", i)
			select {
			case <-kill:
				return
			default:
				if blk.CheckNonce(uint32(i)) {
					log.Println(i, "Success!")
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

func (nd *Node) HandleBlk(blk Block) {

}
func (nd *Node) AddPeer(peer string) {
	nd.peers = append(nd.peers, peer)
}

//Function that handles listening for incoming blocks. Handles starting and killing mining
func (nd *Node) BlListenerOLD() {
	// Listen for incoming connections.
	l, _ := net.Listen(constants.CONN_TYPE, constants.NETWORK_INT)
	// Close the listener when the application closes.
	defer l.Close()
	//Create channel to kill mining
	kill := make(chan bool)
	//Create slice of transactions to mine
	miningtxs := make([][constants.HASHSIZE]byte, 0)
	//Take whole tomine stack and load it in a list to mine
	for nd.mineQ.Len() > 0 {
		tx := heap.Pop(nd.mineQ)
		miningtxs = append(miningtxs, tx.(Transaction).Hash)
	}
	//Start mining the list
	go nd.StartMining(kill, miningtxs)
	fmt.Println("Blocks listening on " + constants.NETWORK_INT)

	//This will be the buffer to prepare new mining routines
	blkbuffer := make([]byte, constants.MAXBLKNETSIZE)
	for {
		// Listen for an incoming connection.
		conn, _ := l.Accept()
		conn.Read(blkbuffer)
		//Loading received block
		blk := LoadBlk(blkbuffer)
		//List of transactions that have already been mined
		wasmined := make([][constants.HASHSIZE]byte, 0)
		//List of transactions that need to be removed from the heap
		removefromheap := make([][constants.HASHSIZE]byte, 0)
		//If we have seen this block before we will ignore it
		if !reflect.DeepEqual(blk, getBlkFromHash(blk.Hash)) {
			//Try to verify it, ignore otherwise
			verified := verifyBlk(blk)
			if verified == nil {
				for _, blktx := range blk.Txs {
					//We know of this transaction
					isKnown := reflect.DeepEqual(Transaction{}, GetTxFromHash(blktx))
					//We are prepped to mine this transaction
					inToMine := nd.mineQ.Contains(blktx)
					if isKnown {
						for _, miningtx := range miningtxs {
							isMining := miningtx == blktx
							if isMining {
								wasmined = append(wasmined, blktx)
								break
							}
						}
					} else {
						SaveTx(requestTxn(blktx))
					}
					if inToMine {
						removefromheap = append(removefromheap, blktx)
					}
				}
				//This will become the new queue of transactions to mine
				minebuffer := new(TxHeap)
				var isin bool
				for nd.mineQ.Len() > 0 {
					//Pop off all the items in the current queue
					temptx := heap.Pop(nd.mineQ)
					//if we need to remove it from te heap
					isin = true
					for _, toremove := range removefromheap {
						if reflect.DeepEqual(temptx, toremove) {
							isin = false
							break
						}
					}
					//Push it onto the new queue
					if isin {
						heap.Push(minebuffer, temptx)
					}
				}
				// Set the new queue
				nd.mineQ = minebuffer
				//If we are currently mining something that has been mined already
				if len(wasmined) > 0 {
					kill <- true
					//start mining again
					miningtxs = make([][constants.HASHSIZE]byte, 0)
					for nd.mineQ.Len() > 0 {
						tx := heap.Pop(nd.mineQ)
						miningtxs = append(miningtxs, tx.(Transaction).Hash)
					}
					go nd.StartMining(kill, miningtxs)
				}
				// Save the Block
				SaveBlk(blk)
			} else {
				fmt.Println("Invalid Block Recieved")
				fmt.Println(verified)
			}
		}
	}
}

func requestTxn(tx [constants.HASHSIZE]byte) Transaction {
	return Transaction{}
}

func (nd *Node) MostTrustedBlock() Block {
	var mostTrusted Block
	mostTrusted.blockchainlength = 0
	for _, bl := range nd.chainHeads {
		if bl.blockchainlength > mostTrusted.blockchainlength {
			mostTrusted = bl
		}
	}
	return mostTrusted
}
