package articleGenerator

import (
	"net/http"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"time"
)

type NewsAPI struct {
	sources string
	base string
	sort string
	api string
}

type Sources struct {
	Status string `json:"status"`
	Sources []struct {
		ID string `json:"id"`
		Name string `json:"name"`
		Description string `json:"description"`
		URL string `json:"url"`
		Category string `json:"category"`
		Language string `json:"language"`
		Country string `json:"country"`
		UrlsToLogos struct {
			Small string `json:"small"`
			Medium string `json:"medium"`
			Large string `json:"large"`
		} `json:"urlsToLogos"`
		SortBysAvailable []string `json:"sortBysAvailable"`
	} `json:"sources"`
}

type Articles struct {
	Status string `json:"status"`
	Source string `json:"source"`
	SortBy string `json:"sortBy"`
	Articles []Article `json:"articles"`
}

type Article struct {
	Author string `json:"author"`
	Title string `json:"title"`
	Description string `json:"description"`
	URL string `json:"url"`
	URLToImage string `json:"urlToImage"`
	PublishedAt time.Time `json:"publishedAt"`
}

func (n NewsAPI) buildURLs() {
	n.api = "678e8e21125d49c9a081c49856ca0bb8"
	n.base = "https://newsapi.org/v1/articles?source="
	n.sort = "&sortBy=latest&apiKey="
	n.sources = "https://newsapi.org/v1/sources?language=en"
}

func (n NewsAPI) getSources() *Sources {
	requester := &http.Client{}
	req, _ := http.NewRequest("GET", n.sources, nil)
	resp, err := requester.Do(req)
	if err != nil {
		fmt.Println("sources request error")
	}

	s := new(Sources)
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &s)

	return s
}

func (n NewsAPI) getSourceArticles(source string) *Articles {
	requester := &http.Client{}
	url := n.base + source + n.sort + n.api
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := requester.Do(req)
	if err != nil {
		fmt.Println("articles request error")
	}

	a := new(Articles)
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &a)

	return a
}

func (n NewsAPI) GetAllArticles() []Article {
	s := n.getSources()
	all := make([]Article, 0)
	for _, source := range s.Sources {
		a := n.getSourceArticles(source.Name)
		all = append(all, a.Articles...)
	}
	return all
}