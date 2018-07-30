//Node implements the networking between blocks. It should be used by instantiating, and then loading in a wallet
package cnet

import (
	"bytes"
	"container/heap"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net"
	"os"
	"reflect"
	"strconv"
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
		os.Exit(1)
	}
}

//New make a node that has no wallets it is using
func New(peerips []string) (nd Node) {

	nd.peers = make([]string, len(peerips))
	for i, p := range peerips {
		nd.peers[i] = p + ":" + constants.TRANSRXPORT
	}
	nd.Wallets = make([]*wallet.Wallet, 0)
	nd.mineQ = new(TxHeap)
	heap.Init(nd.mineQ)
	return
}

func (nd Node) SendBlk(blk Block) (reterr error) {
	for _, peer := range nd.peers {
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
func (nd Node) SendTx(tx Transaction) (reterr error) {
	for _, peer := range nd.peers {
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

func (nd *Node) RequestToJoin(nodeip string, netip string, newNet bool) (err error) {
	if newNet {
		fmt.Println("Creating New Network!")
		go nd.joinService(nodeip)
		return nil
	}
	nd.peers = append(nd.peers, netip+":"+constants.JOINPORT)
	conn, err := net.Dial("tcp", netip+":"+constants.JOINPORT)
	if err != nil {
		return err
	}
	conn.Write([]byte(nodeip))
	fmt.Println("JOINED NETWORK!")
	go nd.joinService(nodeip)
	return nil
}

func (nd *Node) joinService(ip string) {
	l, err := net.Listen(constants.CONN_TYPE, constants.NETWORK_INT+":"+constants.JOINPORT)
	tcpError(err)
	defer l.Close()
	txbuffer := make([]byte, constants.MAXTRANSNETSIZE)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		tcpError(err)
		conn.Read(txbuffer)
		nd.peers = append(nd.peers, string(txbuffer)+":"+constants.JOINPORT)
		fmt.Println("Someone Joined our network!")
	}

}

//Function that handles listening for incoming blocks. Handles starting and killing mining
func (nd *Node) BlListener() {
	// Listen for incoming connections.
	l, err := net.Listen(constants.CONN_TYPE, constants.NETWORK_INT+":"+constants.BLOCKRXPORT)
	tcpError(err)
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
	fmt.Println("Blocks listening on " + constants.NETWORK_INT + ":" + constants.BLOCKRXPORT)

	//This will be the buffer to prepare new mining routines
	blkbuffer := make([]byte, constants.MAXBLKNETSIZE)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		tcpError(err)
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
					isKnown := reflect.DeepEqual(Transaction{}, getTxFromHash(blktx))
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

func tcpError(err error) {
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
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

// Listener function for the transactions
func (nd *Node) TxListener() {
	// Listen for incoming connections.
	l, err := net.Listen(constants.CONN_TYPE, constants.NETWORK_INT+":"+constants.TRANSRXPORT)
	tcpError(err)
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Transaction listening on " + constants.NETWORK_INT + ":" + constants.TRANSRXPORT)
	txbuffer := make([]byte, constants.MAXTRANSNETSIZE)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		tcpError(err)
		conn.Read(txbuffer)
		tx := LoadTX(txbuffer)
		// Handle connections in a new goroutine.
		if !reflect.DeepEqual(getTxFromHash(tx.Hash), tx) {
			go nd.handleTX(*tx)
		}
	}
}

func (nd *Node) CmdListener() {
	l, err := net.Listen(constants.CONN_TYPE, constants.NETWORK_INT+":"+constants.CMDRXPORT)
	tcpError(err)
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Transaction listening on " + constants.NETWORK_INT + ":" + constants.CMDRXPORT)
	txbuffer := make([]byte, constants.MAXTRANSNETSIZE)
	fmt.Println("Listening for Commands...")
	done := false
	for !done {
		/*[1] Send Transaction From File")
		  [2] Manually Prepare Transaction")
		  [3] View Current Balence")
		  [4] Make A New Wallet")
		  [5] Load A Wallet")
		  [10] Exit")*/
		conn, err := l.Accept()
		tcpError(err)
		conn.Read(txbuffer)
		text := string(bytes.Trim(txbuffer, "\x00"))
		i, err := strconv.ParseInt(text, 10, 8)
		if err != nil {
			fmt.Println("INPUT IS NOT A NUMBER")
			fmt.Println(err)
			continue
		}
		switch i {
		case 1:
			if len(nd.Wallets) == 0 {
				fmt.Println("Load a Wallet First!")
				break
			}
			fmt.Println("Select a Wallet index to use")
			for i := range nd.Wallets {
				fmt.Println(fmt.Sprintf("[%v]", i))
			}
			conn, err := l.Accept()
			tcpError(err)
			conn.Read(txbuffer)
			text = string(bytes.Trim(txbuffer, "\x00"))
			windex, err := strconv.ParseInt(text, 10, 8)
			if err != nil {
				fmt.Println("INPUT IS NOT A NUMBER")
				fmt.Println(err)
				continue
			}
			w := nd.Wallets[windex]
			fmt.Println("Select an Address to use")
			for i, addr := range w.Address {
				fmt.Println(fmt.Sprintf("[%v]: %v", i, hex.EncodeToString(addr[:])))
			}
			conn, err = l.Accept()
			tcpError(err)
			conn.Read(txbuffer)
			text = string(bytes.Trim(txbuffer, "\x00"))
			addidx, err := strconv.ParseInt(text, 10, 8)
			if err != nil {
				fmt.Println("INPUT IS NOT A NUMBER")
				fmt.Println(err)
				continue
			}
			addrtouse := w.Address[addidx]
			fmt.Println("Using ", hex.EncodeToString(addrtouse[:]))
			fmt.Println("Select Filename of Transaction to send")
			conn, err = l.Accept()
			tcpError(err)
			conn.Read(txbuffer)
			filename := string(bytes.Trim(txbuffer, "\x00"))

			txtosend := LoadFTX(filename)
			txtosend.Meta.Address = addrtouse
			txtosend.Meta.Pubkey = w.Keys[addidx].PublicKey
			txtosend.Meta.TimePrepared = time.Now().Unix()
			txtosend.SetHash()

			err = nd.SendTx(*txtosend)
			CheckError(err)
			fmt.Println("Sent")
		case 2:
			savetx := new(AnonTransaction)
			savetx.Outputs = make([]Output, 0)
			savetx.Inputs = make([]Input, 0)
			fmt.Println("Select Filename of Transaction to prepare")
			conn, err = l.Accept()
			tcpError(err)
			conn.Read(txbuffer)
			filename := string(bytes.Trim(txbuffer, "\x00"))
			var inp byte
			for inp != 'd' {
				fmt.Println("Select An Option")
				fmt.Println("[i] Add Input")
				fmt.Println("[o] Add Output")
				fmt.Println("[d] Finish and Save")
				conn, err = l.Accept()
				tcpError(err)
				conn.Read(txbuffer)
				rawinp := string(bytes.Trim(txbuffer, "\x00"))
				inp = rawinp[0]
				switch inp {
				case 'i':
					var newInput Input
					var hashbuffer [constants.HASHSIZE]byte
					fmt.Println("Hash for Previous Transaction for input?")
					conn, err = l.Accept()
					tcpError(err)
					conn.Read(txbuffer)
					prev := string(bytes.Trim(txbuffer, "\x00"))
					CheckError(err)
					newHash, err := base64.StdEncoding.DecodeString(prev)
					if err != nil {
						fmt.Println("Not Valid b64")
						continue
					}
					fmt.Println("Index for Previous Transaction for input?")
					conn, err = l.Accept()
					tcpError(err)
					conn.Read(txbuffer)
					prev = string(bytes.Trim(txbuffer, "\x00"))
					CheckError(err)
					transidx, err := strconv.ParseInt(prev, 10, 8)
					if err != nil {
						fmt.Println("INPUT IS NOT A NUMBER")
						fmt.Println(err)
						continue
					}
					newInput.OutIdx = uint32(transidx)
					copy(hashbuffer[:], newHash)
					newInput.PrevTransHash = hashbuffer
					savetx.Inputs = append(savetx.Inputs, newInput)

				case 'o':
					var newOutput Output
					var hashbuffer [constants.ADDRESSSIZE]byte
					fmt.Println("Hash for Address for output?")
					conn, err = l.Accept()
					tcpError(err)
					conn.Read(txbuffer)
					out := strings.TrimSpace(string(txbuffer))
					CheckError(err)
					newHash, err := base64.StdEncoding.DecodeString(out)
					if err != nil {
						fmt.Println("Not Valid b64")
						continue
					}
					fmt.Println("Amount to send?")
					conn, err = l.Accept()
					tcpError(err)
					conn.Read(txbuffer)
					out = strings.TrimSpace(string(txbuffer))
					CheckError(err)
					amount, err := strconv.ParseInt(out, 10, 8)
					if err != nil {
						fmt.Println("INPUT IS NOT A NUMBER")
						fmt.Println(err)
						continue
					}
					newOutput.Amount = uint32(amount)
					copy(hashbuffer[:], newHash)
					newOutput.Signature = make([]byte, 0)
					newOutput.Addr = hashbuffer
					savetx.Outputs = append(savetx.Outputs, newOutput)
				case 'd':
					data, err := json.Marshal(savetx)
					CheckError(err)
					ioutil.WriteFile(filename, data, 0644)
				default:
					fmt.Println(inp, " is not a valid input")
				}
			}

		case 3:
			fmt.Println("Balence")

		case 4:
			w := wallet.NewWallet(3)
			conn, err = l.Accept()
			tcpError(err)
			conn.Read(txbuffer)
			filename := string(bytes.Trim(txbuffer, "\x00"))
			ioutil.WriteFile(filename, w.Dump(), 0644)
			nd.Wallets = append(nd.Wallets, w)

		case 5:
			fmt.Println("Select Filename of wallet")
			conn, err = l.Accept()
			tcpError(err)
			conn.Read(txbuffer)
			filename := string(bytes.Trim(txbuffer, "\x00"))
			rawdata, err := ioutil.ReadFile(filename)
			wallet.CheckError(err)
			nd.Wallets = append(nd.Wallets, wallet.LoadWallet(rawdata))

		case 10:
			fmt.Println("Exiting")
			done = true
			break

		default:
			fmt.Println("Select a valid number")
		}
	}
}
