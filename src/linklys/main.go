package linklys

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"linklys/server/setup"
)

func main() {
	router := httprouter.New()
	setup.AddHandlers(router)
	http.ListenAndServe("localhost:8000", router)
}
