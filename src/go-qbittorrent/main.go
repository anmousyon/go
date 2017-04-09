package main

import (
	"fmt"
	"go-qbittorrent/qbit"
	"time"
)

func main() {
	qb := qbit.NewClient("http://localhost:8080/")
	//time.Sleep(time.Second)
	qb.Login("pavo", "buffalo12")
	time.Sleep(time.Second * 2)
	s, err := qb.Sync("0")
	if err != nil {
		fmt.Print("sync: ")
		fmt.Println(s)
	} else {
		fmt.Print("error: ")
		fmt.Println(err)
	}
	/*
		params := make(map[string]string)
		params["filter"] = "all"
		torrents := qb.Torrents(params)
		for i, t:= range torrents {
			fmt.Print("Hash " + strconv.Itoa(i) + ": ")
			fmt.Println(t.Hash)
		}

		time.Sleep(time.Second * 2)

		link := "magnet:?xt=urn:btih:accd778e8ef86005a9b3e8b9407675862e306a90&dn=Mr+Robot+Season+2+Complete+720p+WEB-DL+EN-SUB+x264-%5BMULVAcoded%5D+&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Fzer0day.ch%3A1337&tr=udp%3A%2F%2Fopen.demonii.com%3A1337&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Fexodus.desync.com%3A6969"
		options := make(map[string]string)
		qb.DownloadFromLink(link, options)
		fmt.Println("downloaded")

		time.Sleep(time.Second * 2)

		newTorrents := qb.Torrents(params)
		for i, t:= range newTorrents {
			fmt.Print("Hash " + strconv.Itoa(i) + ": ")
			fmt.Println(t.Hash)
		}
	*/
}
