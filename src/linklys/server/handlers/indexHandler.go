package handlers

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
)

var indexTemplate = template.Must(template.ParseFiles("index.html"))

type index struct {
	Title string
	Body  string
	Links []string
}

func IndexHandler(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	links := []string{"1", "2"}
	data := &index{
		Title: "title",
		Body:  "body",
		Links: links,
	}
	indexTemplate.Execute(w, data)
}
