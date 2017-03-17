package handlers

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
)

var radioTemplate = template.Must(template.ParseFiles("radio.html"))

type Radio struct {
	Title    string
	User     string
	Articles []string
}

func RadioHandler(w http.ResponseWriter, r *http.Request, id httprouter.Params) {
	fmt.Println(id)
	//query database for id

	//if id found, put into playlist struct
	data := &Playlist{
		Title:    "title",
		User:     "user",
		Articles: []string{"article1", "article2"},
	}

	radioTemplate.Execute(w, data)
}
