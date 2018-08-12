package main

import (
	"log"
	"net/http"
	"os"

	"github.com/AnotherOctopus/goin/cnet"
)

func exposefiles() {
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)
}

func main() {
	go exposefiles()
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
	http.ListenAndServe(":1945", nil)
	<-make(chan bool)
}
