package server

import (
	"log"
	"net/http"

	"blockchain_go/jsonrpc/service"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
)

func Run() {
	port := ":1122"

	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(new(service.Server), "")

	r := mux.NewRouter()
	r.Handle("/rpc", s)

	log.Printf("Ouvindo o servidor na porta %s", port)
	log.Fatal(http.ListenAndServe(port, r))

}
