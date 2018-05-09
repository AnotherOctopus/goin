package node

import (
	"cnet"
	"fmt"
	"net"
	"os"
)

const (
	TRANSBROADPORT = "1917"
	TRANSRXPORT    = "1943"
	BLOCKRXPORT    = "1918"
	CONN_TYPE      = "tcp"
	NETWORK_INT    = "0.0.0.0"
)


type Node struct {
	peers   []*net.TCPAddr
	wallets []*Wallet
}

func New()(nd Node){
	nd.peers =	nil
	nd.wallets = make([]*Wallet,10)
	nd.wallets[0] = NewWallet()
	return
}

func (wal *Wallet) ClaimTx(tx cnet.Transaction) (reterr error) {

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
		saveTx(tx)
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
