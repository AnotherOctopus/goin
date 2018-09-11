package network

import (
	"log"
	"net/http"

	"github.com/AnotherOctopus/goin/cnet"
	"github.com/gorilla/mux"
)

type HandlerContext struct {
	Nd *cnet.Node
}

type routeHandler struct {
	context *HandlerContext
	H       func(*HandlerContext, http.ResponseWriter, *http.Request) (int, error)
}

func (rh routeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := rh.H(rh.context, w, r)
	if err != nil {
		log.Println("HTTP %d", status)
	}
}
func Handlers(context *HandlerContext) *mux.Router {
	r := mux.NewRouter()
	r.Handle("/", routeHandler{context, RootHandler}).Methods("GET")
	r.Handle("/transaction", routeHandler{context, TransactionHandler}).Methods("PUT")
	r.Handle("/block", routeHandler{context, BlockHandler}).Methods("PUT")
	r.Handle("/command", routeHandler{context, CmdHandler}).Methods("POST")
	r.Handle("/addresses", routeHandler{context, AddressHandler}).Methods("GET")
	r.Handle("/claimedtxs", routeHandler{context, ClaimedTxHandler}).Methods("GET")
	r.Handle("/join", routeHandler{context, JoinHandler}).Methods("POST")

	return r
}
