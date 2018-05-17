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
	"encoding/hex"
	"time"
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
	nd := cnet.New(peerips)
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
			if len(nd.Wallets) == 0{
				fmt.Println("Load a Wallet First!")
				break
			}
			fmt.Println("Select a Wallet index to use")
			for i := range nd.Wallets {
				fmt.Println(fmt.Sprintf("[%v]",i))
			}
			text, _ := reader.ReadString('\n')
			text = strings.TrimSpace(text)
			windex,err := strconv.ParseInt(text,10,8)
			if err != nil {
				fmt.Println("INPUT IS NOT A NUMBER")
				fmt.Println(err)
				continue
			}
			w := nd.Wallets[windex]
			fmt.Println("Select an Address to use")
			for i,addr := range w.Address {
				fmt.Println(fmt.Sprintf("[%v]: %v",i,hex.EncodeToString(addr[:])))
			}
			text, _ = reader.ReadString('\n')
			text = strings.TrimSpace(text)
			addidx,err := strconv.ParseInt(text,10,8)
			if err != nil {
				fmt.Println("INPUT IS NOT A NUMBER")
				fmt.Println(err)
				continue
			}
			addrtouse := w.Address[addidx]
			fmt.Println("Using ",hex.EncodeToString(addrtouse[:]))
			fmt.Println("Select Filename of Transaction to send")
			filename, _ := reader.ReadString('\n')
			filename = strings.TrimSpace(filename)

			txtosend := cnet.LoadFTX(filename)
			txtosend.Meta.Address = addrtouse
			txtosend.Meta.Pubkey = w.Keys[addidx].PublicKey
			txtosend.Meta.TimePrepared = time.Now().Unix()
			txtosend.SetHash()

			log.Println(txtosend)
			err = nd.SendTx(txtosend)
			CheckError(err)
			fmt.Println("Sent")
		case 2:
			fmt.Println("Prep")

		case 3:
			fmt.Println("Balence")

		case 4:
			w := wallet.NewWallet(3)
			fmt.Println("Select Filename of where to save wallet")
			filename, _ := reader.ReadString('\n')
			filename = strings.TrimSpace(filename)
			ioutil.WriteFile(filename,w.Dump(),0644)
			nd.Wallets = append(nd.Wallets,w)

		case 5:
			fmt.Println("Select Filename of wallet")
			filename, _ := reader.ReadString('\n')
			filename = strings.TrimSpace(filename)
			rawdata,err := ioutil.ReadFile(filename)
			wallet.CheckError(err)
			nd.Wallets = append(nd.Wallets, wallet.LoadWallet(rawdata))

		case 10:
			fmt.Println("Exiting")
			done = true
			break

		default:
			fmt.Println("Select a valid number")
		}
	}
}