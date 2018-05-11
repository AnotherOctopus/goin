package main

import (
	"node"
	"fmt"
	"bufio"
	"os"
	"strconv"
	"strings"
	"cnet"
)
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}
func main(){
	peerips := []string{"127.0.0.1"}
	nd := node.New(peerips)
	go nd.TxListener()
	fmt.Println("Listening for Transactions...")
	fmt.Println("Launching Goin CLI! ")
	done := false
	for !done{
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter Command: ")
		fmt.Println("[1] Send Transaction From File")
		fmt.Println("[2] Manually Prepare Transaction")
		fmt.Println("[3] View Current Balence")
		fmt.Println("[4] Exit")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		i,err := strconv.ParseInt(text,10,8)
		if err != nil {
			fmt.Println("INPUT IS NOT A NUMBER")
			fmt.Println(err)
			continue
		}
		switch i {
		case 1:
			fmt.Println("Select Filename of Transaction to send")
			filename, _ := reader.ReadString('\n')
			filename = strings.TrimSpace(filename)
			txtosend := cnet.LoadFTX(filename)
			err = nd.SendTx(txtosend)
			CheckError(err)
			fmt.Println("Sent")
		case 2:
			fmt.Println("Prep")
		case 3:
			fmt.Println("Balence")
		case 4:
			fmt.Println("Exiting")
			done = true
			break
		default:
			fmt.Println("Select a valid number")
		}
	}
}