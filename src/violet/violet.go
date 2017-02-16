package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	client := &Client{}
	client.config.token = getToken()
	fmt.Println(client.getSubreddit("python", "new"))
}

//Client struct
type Client struct {
	config struct {
		token string
	}
}

//Token struct
type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   string `json:"expires_in"`
	Scope       string `json:"scope"`
}

//Submission struct
type Submission struct {
	Author       string  `json:"author"`
	Title        string  `json:"title"`
	URL          string  `json:"url"`
	Domain       string  `json:"domain"`
	Subreddit    string  `json:"subreddit"`
	SubredditID  string  `json:"subreddit_id"`
	FullID       string  `json:"name"`
	ID           string  `json:"id"`
	Permalink    string  `json:"permalink"`
	Selftext     string  `json:"selftext"`
	ThumbnailURL string  `json:"thumbnail"`
	DateCreated  float64 `json:"created_utc"`
	NumComments  int     `json:"num_comments"`
	Score        int     `json:"score"`
	Ups          int     `json:"ups"`
	Downs        int     `json:"downs"`
	IsNSFW       bool    `json:"over_18"`
	IsSelf       bool    `json:"is_self"`
	WasClicked   bool    `json:"clicked"`
	IsSaved      bool    `json:"saved"`
	BannedBy     *string `json:"banned_by"`
}

//Response struct
type Response struct {
	Data struct {
		Children []struct {
			Data *Submission
		}
	}
}

func requestToken(url string, form url.Values) *http.Response {
	requester := &http.Client{}
	req, _ := http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	req = setAuth(req)
	req = setBasicHeader(req)
	resp, err := requester.Do(req)
	if err != nil {
		//TODO: Handle bad requests correctly
		fmt.Println("bad request")
	}
	return resp
}

func (c Client) request(url string) *http.Response {
	requester := &http.Client{}
	token := "bearer " + c.config.token
	req, _ := http.NewRequest("GET", url, nil)
	req = setFullHeader(req, token)
	resp, err := requester.Do(req)
	if err != nil {
		//TODO: Handle bad requests correctly
		fmt.Println("bad request")
	}
	return resp
}

func setBasicHeader(req *http.Request) *http.Request {
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "violet v.01 wrapper for reddit api in golang")
	return req
}

func setFullHeader(req *http.Request, token string) *http.Request {
	req.Header.Set("Authorization", token)
	req = setBasicHeader(req)
	return req
}

func setAuth(req *http.Request) *http.Request {
	req.SetBasicAuth("-400hLn5ypKJhg", "KuiYAJxYa1gqDBd4eg_Y-A3fuTw")
	return req
}

func setForm() url.Values {
	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	return form
}

func getToken() string {
	tokenURL := "https://www.reddit.com/api/v1/access_token"
	form := setForm()
	resp := requestToken(tokenURL, form)
	token := tokenStruct(resp)
	return token.AccessToken
}

func tokenString(resp *http.Response) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	return buf.String()
}

func tokenStruct(resp *http.Response) *Token {
	body := []byte(tokenString(resp))
	token := new(Token)
	json.Unmarshal(body, &token)
	return token
}

func submissionStruct(resp *http.Response) *Submission {
	body := []byte(tokenString(resp))
	response := new(Response)
	json.Unmarshal(body, &response)
	submission := new(Submission)
	submission = response.Data.Children[0].Data
	return submission
}

func (c Client) getSubreddit(subreddit string, sort string) string {
	submissionURL := "https://oauth.reddit.com/"
	url := submissionURL + "r/" + subreddit + "/" + sort
	resp := c.request(url)
	submission := submissionStruct(resp)
	return submission.Title
	//TODO: create list of posts from response
}

func (c Client) getPost(post string) string {
	submissionURL := "https://oauth.reddit.com/"
	url := submissionURL + "r/" + "TODO: get subreddit of post" + "/comments/" + post
	resp := c.request(url)
	submission := submissionStruct(resp)
	return submission.Title
	//TODO: create post from response
}

func getComments(post string) {
	//TODO: create list of comments from post
}
