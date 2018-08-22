package cnet

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/AnotherOctopus/goin/wallet"
)

type Cmd struct {
	Job      int
	Filename string
	Walidx   int
	Addridx  int
	Inputs   []struct {
		Hash string
		Idx  int
	}
	Outputs []struct {
		Hash string
		Amt  int
	}
}

func (nd *Node) HandleCmd(w http.ResponseWriter, r *http.Request) {
	/*[1] Send Transaction From File"
	  [2] Manually Prepare Transaction To File"
	  [3] View Current Balence To File"
	  [4] Make A New Wallet To File"
	  [5] Load A Wallet From File"
	  [6] Save A Wallet To File
	  Other exits

		expects json of
		{
			<key>: <valuetype>, <which tasks require>
			"job": int, [1,2,3,4,5,6]
			"filename":string, [1,2,3,4,5,6]
			"walindex":int,[1,3,6]
			"addridx": int, [1]
			"inputs": [(Hash,int)], [2]
			"outputs": [(Hash,int)] [2]
		}
	*/
	var rawjson Cmd

	raw := make([]byte, r.ContentLength)
	r.Body.Read(raw)
	json.Unmarshal(raw, &rawjson)
	jobnum := rawjson.Job
	filename := rawjson.Filename

	switch jobnum {
	//[1] Send Transaction From File"
	case 1:
		walidx := rawjson.Walidx
		addridx := rawjson.Addridx
		if walidx > len(nd.Wallets) {
			w.Write([]byte("X"))
			break
		}
		wall := nd.Wallets[rawjson.Walidx]
		if addridx > len(wall.Keys) {
			w.Write([]byte("X"))
			break
		}
		addrtouse := wall.Address[addridx]
		txtosend := LoadFTX(filename)
		for i, o := range txtosend.Outputs {
			sign := o.GenSignature(wall.Keys[addridx])
			txtosend.Outputs[i].Signature = sign
		}
		txtosend.Meta.Address = addrtouse
		txtosend.Meta.Pubkey = wall.Keys[addridx].PublicKey
		txtosend.Meta.TimePrepared = time.Now().Unix()
		txtosend.SetHash()
		err := nd.SendTx(*txtosend)
		CheckError(err)
		if err == nil {
			fmt.Println("Sent")
		}
		w.Write([]byte(fmt.Sprintf("{\"tx\":%v}", txtosend.String())))
	//[2] Manually Prepare Transaction To File"
	case 2:
		inputs := rawjson.Inputs
		outputs := rawjson.Outputs
		savetx := new(Transaction)
		savetx.Outputs = make([]Output, 0)
		savetx.Inputs = make([]Input, 0)
		for _, v := range inputs {
			var newInput Input
			prevTransHash := v.Hash
			prevTransIdx := uint32(v.Idx)
			raw, err := base64.StdEncoding.DecodeString(prevTransHash)
			checkerror(err)
			copy(newInput.PrevTransHash[:], raw)
			newInput.OutIdx = prevTransIdx
			savetx.Inputs = append(savetx.Inputs, newInput)
		}

		for _, v := range outputs {
			var newOutput Output
			sendAddr := v.Hash
			amount := uint32(v.Amt)
			raw, err := base64.StdEncoding.DecodeString(sendAddr)
			checkerror(err)
			copy(newOutput.Addr[:], raw)
			newOutput.Amount = amount
			newOutput.Signature = make([]byte, 0)
			savetx.Outputs = append(savetx.Outputs, newOutput)
		}
		_, provtx := savetx.Dump()
		ioutil.WriteFile(filename, provtx, 0644)
		w.Write([]byte(fmt.Sprintf("{\"tx\":%v}", savetx.String())))
	//[3] View Current Balence To File"
	case 3:
		idx := rawjson.Walidx
		if idx > len(nd.Wallets) {
			w.Write([]byte("X"))
			break
		}
		wall := nd.Wallets[idx]
		w.Write([]byte(fmt.Sprintf("{\"value\":%v}", wall.GetTotal())))
	//[4] Make A New Wallet To File"
	case 4:
		wall := wallet.NewWallet(1)
		ioutil.WriteFile(filename, wall.Dump(), 0644)
		w.Write([]byte("{\"wallet\":\"written\"}"))
	//[5] Load A Wallet From File"
	case 5:
		rawdata, err := ioutil.ReadFile(filename)
		wallet.CheckError(err)
		nd.Wallets = append(nd.Wallets, wallet.LoadWallet(rawdata))
		w.Write([]byte("{\"wallet\":\"loaded\"}"))
	//[6] Save A Wallet To File
	case 6:
		idx := rawjson.Walidx
		if idx > len(nd.Wallets) {
			w.Write([]byte("X"))
			break
		}
		wall := nd.Wallets[idx]
		ioutil.WriteFile(filename, wall.Dump(), 0644)
		w.Write([]byte("{\"wallet\":\"saved\"}"))
	default:
		w.Write([]byte("{\"done\":\"done\"}"))
	}
}

func (nd *Node) ServeAddresses(w http.ResponseWriter, r *http.Request) {
	for _, wall := range nd.Wallets {
		for _, addr := range wall.Address {
			w.Write([]byte(base64.StdEncoding.EncodeToString(addr[:])))
			w.Write([]byte("#"))
		}
		w.Write([]byte("%"))
	}
}
func (nd *Node) ServeClaimedTx(w http.ResponseWriter, r *http.Request) {
	for _, wall := range nd.Wallets {
		for _, txs := range wall.ClaimedTxs {
			w.Write([]byte(base64.StdEncoding.EncodeToString(txs[:])))
			w.Write([]byte("#"))
		}
		w.Write([]byte("%"))
	}
}
func (nd *Node) CmdListener() {
	http.HandleFunc("/cmd", nd.HandleCmd)
}

func (nd *Node) ExposeInfo() {
	http.HandleFunc("/addresses", nd.ServeAddresses)
	http.HandleFunc("/claimedtxs", nd.ServeClaimedTx)
}
