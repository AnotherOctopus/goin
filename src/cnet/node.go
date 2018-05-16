package cnet

import (
	"wallet"
	"fmt"
	"net"
	"os"
	"reflect"
	"constants"
	"encoding/json"
	"crypto/sha256"
)

type Node struct {
	peers   []string
	wallet *wallet.Wallet
	isMiner bool
}

func New(peerips []string)(nd Node){

	nd.peers =	make([]string,len(peerips))
	for i, p := range peerips {
		nd.peers[i] = p + ":" + constants.TRANSBROADPORT
	}
	nd.wallet = new(wallet.Wallet)
	return
}

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
func getTxFromHash([constants.HASHSIZE] byte)(tx Transaction) {

}
func verifyTx(tx Transaction)(err error) {
	// Check if the the Transaction is valid
	// Check whole Hash
	origHash := tx.Hash
	tx.SetHash()
	if tx.Hash != origHash {
		return tx
	}
	totalOut := 0
	totalIn := 0
	for _, outp := range tx.Outputs {
		//Check Signature of outputs
	}


	for _, inp := range tx.Inputs {
		prevTx := getTxFromHash(inp.PrevTransHash)
		totalIn += int(prevTx.Outputs[inp.OutIdx].Amount)
		verifyTx(prevTx)
	}
	if totalIn != totalOut {
		return tx
	}
	return nil
}

func saveTx(tx Transaction)(err error){

	return nil
}

func (nd *Node) handleTX(tx Transaction) {
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
			nd.wallet.ClaimedTxs = append(nd.wallet.ClaimedTxs, tx.Hash)
		}
	}else {
		fmt.Println("Invalid transaction Recieved")
	}
}

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
