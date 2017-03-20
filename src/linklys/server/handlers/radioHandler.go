package handlers

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
)

var radioTemplate = template.Must(template.ParseFiles("radio.html"))

type radio struct {
	Title    string
	User     string
	Articles []string
}

func RadioHandler(w http.ResponseWriter, _ *http.Request, id httprouter.Params) {
	fmt.Println(id)
	//query database's radios table for id

	//if id found, put into playlist struct
	data := &radio{
		Title:    "title",
		User:     "user",
		Articles: []string{"article1", "article2"},
	}

	radioTemplate.Execute(w, data)
}
