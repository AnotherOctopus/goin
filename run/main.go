package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/AnotherOctopus/goin/cnet"
	"github.com/AnotherOctopus/goin/network"
	"github.com/AnotherOctopus/goin/wallet"
	"github.com/gorilla/handlers"
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
			log.Println(os.Getenv("NETINT"), nd.GetPeers())
			out, _ := exec.Command("ping", "172.18.0.3", "-c 5", "-i 3", "-w 10").Output()
			if strings.Contains(string(out), "Destination Host Unreachable") {
				fmt.Println("TANGO DOWN")
			} else {
				fmt.Println("IT'S ALIVEEE")
				req, err := http.Get("http://172.18.0.3:1945/addresses")
				if err != nil {
					log.Println(err)
				}
				body, err := ioutil.ReadAll(req.Body)
				log.Println("REQ BODY", req.StatusCode, body)
				req.Body.Close()
			}
			time.Sleep(10 * time.Second)
		}
	}
}
