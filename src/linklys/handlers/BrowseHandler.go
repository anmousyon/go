package handlers

import (
	"net/http"
	"html/template"
	"github.com/julienschmidt/httprouter"
)

var browseTemplate = template.Must(template.ParseFiles("browse.html"))

type Browse struct {
	Title string
	Body string
	Links []string
}

func BrowseHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	links := []string{"1", "2"}
	data := &Browse {
		Title: "title",
		Body: "body",
		Links: links,
	}
	browseTemplate.Execute(w, data)
}