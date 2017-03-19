package articleGenerator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"linklys/models"
	"net/http"
)

type NewsAPI struct {
	sources string
	base    string
	sort    string
	api     string
}

type Sources struct {
	Status  string          `json:"status"`
	Sources []models.Source `json:"sources"`
}

type Articles struct {
	Status   string           `json:"status"`
	Source   string           `json:"source"`
	SortBy   string           `json:"sortBy"`
	Articles []models.Article `json:"articles"`
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

func (n NewsAPI) GetAllArticles() []models.Article {
	s := n.getSources()
	all := make([]models.Article, 0)
	for _, source := range s.Sources {
		a := n.getSourceArticles(source.Name)
		all = append(all, a.Articles...)
	}
	return all
}
