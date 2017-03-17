package linklys

import (
	"net/http"
	"linklys/handlers"
	"github.com/julienschmidt/httprouter"
)


func main() {
	router := httprouter.New()
	router.GET("/", handlers.IndexHandler)
	router.GET("/browse/", handlers.BrowseHandler)
	router.GET("/playlist/:id", handlers.PlaylistHandler)
	router.GET("/radio/:id", handlers.RadioHandler)
	http.ListenAndServe("localhost:8000", router)
}
