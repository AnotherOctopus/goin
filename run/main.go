package main

import (
	"github.com/AnotherOctopus/goin/cnet"
)
func main(){
	//test()
	peerips := []string{}
	nd := cnet.New(peerips)
	nd.requestToJoin("192.168.1.127",nil,true)
	go nd.TxListener()
	go nd.BlListener()
	go nd.CmdListener()
	<-make(chan bool)
}
