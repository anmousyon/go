package linklys

import (
	"net/http"
	"linklys/server"
)

func main() {
	router := server.Setup()
	http.ListenAndServe("localhost:8000", router)
}
