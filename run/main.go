package main

import (
	"encoding/base64"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/AnotherOctopus/goin/cnet"
	"github.com/AnotherOctopus/goin/wallet"
	mgo "gopkg.in/mgo.v2"
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
	return base64.StdEncoding.EncodeToString(genesisBlock.Hash[:])

}

func main() {
	creategen := flag.Bool("c", false, "generate new genesis files")
	flag.Parse()
	if *creategen {
		genHash := create()
		exec.Command("cp", "-r", "/var/lib/mongodb", ".").Run()
		ioutil.WriteFile("networkfiles/genhash", []byte(genHash), 0644)
	} else {
		hostdomain, err := os.Hostname()
		if err != nil {
			log.Println("Something up nigga")
			os.Exit(1)
		}
		peerips := []string{}
		nd := cnet.New(peerips)
		nd.RequestToJoin(hostdomain, os.Getenv("NETNODE"), os.Getenv("NETNODE") == "")
		go nd.TxListener()
		go nd.BlListener()
		go nd.CmdListener()
		go nd.ExposeInfo()
		http.ListenAndServe(":1945", nil)
		<-make(chan bool)
	}
}
