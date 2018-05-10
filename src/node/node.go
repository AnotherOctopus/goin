package node

import (
	"cnet"
	"fmt"
	"net"
	"os"
	"reflect"
)

const (
	TRANSBROADPORT = "1917"
	TRANSRXPORT    = "1943"
	BLOCKRXPORT    = "1918"
	CONN_TYPE      = "tcp"
	NETWORK_INT    = "0.0.0.0"
)


type Node struct {
	peers   []string
	wallet *Wallet
}

func New(peerips []string)(nd Node){

	nd.peers =	make([]string,len(peerips))
	for i, p := range peerips {
		nd.peers[i] = p + ":" + TRANSBROADPORT
	}
	nd.wallet = new(Wallet)
	return
}

func (nd Node) SendTx (tx cnet.Transaction)(reterr error){
	for _,peer := range nd.peers{
		conn, err := net.Dial("tcp", peer)
		if err != nil {
			return err
		}
		conn.Write(tx.Dump())
		conn.Close()
	}
	return nil
}

func verifyTx(tx cnet.Transaction)(err error) {
	// Check if the the Transaction is valid
	return nil
}

func saveTx(tx cnet.Transaction)(err error){

	return nil
}

func (nd *Node) handleTX(tx cnet.Transaction) {
	if verifyTx(tx) == nil {
		canclaim := false
		saveTx(tx)
		for _, txOut := range tx.Outputs {
			for _, addr := range nd.wallet.Address {
				if reflect.DeepEqual(txOut.Addr, addr) {
					canclaim = true
				}
			}
		}
		if canclaim {
			nd.wallet.ClaimedTxs = append(nd.wallet.ClaimedTxs, tx)
		}
	}else {
		fmt.Println("Invalid transaction Recieved")
	}
}

func (nd *Node) TxListener() {
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, NETWORK_INT+":"+TRANSRXPORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + NETWORK_INT + ":" + TRANSRXPORT)
	txbuffer := make([]byte,cnet.MAXTRANSNETSIZE)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		conn.Read(txbuffer)
		tx:= cnet.LoadTX(txbuffer)
		// Handle connections in a new goroutine.
		go nd.handleTX(tx)
	}
}
