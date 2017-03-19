package handlers

import (
	"net/http"
	"html/template"
)

var browseTemplate = template.Must(template.ParseFiles("browse.html"))

type Browse struct {
	Title string
	Body string
	Links []string
}

func BrowseHandler(w http.ResponseWriter, r *http.Request) {
	links := []string{"1", "2"}
	data := &Browse {
		Title: "title",
		Body: "body",
		Links: links,
	}
	browseTemplate.Execute(w, data)
}