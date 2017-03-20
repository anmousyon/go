package handlers

import (
	"net/http"
	"html/template"
	"github.com/julienschmidt/httprouter"
)

var indexTemplate = template.Must(template.ParseFiles("index.html"))

type Index struct {
	Title string
	Body string
	Links []string
}

func IndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	links := []string{"1", "2"}
	data := &Index {
		Title: "title",
		Body: "body",
		Links: links,
	}
	indexTemplate.Execute(w, data)
}