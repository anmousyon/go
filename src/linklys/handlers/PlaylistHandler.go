package handlers

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"linklys/models"
)

var playlistTemplate = template.Must(template.ParseFiles("playlist.html"))

type Playlists struct {
	Title    string
	User     string
	Articles []string
}

func PlaylistHandler(w http.ResponseWriter, r *http.Request, id httprouter.Params) {
	fmt.Println(id)
	//query database for id

	//if id found, put into playlist struct
	data := &Playlists{
		Title:    "title",
		User:     "body",
		Articles: []string{"1", "2"},
	}

	playlistTemplate.Execute(w, data)
}
