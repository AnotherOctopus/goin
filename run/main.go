package main

import (
	"github.com/AnotherOctopus/goin/cnet"
	"log"
	"os"
)
func main(){
	hostdomain,err := os.Hostname()
	if err != nil {
		log.Println("Something up nigga")
		os.Exit(1)
	}
	peerips := []string{}
	nd := cnet.New(peerips)
	nd.RequestToJoin(hostdomain,os.Getenv("NETNODE"),os.Getenv("NETNODE")=="")
	go nd.TxListener()
	go nd.BlListener()
	go nd.CmdListener()
	<-make(chan bool)
}
