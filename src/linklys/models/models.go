package models

import "time"

type User struct {
	ID string
	Name string
	Email string
}

type Article struct {
	Author string
	Title string
	Description string
	URL string
	URLToImage string
	Created time.Time
}

type Source struct {
	ID string
	Name string
	Description string
	URL string
	Category string
	Language string
	Country string
}

type Playlist struct {
	ID string
	Name string
	User_id string
	Created string
}

type Playlist_Articles struct {
	Playlist_id string
	Article_id string
}