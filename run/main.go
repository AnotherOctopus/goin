package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/AnotherOctopus/goin/cnet"
)

func main() {
	creategen := flag.Bool("c", false, "generate new genesis files")
	flag.Parse()
	if *creategen {
		genHash := create()
		exec.Command("mongodump", "dump").Run()
		ioutil.WriteFile("genhash", []byte(genHash), 0644)
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
