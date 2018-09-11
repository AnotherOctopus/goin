package network

import "net/http"

func RootHandler(ctx *HandlerContext, w http.ResponseWriter, r *http.Request) (int, error) {
	w.Write([]byte("foo"))
	return http.StatusOK, nil
}
