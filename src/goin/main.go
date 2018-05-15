package main

import (
	"wallet"
	"fmt"
	"bufio"
	"os"
	"strconv"
	"strings"
	"cnet"
	"io/ioutil"
	"log"
)
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}
func main(){
	peerips := []string{"127.0.0.1"}
	nd := wallet.New(peerips)
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
		fmt.Println("[4] Make A New Wallet")
		fmt.Println("[5] Load A Wallet")
		fmt.Println("[10] Exit")
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
			wallet := wallet.NewWallet(3)
			fmt.Println("Select Filename of where to save wallet")
			filename, _ := reader.ReadString('\n')
			filename = strings.TrimSpace(filename)
			ioutil.WriteFile(filename,wallet.Dump(),0644)
		case 5:
			fmt.Println("Select Filename of wallet")
			filename, _ := reader.ReadString('\n')
			filename = strings.TrimSpace(filename)
			rawdata,err := ioutil.ReadFile(filename)
			wallet.CheckError(err)
			wallet := wallet.LoadWallet(rawdata)
			log.Println(wallet)
			log.Println(cnet.CreateGenesisBlock(wallet))
		case 10:
			fmt.Println("Exiting")
			done = true
			break

		default:
			fmt.Println("Select a valid number")
		}
	}
}