package linklys

import (
	"net/http"
	"linklys/handlers"
)


func main() {
	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/browse/", handlers.BrowseHandler)
	http.ListenAndServe("localhost:8000", nil)
}
