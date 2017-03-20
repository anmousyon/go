package handlers

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
)

var browseTemplate = template.Must(template.ParseFiles("browse.html"))

type browse struct {
	Title    string
	Body     string
	Articles []string
}

func BrowseHandler(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	links := []string{"1", "2"}
	data := &browse{
		Title:    "title",
		Body:     "body",
		Articles: links,
	}
	browseTemplate.Execute(w, data)
}
