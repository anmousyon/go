package setup

import (
	"github.com/julienschmidt/httprouter"
	"linklys/server/handlers"
)

func AddHandlers(r *httprouter.Router) {
	r.GET("/", handlers.IndexHandler)
	r.GET("/browse/", handlers.BrowseHandler)
	r.GET("/playlist/:id", handlers.PlaylistHandler)
	r.GET("/radio/:id", handlers.RadioHandler)
}
