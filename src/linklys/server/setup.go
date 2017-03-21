package server

import (
	"github.com/julienschmidt/httprouter"
	"linklys/server/handlers"
)

func Setup() *httprouter.Router {
	r := httprouter.New()
	r.GET("/", handlers.IndexHandler)
	r.GET("/browse/", handlers.BrowseHandler)
	r.GET("/playlist/:id", handlers.PlaylistHandler)
	r.GET("/radio/:id", handlers.RadioHandler)
	return r
}
