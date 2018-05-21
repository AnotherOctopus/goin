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
	"encoding/json"
	"encoding/base64"
	"constants"
)
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}
func test(){
	peerips := []string{"127.0.0.1"}
	nd := cnet.New(peerips)
	go nd.TxListener()
	wait :=	time.NewTimer(time.Millisecond*100)
	<-wait.C
	fmt.Println("Listening for Transactions...")
	fmt.Println("Launching Goin CLI! ")

	filename := "firstwallet"
	filename = strings.TrimSpace(filename)
	rawdata,err := ioutil.ReadFile(filename)
	wallet.CheckError(err)
	nd.Wallets = append(nd.Wallets, wallet.LoadWallet(rawdata))

	addidx := 0
	windex := 0
	filename = "tx1.json"
	w := nd.Wallets[windex]
	addrtouse := w.Address[addidx]

	txtosend := cnet.LoadFTX(filename)
	txtosend.Meta.Address = addrtouse
	txtosend.Meta.Pubkey = w.Keys[addidx].PublicKey
	txtosend.Meta.TimePrepared = time.Now().Unix()
	txtosend.Outputs[0].Signature = txtosend.Outputs[0].GenSignature(w.Keys[addidx])
	txtosend.SetHash()

	err = nd.SendTx(*txtosend)
	CheckError(err)
	fmt.Println("Sent")

	wait =	time.NewTimer(time.Second*2)
	<-wait.C
	panic("Done")
}
func main(){
	//test()
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

			err = nd.SendTx(*txtosend)
			CheckError(err)
			fmt.Println("Sent")
		case 2:
			savetx := new(cnet.AnonTransaction)
			savetx.Outputs = make([]cnet.Output,0)
			savetx.Inputs = make([]cnet.Input,0)
			fmt.Println("Select Filename of Transaction to prepare")
			filename, _ := reader.ReadString('\n')
			filename = strings.TrimSpace(filename)
			var inp byte
			for inp != 'd'{
				fmt.Println("Select An Option")
				fmt.Println("[i] Add Input")
				fmt.Println("[o] Add Output")
				fmt.Println("[d] Finish and Save")
				rawinp, _ := reader.ReadString('\n')
				inp = rawinp[0]
				switch inp {
				case 'i':
					var newInput cnet.Input
					var hashbuffer [constants.HASHSIZE]byte
					fmt.Println("Hash for Previous Transaction for input?")
					prev, _ := reader.ReadString('\n')
					prev = strings.TrimSpace(prev)
					CheckError(err)
					newHash, err := base64.StdEncoding.DecodeString(prev)
					if err != nil {
						fmt.Println("Not Valid b64")
						continue
					}
					fmt.Println("Index for Previous Transaction for input?")
					prev, _ = reader.ReadString('\n')
					prev = strings.TrimSpace(prev)
					CheckError(err)
					transidx,err := strconv.ParseInt(prev,10,8)
					if err != nil {
						fmt.Println("INPUT IS NOT A NUMBER")
						fmt.Println(err)
						continue
					}
					newInput.OutIdx = uint32(transidx)
					copy(hashbuffer[:],newHash)
					newInput.PrevTransHash = hashbuffer
					savetx.Inputs = append(savetx.Inputs,newInput)

				case 'o':
					var newOutput cnet.Output
					var hashbuffer [constants.ADDRESSSIZE]byte
					fmt.Println("Hash for Address for output?")
					out, _ := reader.ReadString('\n')
					out = strings.TrimSpace(out)
					CheckError(err)
					newHash, err := base64.StdEncoding.DecodeString(out)
					if err != nil {
						fmt.Println("Not Valid b64")
						continue
					}
					fmt.Println("Amount to send?")
					out, _ = reader.ReadString('\n')
					out = strings.TrimSpace(out)
					CheckError(err)
					amount,err := strconv.ParseInt(out,10,8)
					if err != nil {
						fmt.Println("INPUT IS NOT A NUMBER")
						fmt.Println(err)
						continue
					}
					newOutput.Amount = uint32(amount)
					copy(hashbuffer[:],newHash)
					newOutput.Signature = make([]byte,0)
					newOutput.Addr = hashbuffer
					savetx.Outputs = append(savetx.Outputs,newOutput)
				case 'd':
					data,err := json.Marshal(savetx)
					CheckError(err)
					ioutil.WriteFile(filename,data,0644)
				default:
					fmt.Println(inp, " is not a valid input")
				}
			}


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