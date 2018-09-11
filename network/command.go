package network

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/AnotherOctopus/goin/cnet"
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

func CmdHandler(ctx *HandlerContext, w http.ResponseWriter, r *http.Request) (int, error) {
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
	nd := ctx.Nd

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
		txtosend := cnet.LoadFTX(filename)
		for i, o := range txtosend.Outputs {
			sign := o.GenSignature(wall.Keys[addridx])
			txtosend.Outputs[i].Signature = sign
		}
		txtosend.Meta.Address = addrtouse
		txtosend.Meta.Pubkey = wall.Keys[addridx].PublicKey
		txtosend.Meta.TimePrepared = time.Now().Unix()
		txtosend.SetHash()
		err := nd.SendTx(*txtosend)
		cnet.CheckError(err)
		if err == nil {
			fmt.Println("Sent")
		}
		w.Write([]byte(fmt.Sprintf("{\"tx\":%v}", txtosend.String())))
	//[2] Manually Prepare Transaction To File"
	case 2:
		inputs := rawjson.Inputs
		outputs := rawjson.Outputs
		savetx := new(cnet.Transaction)
		savetx.Outputs = make([]cnet.Output, 0)
		savetx.Inputs = make([]cnet.Input, 0)
		for _, v := range inputs {
			var newInput cnet.Input
			prevTransHash := v.Hash
			prevTransIdx := uint32(v.Idx)
			raw, err := base64.StdEncoding.DecodeString(prevTransHash)
			cnet.CheckError(err)
			copy(newInput.PrevTransHash[:], raw)
			newInput.OutIdx = prevTransIdx
			savetx.Inputs = append(savetx.Inputs, newInput)
		}

		for _, v := range outputs {
			var newOutput cnet.Output
			sendAddr := v.Hash
			amount := uint32(v.Amt)
			raw, err := base64.StdEncoding.DecodeString(sendAddr)
			cnet.CheckError(err)
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
	return http.StatusOK, nil
}
