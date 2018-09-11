package network

import (
	"io/ioutil"
	"net/http"
)

func JoinHandler(ctx *HandlerContext, w http.ResponseWriter, r *http.Request) (int, error) {
	peer, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return http.StatusOK, nil
	}
	ctx.Nd.AddPeer(string(peer))
	w.Write([]byte("Joined"))
	return http.StatusOK, nil
}
