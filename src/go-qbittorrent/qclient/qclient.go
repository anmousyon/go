package qclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/publicsuffix"
)

//Client creates a connection to qbittorrent and performs requests
type Client struct {
	http          *http.Client
	URL           string
	Authenticated bool
	Session       string //replace with session type
	Jar           http.CookieJar
}

func printResponse(resp *http.Response) {
	x := make([]byte, 256)
	x, _ = ioutil.ReadAll(resp.Body)
	fmt.Println("response: " + string(x))
}

func printRequest(req *http.Request) {
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))
}

//NewClient creates a new client connection to qbittorrent
func NewClient(url string) *Client {
	c := &Client{}

	if url[len(url)-1:] != "/" {
		url = url + "/"
	}

	c.URL = url

	c.Jar, _ = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	c.http = &http.Client{
		Jar: c.Jar,
	}
	return c
}

func (c *Client) get(endpoint string) *http.Response {
	req, _ := http.NewRequest("GET", c.URL+endpoint, nil)
	req.Header.Set("User-Agent", "autodownloader v0.1")
	//printRequest(req)

	resp, err := c.http.Do(req)
	if err != nil {
		fmt.Println(err)
		fmt.Println("error on performing get request at endpoint: " + endpoint)
	}

	//printResponse(resp)

	return resp
}

func (c *Client) getWithParams(endpoint string, params map[string]string) *http.Response {

	req, err := http.NewRequest("GET", c.URL+endpoint, nil)
	if err != nil {
		fmt.Println("error on creating get request at endpoint: " + endpoint)
	}

	req.Header.Set("User-Agent", "autodownloader v0.1")

	//add parameters to url
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	//printRequest(req)

	resp, err := c.http.Do(req)
	if err != nil {
		fmt.Println(err)
		fmt.Println("error on performing get request (with params) at endpoint: " + endpoint)
	}

	//printResponse(resp)

	req.Close = true

	return resp
}

func addForm(req *http.Request, params map[string]string) *http.Request {
	form := url.Values{}
	for k, v := range params {
		form.Add(k, v)
	}
	req.PostForm = form
	return req
}

func (c *Client) post(endpoint string, data map[string]string) *http.Response {
	req, err := http.NewRequest("POST", c.URL+endpoint, nil)
	if err != nil {
		fmt.Println("error on creating post request at endpoint: " + endpoint)
	}

	req.Header.Set("User-Agent", "autodownloader v0.1")

	req = addForm(req, data)

	resp, err := c.http.Do(req)
	if err != nil {
		fmt.Println(err)
		fmt.Println("error on performing post request at endpoint: " + endpoint)
	}

	//printResponse(resp)

	return resp

}

func (c *Client) postWithHeaders(endpoint string, data map[string]string) *http.Response {
	req, err := http.NewRequest("POST", c.URL+endpoint, nil)
	if err != nil {
		fmt.Println("error on creating post request at endpoint: " + endpoint)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "autodownloader v0.1")

	req = addForm(req, data)

	resp, err := c.http.Do(req)
	if err != nil {
		fmt.Println(err)
		fmt.Println("error on performing post (with headers) request at endpoint: " + endpoint)
	}

	//printResponse(resp)

	return resp
}

func (c *Client) postMultipart(endpoint string, data map[string]string) *http.Response {

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, val := range data {
		w.WriteField(key, val)
	}

	err := w.Close()
	if err != nil {
		fmt.Println("error on closing multipart form writer")
	}

	req, err := http.NewRequest("POST", c.URL+endpoint, &b)
	if err != nil {
		fmt.Println("error on creating post request at endpoint: " + endpoint)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	//printRequest(req)

	resp, err := c.http.Do(req)
	if err != nil {
		fmt.Println(err)
		fmt.Println("error on performing multipart post request at endpoint: " + endpoint)
	}

	printResponse(resp)

	return resp
}

//Login logs you in to the qbittorrent client
func (c *Client) Login(username string, password string) (loggedIn bool) {
	creds := make(map[string]string)
	creds["username"] = username
	creds["password"] = password

	resp := c.post("login", creds)
	cookieURL, _ := url.Parse("http://localhost:8080")

	if cookies := resp.Cookies(); len(cookies) > 0 {
		c.Jar.SetCookies(cookieURL, cookies)
	}

	c.http = &http.Client{
		Jar: c.Jar,
	}

	if resp.Status == "200 OK" {
		c.Authenticated = true
	} else {
		fmt.Println(resp.Status)
	}
	return c.Authenticated
}

//Logout logs you out of the qbittorrent client
func (c *Client) Logout() (loggedOut bool) {
	resp := c.get("logout")
	fmt.Println(resp)
	if resp.Status == "200 OK" {
		c.Authenticated = false
	} else {
		fmt.Println(resp.Status)
	}
	return !c.Authenticated
}

//Shutdown shuts down the qbittorrent client
func (c *Client) Shutdown() (shuttingDown bool) {
	resp := c.get("command/shutdown")
	return resp.Status == "200 OK"
}

//BasicTorrent holds a basic torrent object from qbittorrent
type BasicTorrent struct {
	AddedOn       int    `json:"added_on"`
	Category      string `json:"category"`
	CompletionOn  int64  `json:"completion_on"`
	Dlspeed       int    `json:"dlspeed"`
	Eta           int    `json:"eta"`
	ForceStart    bool   `json:"force_start"`
	Hash          string `json:"hash"`
	Name          string `json:"name"`
	NumComplete   int    `json:"num_complete"`
	NumIncomplete int    `json:"num_incomplete"`
	NumLeechs     int    `json:"num_leechs"`
	NumSeeds      int    `json:"num_seeds"`
	Priority      int    `json:"priority"`
	Progress      int    `json:"progress"`
	Ratio         int    `json:"ratio"`
	SavePath      string `json:"save_path"`
	SeqDl         bool   `json:"seq_dl"`
	Size          int    `json:"size"`
	State         string `json:"state"`
	SuperSeeding  bool   `json:"super_seeding"`
	Upspeed       int    `json:"upspeed"`
}

//Torrents gets a list of all torrents in qbittorrent matching your filter
func (c *Client) Torrents(filters map[string]string) (torrentList []BasicTorrent) {
	params := make(map[string]string)
	for k, v := range filters {
		if k == "status" {
			k = "filter"
		}
		params[k] = v
	}
	resp := c.getWithParams("query/torrents", params)
	var t []BasicTorrent
	json.NewDecoder(resp.Body).Decode(&t)
	return t
}

//Torrent hold torrent objects from qbittorrent
type Torrent struct {
	AdditionDate           int     `json:"addition_date"`
	Comment                string  `json:"comment"`
	CompletionDate         int     `json:"completion_date"`
	CreatedBy              string  `json:"created_by"`
	CreationDate           int     `json:"creation_date"`
	DlLimit                int     `json:"dl_limit"`
	DlSpeed                int     `json:"dl_speed"`
	DlSpeedAvg             int     `json:"dl_speed_avg"`
	Eta                    int     `json:"eta"`
	LastSeen               int     `json:"last_seen"`
	NbConnections          int     `json:"nb_connections"`
	NbConnectionsLimit     int     `json:"nb_connections_limit"`
	Peers                  int     `json:"peers"`
	PeersTotal             int     `json:"peers_total"`
	PieceSize              int     `json:"piece_size"`
	PiecesHave             int     `json:"pieces_have"`
	PiecesNum              int     `json:"pieces_num"`
	Reannounce             int     `json:"reannounce"`
	SavePath               string  `json:"save_path"`
	SeedingTime            int     `json:"seeding_time"`
	Seeds                  int     `json:"seeds"`
	SeedsTotal             int     `json:"seeds_total"`
	ShareRatio             float64 `json:"share_ratio"`
	TimeElapsed            int     `json:"time_elapsed"`
	TotalDownloaded        int     `json:"total_downloaded"`
	TotalDownloadedSession int     `json:"total_downloaded_session"`
	TotalSize              int     `json:"total_size"`
	TotalUploaded          int     `json:"total_uploaded"`
	TotalUploadedSession   int     `json:"total_uploaded_session"`
	TotalWasted            int     `json:"total_wasted"`
	UpLimit                int     `json:"up_limit"`
	UpSpeed                int     `json:"up_speed"`
	UpSpeedAvg             int     `json:"up_speed_avg"`
}

//Torrent gets a specific torrent
func (c *Client) Torrent(infoHash string) Torrent {
	resp := c.get("query/propertiesGeneral/" + strings.ToLower(infoHash))
	var t Torrent
	json.NewDecoder(resp.Body).Decode(&t)
	return t
}

//Tracker holds a tracker object from qbittorrent
type Tracker struct {
	Msg      string `json:"msg"`
	NumPeers int    `json:"num_peers"`
	Status   string `json:"status"`
	URL      string `json:"url"`
}

//TorrentTrackers gets all trackers for a specific torrent
func (c *Client) TorrentTrackers(infoHash string) []Tracker {
	resp := c.get("query/propertiesTrackers/" + strings.ToLower(infoHash))
	var t []Tracker
	json.NewDecoder(resp.Body).Decode(&t)
	return t
}

//WebSeed holds a webseed object from qbittorrent
type WebSeed struct {
	URL string `json:"url"`
}

//TorrentWebSeeds gets seeders for a specific torrent
func (c *Client) TorrentWebSeeds(infoHash string) []WebSeed {
	resp := c.get("query/propertiesWebSeeds/" + strings.ToLower(infoHash))
	var w []WebSeed
	json.NewDecoder(resp.Body).Decode(&w)
	return w
}

//TorrentFile holds a torrent file object from qbittorrent
type TorrentFile struct {
	IsSeed   bool   `json:"is_seed"`
	Name     string `json:"name"`
	Priority int    `json:"priority"`
	Progress int    `json:"progress"`
	Size     int    `json:"size"`
}

//TorrentFiles gets the files of a specifc torrent
func (c *Client) TorrentFiles(infoHash string) []TorrentFile {
	resp := c.get("query/propertiesFiles" + strings.ToLower(infoHash))
	var t []TorrentFile
	json.NewDecoder(resp.Body).Decode(&t)
	return t
}

//Status is the status of a torrent
type Status struct {
	status string
}

//Sync is the result of syncing your qbittorrent client
type Sync struct {
	rid     int
	Torrent map[string]Status
}

//Sync syncs your maindata
func (c *Client) Sync(rid string) Sync {
	params := make(map[string]string)
	params["rid"] = rid
	resp := c.getWithParams("sync/maindata", params)
	var s map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&s)
	var sync Sync
	//TODO: put map into Sync struct
	return sync
}

//DownloadFromLink starts downloading a torrent from a link
func (c *Client) DownloadFromLink(link string, options map[string]string) {
	options["urls"] = link
	c.postMultipart("command/download", options)
}

//DownloadFromFile downloads a torrent from a file
func (c *Client) DownloadFromFile(fileName string, options map[string]string) {
	//TODO: implement a way to download using file buffer and new post method?
	c.post("command/download", options)
}

//AddTrackers adds trackers to a specific torrent
func (c *Client) AddTrackers(infoHash string, trackers string) {
	params := make(map[string]string)
	params["hash"] = strings.ToLower(infoHash)
	params["urls"] = trackers
	c.post("command/addTrackers", params)
}

//process the hash list and put it into a combined (single element) map with all hashes connected with '|'
func (Client) processInfoHashList(infoHashList []string) (hashMap map[string]string) {
	d := map[string]string{}
	infoHash := ""
	for _, v := range infoHashList {
		infoHash = infoHash + "|" + v
	}
	d["hashes"] = infoHash
	return d
}

//Pause pauses a specific torrent
func (c *Client) Pause(infoHash string) {
	params := make(map[string]string)
	params["hash"] = strings.ToLower(infoHash)
	c.post("command/pause", params)
}

//PauseAll pauses all torrents
func (c *Client) PauseAll() {
	c.get("command/pauseAll")
}

//PauseMultiple pauses a list of torrents
func (c *Client) PauseMultiple(infoHashList []string) {
	params := c.processInfoHashList(infoHashList)
	c.post("command/pauseAll", params)
}

//SetLabel sets the labels for a list of torrents
func (c *Client) SetLabel(infoHashList []string, label string) {
	params := c.processInfoHashList(infoHashList)
	params["label"] = label
	c.post("command/setLabel", params)
}

//SetCategory sets the category for a list of torrents
func (c *Client) SetCategory(infoHashList []string, category string) {
	params := c.processInfoHashList(infoHashList)
	params["category"] = category
	c.post("command/setLabel", params)
}

//Resume resumes a specific torrent
func (c *Client) Resume(infoHash string) {
	params := make(map[string]string)
	params["hash"] = strings.ToLower(infoHash)
	c.post("command/resume", params)
}

//ResumeAll resumes all torrents
func (c *Client) ResumeAll(infoHashList []string) {
	c.get("command/resumeAll")
}

//ResumeMultiple resumes a list of torrents
func (c *Client) ResumeMultiple(infoHashList []string) {
	params := c.processInfoHashList(infoHashList)
	c.post("command/resumeAll", params)
}

//DeleteTemp deletes the temporary files for a list of torrents
func (c *Client) DeleteTemp(infoHashList []string) {
	params := c.processInfoHashList(infoHashList)
	c.post("command/delete", params)
}

//DeletePermanently deletes all files for a list of torrents
func (c *Client) DeletePermanently(infoHashList []string) {
	params := c.processInfoHashList(infoHashList)
	c.post("command/deletePerm", params)
}

//Recheck rechecks a list of torrents
func (c *Client) Recheck(infoHashList []string) {
	params := c.processInfoHashList(infoHashList)
	c.post("command/recheck", params)
}

//IncreasePriority increases the priority of a list of torrents
func (c *Client) IncreasePriority(infoHashList []string) {
	params := c.processInfoHashList(infoHashList)
	c.post("command/increasePrio", params)
}

//DecreasePriority decreases the priority of a list of torrents
func (c *Client) DecreasePriority(infoHashList []string) {
	params := c.processInfoHashList(infoHashList)
	c.post("command/decreasePrio", params)
}

//SetMaxPriority sets the max priority for a list of torrents
func (c *Client) SetMaxPriority(infoHashList []string) {
	params := c.processInfoHashList(infoHashList)
	c.post("command/topPrio", params)
}

//SetMinPriority sets the min priority for a list of torrents
func (c *Client) SetMinPriority(infoHashList []string) {
	params := c.processInfoHashList(infoHashList)
	c.post("command/bottomPrio", params)
}

//SetFilePriority sets the priority for a specific torrent file
func (c *Client) SetFilePriority(infoHash string, fileID string, priority string) {
	//TODO: find a way to work with files
	priorities := [...]string{"0", "1", "2", "7"}
	for _, v := range priorities {
		if v == priority {
			fmt.Println("error, priority no tavailable")
		}
	}
	params := make(map[string]string)
	params["hash"] = infoHash
	params["id"] = fileID
	params["priority"] = priority
	c.post("command/setFilePriority", params)
}

//GetGlobalDownloadLimit gets the global download limit of your qbittorrent client
func (c *Client) GetGlobalDownloadLimit() (limit int) {
	resp := c.get("command/getGlobalDlLimit")
	var l int
	json.NewDecoder(resp.Body).Decode(&l)
	return l
}

//SetGlobalDownloadLimit sets the global download limit of your qbittorrent client
func (c *Client) SetGlobalDownloadLimit(limit string) {
	params := make(map[string]string)
	params["limit"] = limit
	c.post("command/setGlobalDlLimit", params)
}

//GetGlobalUploadLimit gets the global upload limit of your qbittorrent client
func (c *Client) GetGlobalUploadLimit() (limit int) {
	resp := c.get("command/getGlobalUpLimit")
	var l int
	json.NewDecoder(resp.Body).Decode(&l)
	return l
}

//SetGlobalUploadLimit sets the global upload limit of your qbittorrent client
func (c *Client) SetGlobalUploadLimit(limit string) {
	params := make(map[string]string)
	params["limit"] = limit
	c.post("command/setGlobalUpLimit", params)
}

//GetTorrentDownloadLimit gets the download limit for a list of torrents
func (c *Client) GetTorrentDownloadLimit(infoHashList []string) (limits map[string]string) {
	params := c.processInfoHashList(infoHashList)
	resp := c.post("command/getTorrentsDlLimit", params)
	var l map[string]string
	json.NewDecoder(resp.Body).Decode(&l)
	return l
}

//SetTorrentDownloadLimit sets the download limit for a list of torrents
func (c *Client) SetTorrentDownloadLimit(infoHashList []string, limit string) {
	params := c.processInfoHashList(infoHashList)
	params["limit"] = limit
	c.post("command/setTorrentsDlLimit", params)
}

//GetTorrentUploadLimit gets the upload limit for a list of torrents
func (c *Client) GetTorrentUploadLimit(infoHashList []string) (limits map[string]string) {
	params := c.processInfoHashList(infoHashList)
	resp := c.post("command/getTorrentsUpLimit", params)
	var l map[string]string
	json.NewDecoder(resp.Body).Decode(&l)
	return l
}

//SetTorrentUploadLimit sets the upload limit of a list of torrents
func (c *Client) SetTorrentUploadLimit(infoHashList []string, limit string) {
	params := c.processInfoHashList(infoHashList)
	params["limit"] = limit
	c.post("command/setTorrentsUpLimit", params)
}

//SetPreferences sets the preferences of your qbittorrent client
func (c *Client) SetPreferences(params map[string]string) {
	c.postWithHeaders("command/setPreferences", params)
}

//GetAlternativeSpeedStatus toggles the alternative speed status of your qbittorrent client
func (c *Client) GetAlternativeSpeedStatus() (status bool) {
	resp := c.get("command/alternativeSpeedLimitsEnabled")
	var s bool
	json.NewDecoder(resp.Body).Decode(&s)
	return s
}

//ToggleAlternativeSpeed toggles the alternative speed of your qbittorrent client
func (c *Client) ToggleAlternativeSpeed() {
	c.get("command/toggleAlternativeSpeedLimits")
}

//ToggleSequentialDownload toggles the download sequence of a list of torrents
func (c *Client) ToggleSequentialDownload(infoHashList []string) {
	params := c.processInfoHashList(infoHashList)
	c.post("command/toggleSequentialDownload", params)
}

//ToggleFirstLastPiecePriority toggles first last piece priority of a list of torrents
func (c *Client) ToggleFirstLastPiecePriority(infoHashList []string) {
	params := c.processInfoHashList(infoHashList)
	c.post("command/toggleFirstLastPiecePrio", params)
}

//ForceStart force stats a list of torrents
func (c *Client) ForceStart(infoHashList []string, value bool) {
	params := c.processInfoHashList(infoHashList)
	params["value"] = strconv.FormatBool(value)
	c.post("command/setForceStart", params)
}
