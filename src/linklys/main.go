package linklys

import (
	"github.com/julienschmidt/httprouter"
	"linklys/server/handlers"
	"net/http"
)

func main() {
	router := httprouter.New()
	router.GET("/", handlers.IndexHandler)
	router.GET("/browse/", handlers.BrowseHandler)
	router.GET("/playlist/:id", handlers.PlaylistHandler)
	router.GET("/radio/:id", handlers.RadioHandler)
	http.ListenAndServe("localhost:8000", router)
}
