package network

import (
	"encoding/base64"
	"net/http"
)

func ClaimedTxHandler(ctx *HandlerContext, w http.ResponseWriter, r *http.Request) (int, error) {
	for _, wall := range ctx.Nd.Wallets {
		for _, txs := range wall.ClaimedTxs {
			w.Write([]byte(base64.StdEncoding.EncodeToString(txs[:])))
			w.Write([]byte("#"))
		}
		w.Write([]byte("%"))
	}
	return http.StatusOK, nil
}
