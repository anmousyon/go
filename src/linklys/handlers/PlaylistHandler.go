package handlers

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
)

var playlistTemplate = template.Must(template.ParseFiles("browse.html"))

type Playlist struct {
	Title    string
	User     string
	Articles []string
}

func PlaylistHandler(w http.ResponseWriter, r *http.Request, id httprouter.Params) {
	fmt.Println(id)
	//query database for id

	//if id found, put into playlist struct
	data := &Playlist{
		Title:    "title",
		User:     "body",
		Articles: []string{"1", "2"},
	}
	
	playlistTemplate.Execute(w, data)
}
