package main

import (
	"encoding/base64"
	"io/ioutil"
	"log"

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

	w := wallet.NewWallet(1)
	genesisBlock, genesisTx := cnet.CreateGenesisBlock(w)
	cnet.SaveBlock(genesisBlock)
	cnet.SaveTx(genesisTx)
	log.Println(genesisBlock)
	ioutil.WriteFile("genesisWallet", w.Dump(), 0644)
	log.Println("GENISIS HASH")
	log.Println(base64.StdEncoding.EncodeToString(genesisBlock.Hash[:]))

}
