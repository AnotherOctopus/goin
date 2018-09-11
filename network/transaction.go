package network

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/AnotherOctopus/goin/cnet"
)

func TransactionHandler(ctx *HandlerContext, w http.ResponseWriter, r *http.Request) (int, error) {
	// Listen for incoming connections.
	raw, _ := ioutil.ReadAll(r.Body)
	data, _ := base64.StdEncoding.DecodeString(string(raw))
	tx := cnet.LoadTX(data)
	// Handle connections in a new goroutine.
	if !reflect.DeepEqual(cnet.GetTxFromHash(tx.Hash), tx) {
		go ctx.Nd.HandleTX(*tx)
	}
	return http.StatusOK, nil
}
