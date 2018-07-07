package main

import (
	"github.com/AnotherOctopus/goin/cnet"
)
func main(){
	//test()
	peerips := []string{"127.0.0.1"}
	nd := cnet.New(peerips)
	go nd.TxListener()
	go nd.BlListener()
	go nd.CmdListener()
	<-make(chan bool)
}
