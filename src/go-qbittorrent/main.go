package main

import (
	"go-qbittorrent/qclient"
	"fmt"
)

func main() {
	qb := qclient.NewClient("http://127.0.0.1:8080/")
	torrents := qb.Torrents(nil)
	fmt.Println(torrents)
}