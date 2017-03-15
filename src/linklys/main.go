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
	http.ListenAndServe("localhost:8000", router)
}
