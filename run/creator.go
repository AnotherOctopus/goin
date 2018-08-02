package main

import (
	"io/ioutil"

	"github.com/AnotherOctopus/goin/cnet"
	"github.com/AnotherOctopus/goin/wallet"
	mgo "gopkg.in/mgo.v2"
)

func create() {
	var v interface{}
	sess, _ := mgo.Dial("localhost")
	defer sess.Close()
	sess.DB("Goin").C("Transactions").RemoveAll(v)
	sess.DB("Goin").C("Blocks").RemoveAll(v)

	w := wallet.NewWallet(3)
	genesisBlock, genesisTx := cnet.CreateGenesisBlock(w)
	cnet.SaveBlock(genesisBlock)
	cnet.SaveTx(genesisTx)
	ioutil.WriteFile("genesisWallet", w.Dump(), 0644)

}
