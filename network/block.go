package network

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"

	"github.com/AnotherOctopus/goin/cnet"
)

func BlockHandler(ctx *HandlerContext, w http.ResponseWriter, r *http.Request) (int, error) {
	raw, err := ioutil.ReadAll(r.Body)
	cnet.CheckError(err)
	data, err := base64.StdEncoding.DecodeString(string(raw))
	blk := cnet.LoadBlk(data)
	// Handle connections in a new goroutine.
	go ctx.Nd.HandleBlk(blk)
	return http.StatusOK, nil
}
