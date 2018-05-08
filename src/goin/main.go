package main

import (
	"node"
	"fmt"
)
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
func main(){
	nd := node.New()
	nd.TxListener()
}