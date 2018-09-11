package network

import (
	"encoding/base64"
	"net/http"
)

func AddressHandler(ctx *HandlerContext, w http.ResponseWriter, r *http.Request) (int, error) {
	for _, wall := range ctx.Nd.Wallets {
		for _, addr := range wall.Address {
			w.Write([]byte(base64.StdEncoding.EncodeToString(addr[:])))
			w.Write([]byte("#"))
		}
		w.Write([]byte("%"))
	}
	return http.StatusOK, nil
}
