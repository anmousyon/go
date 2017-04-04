package main

import (
	"go-qbittorrent/qclient"
	"time"
	"fmt"
	"strconv"
)

func main() {
	qb := qclient.NewClient("http://localhost:8080/")
	//time.Sleep(time.Second)
	qb.Login("pavo", "buffalo12")
	time.Sleep(time.Second * 2)
	//qb.PauseAll()
	//time.Sleep(time.Second * 2)
	params := make(map[string]string)
	params["filter"] = "all"
	torrents := qb.Torrents(params)
	for i, t:= range torrents {
		fmt.Print("Hash " + strconv.Itoa(i) + ": ")
		fmt.Println(t.Hash)
	}
}