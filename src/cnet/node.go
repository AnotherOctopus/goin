//Node implements the networking between blocks. It should be used by instantiating, and then loading in a wallet
package cnet

import (
	"wallet"
	"fmt"
	"net"
	"os"
	"reflect"
	"constants"
)

type Node struct {
	peers   []string // The nodes this node gossips to
	Wallets []*wallet.Wallet // The wallets this node is using
	isMiner bool // Whether this node should act as a miner
}

//New make a node that has no wallets it is using
func New(peerips []string)(nd Node){

	nd.peers =	make([]string,len(peerips))
	for i, p := range peerips {
		nd.peers[i] = p + ":" + constants.TRANSBROADPORT
	}
	nd.Wallets = make([]*wallet.Wallet,0)
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
	if verifyTx(tx) == nil {
		// Save the transaction
		saveTx(tx)
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
	}
}

// Listener function for the transactions
func (nd *Node) TxListener() {
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
		go nd.handleTX(tx)
	}
}
