package models

import (
	"time"
)

type User struct {
	ID    string
	Name  string
	Email string
}

type Article struct {
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	URLToImage  string    `json:"urlToImage"`
	Created     time.Time `json:"publishedAt"`
}

type Source struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	URL              string   `json:"url"`
	Category         string   `json:"category"`
	Language         string   `json:"language"`
	Country          string   `json:"country"`
	SortBysAvailable []string `json:"sortBysAvailable"`
}

type Playlist struct {
	ID      string
	Name    string
	User_id string
	Created string
}

type Playlist_Articles struct {
	Playlist_id string
	Article_id  string
}

type Radio struct {
	ID      string
	Name    string
	User_id string
	Created string
}
