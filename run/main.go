package main

import (
	"encoding/base64"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/AnotherOctopus/goin/cnet"
	"github.com/AnotherOctopus/goin/network"
	"github.com/AnotherOctopus/goin/wallet"
	"github.com/globalsign/mgo"
	"github.com/gorilla/handlers"
)

func create() string {
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
	ioutil.WriteFile("networkfiles/genesisWallet", w.Dump(), 0644)
	exec.Command("mongodump").Run()
	return base64.StdEncoding.EncodeToString(genesisBlock.Hash[:])

}

func main() {
	creategen := flag.Bool("c", false, "generate new genesis files")
	flag.Parse()
	if *creategen {
		genHash := create()
		ioutil.WriteFile("networkfiles/genhash", []byte(genHash), 0644)
	} else {
		peerips := []string{}
		nd := cnet.New(peerips)
		nd.RequestToJoin(os.Getenv("NETINT"), os.Getenv("NETNODE"), os.Getenv("NETNODE") == "")
		context := &network.HandlerContext{&nd}
		r := network.Handlers(context)
		loggedRouter := handlers.LoggingHandler(os.Stdout, r)
		http.Handle("/", r)
		go http.ListenAndServe(":1945", loggedRouter)
		for {
			time.Sleep(10 * time.Second)
		}
	}
}
