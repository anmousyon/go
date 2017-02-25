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
	posts := client.getSubreddit("askreddit", "hot")
	for _, post := range posts {
		fmt.Println(post.Title)
	}

}

//Client struct
type Client struct {
	config struct {
		token string
	}
}

//Token is used to authorize the user's requests
type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   string `json:"expires_in"`
	Scope       string `json:"scope"`
}

// Comment is a submission to a post
type Comment struct {
	Author              string  //`json:"author"`
	Body                string  //`json:"body"`
	BodyHTML            string  //`json:"body_html"`
	Subreddit           string  //`json:"subreddit"`
	LinkID              string  //`json:"link_id"`
	ParentID            string  //`json:"parent_id"`
	SubredditID         string  //`json:"subreddit_id"`
	FullID              string  //`json:"name"`
	UpVotes             float64 //`json:"ups"`
	DownVotes           float64 //`json:"downs"`
	Created             float64 //`json:"created_utc"`
	Edited              bool    //`json:"edited"`
	BannedBy            *string //`json:"banned_by"`
	ApprovedBy          *string //`json:"approved_by"`
	AuthorFlairTxt      *string //`json:"author_flair_text"`
	AuthorFlairCSSClass *string //`json:"author_flair_css_class"`
	NumReports          *int    //`json:"num_reports"`
	Likes               *int    //`json:"likes"`
	Replies             []*Comment
}

// Post is a submission to a subreddit
type Post struct {
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

const (
	listingKind = "Listing"
	postKind    = "t3"
	commentKind = "t1"
	messageKind = "t4"
)

// author fields and body fields are set to the deletedKey if the user deletes
// their post.
const deletedKey = "[deleted]"

//Response allows for easier parsing of reddit responses
type Response struct {
	Data struct {
		Children []struct {
			Data *Post
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

func respString(resp *http.Response) []byte {
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return nil
	}
	return buf.Bytes()
}

func tokenStruct(resp *http.Response) *Token {
	body := respString(resp)
	token := new(Token)
	json.Unmarshal(body, &token)
	return token
}

func (c Client) getSubreddit(subreddit string, sort string) []*Post {
	submissionURL := "https://oauth.reddit.com/"
	url := submissionURL + "r/" + subreddit + "/" + sort
	resp := c.request(url)
	posts, err := parsePost(resp)
	if err != nil {
		fmt.Println("error on post parse")
	}
	return posts
}

// parsePost parses a post into the user facing Post struct.
func parsePost(resp *http.Response) ([]*Post, error) {
	r := new(Response)
	err := json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return nil, err
	}

	submissions := make([]*Post, len(r.Data.Children))
	for i, child := range r.Data.Children {
		submissions[i] = child.Data
	}

	return submissions, nil
}
