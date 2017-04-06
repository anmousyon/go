package qclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	//"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	"os"

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

/*
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
*/

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

func (c *Client) get(endpoint string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", c.URL+endpoint, nil)
	req.Header.Set("User-Agent", "autodownloader v0.1")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) getWithParams(endpoint string, params map[string]string) (*http.Response, error) {

	req, err := http.NewRequest("GET", c.URL+endpoint, nil)

	req.Header.Set("User-Agent", "autodownloader v0.1")

	//add parameters to url
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	req.Close = true

	return resp, nil
}

func addForm(req *http.Request, params map[string]string) *http.Request {
	form := url.Values{}
	for k, v := range params {
		form.Add(k, v)
	}
	req.PostForm = form
	return req
}

func (c *Client) post(endpoint string, data map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("POST", c.URL+endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "autodownloader v0.1")

	req = addForm(req, data)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil

}

func (c *Client) postWithHeaders(endpoint string, data map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("POST", c.URL+endpoint, nil)
	if err != nil {
		return nil, err
	}

	//add headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "autodownloader v0.1")

	req = addForm(req, data)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) postMultipart(endpoint string, data map[string]string) (*http.Response, error) {

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, val := range data {
		w.WriteField(key, val)
	}
	contentType := w.FormDataContentType()

	err := w.Close()
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(c.URL+endpoint, contentType, &b)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) postMultipartFile(endpoint string, data map[string]string, file string) (*http.Response, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	fileContents, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	f.Close()

	part, err := w.CreateFormFile("file", fi.Name())
	if err != nil {
		return nil, err
	}

	part.Write(fileContents)

	for key, val := range data {
		w.WriteField(key, val)
	}
	contentType := w.FormDataContentType()

	err = w.Close()
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(c.URL+endpoint, contentType, &b)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

//Login logs you in to the qbittorrent client
func (c *Client) Login(username string, password string) (loggedIn bool, err error) {
	creds := make(map[string]string)
	creds["username"] = username
	creds["password"] = password

	resp, err := c.post("login", creds)
	if err != nil {
		return false, err
	}
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
	return c.Authenticated, nil
}

//Logout logs you out of the qbittorrent client
func (c *Client) Logout() (loggedOut bool, err error) {
	resp, err := c.get("logout")
	if err != nil {
		return false, err
	}
	fmt.Println(resp)
	if resp.Status == "200 OK" {
		c.Authenticated = false
	} else {
		fmt.Println(resp.Status)
	}
	return !c.Authenticated, nil
}

//Shutdown shuts down the qbittorrent client
func (c *Client) Shutdown() (shuttingDown bool, err error) {
	resp, err := c.get("command/shutdown")
	return resp.Status == "200 OK", err
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
func (c *Client) Torrents(filters map[string]string) (torrentList []BasicTorrent, err error) {
	var t []BasicTorrent
	params := make(map[string]string)
	for k, v := range filters {
		if k == "status" {
			k = "filter"
		}
		params[k] = v
	}
	resp, err := c.getWithParams("query/torrents", params)
	if err != nil {
		return t, err
	}
	json.NewDecoder(resp.Body).Decode(&t)
	return t, nil
}

//Torrent holds a torrent object from qbittorrent
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
func (c *Client) Torrent(infoHash string) (Torrent, error) {
	var t Torrent
	resp, err := c.get("query/propertiesGeneral/" + strings.ToLower(infoHash))
	if err != nil {
		return t, err
	}
	json.NewDecoder(resp.Body).Decode(&t)
	return t, nil
}

//Tracker holds a tracker object from qbittorrent
type Tracker struct {
	Msg      string `json:"msg"`
	NumPeers int    `json:"num_peers"`
	Status   string `json:"status"`
	URL      string `json:"url"`
}

//TorrentTrackers gets all trackers for a specific torrent
func (c *Client) TorrentTrackers(infoHash string) ([]Tracker, error) {
	var t []Tracker
	resp, err := c.get("query/propertiesTrackers/" + strings.ToLower(infoHash))
	if err != nil {
		return t, err
	}
	json.NewDecoder(resp.Body).Decode(&t)
	return t, nil
}

//WebSeed holds a webseed object from qbittorrent
type WebSeed struct {
	URL string `json:"url"`
}

//TorrentWebSeeds gets seeders for a specific torrent
func (c *Client) TorrentWebSeeds(infoHash string) ([]WebSeed, error) {
	var w []WebSeed
	resp, err := c.get("query/propertiesWebSeeds/" + strings.ToLower(infoHash))
	if err != nil {
		return w, err
	}
	json.NewDecoder(resp.Body).Decode(&w)
	return w, nil
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
func (c *Client) TorrentFiles(infoHash string) ([]TorrentFile, error) {
	var t []TorrentFile
	resp, err := c.get("query/propertiesFiles" + strings.ToLower(infoHash))
	if err != nil {
		return t, err
	}
	json.NewDecoder(resp.Body).Decode(&t)
	return t, nil
}

//Sync holds the sync response struct
type Sync struct {
	Categories  []string `json:"categories"`
	FullUpdate  bool     `json:"full_update"`
	Rid         int      `json:"rid"`
	ServerState struct {
		ConnectionStatus  string `json:"connection_status"`
		DhtNodes          int    `json:"dht_nodes"`
		DlInfoData        int    `json:"dl_info_data"`
		DlInfoSpeed       int    `json:"dl_info_speed"`
		DlRateLimit       int    `json:"dl_rate_limit"`
		Queueing          bool   `json:"queueing"`
		RefreshInterval   int    `json:"refresh_interval"`
		UpInfoData        int    `json:"up_info_data"`
		UpInfoSpeed       int    `json:"up_info_speed"`
		UpRateLimit       int    `json:"up_rate_limit"`
		UseAltSpeedLimits bool   `json:"use_alt_speed_limits"`
	} `json:"server_state"`
	Torrents map[string]Torrent `json:"torrents"`
}

//Sync syncs your maindata
func (c *Client) Sync(rid string) (Sync, error) {
	var s Sync
	params := make(map[string]string)
	params["rid"] = rid
	resp, err := c.getWithParams("sync/maindata", params)
	if err != nil {
		return s, err
	}
	//printResponse(resp)
	json.NewDecoder(resp.Body).Decode(&s)
	return s, nil
}

//DownloadFromLink starts downloading a torrent from a link
func (c *Client) DownloadFromLink(link string, options map[string]string) (*http.Response, error) {
	options["urls"] = link
	return c.postMultipart("command/download", options)
}

//DownloadFromFile downloads a torrent from a file
func (c *Client) DownloadFromFile(fileName string, options map[string]string) (*http.Response, error) {
	return c.postMultipartFile("command/download", options, fileName)
}

//AddTrackers adds trackers to a specific torrent
func (c *Client) AddTrackers(infoHash string, trackers string) (*http.Response, error) {
	params := make(map[string]string)
	params["hash"] = strings.ToLower(infoHash)
	params["urls"] = trackers
	return c.post("command/addTrackers", params)
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
func (c *Client) Pause(infoHash string) (*http.Response, error) {
	params := make(map[string]string)
	params["hash"] = strings.ToLower(infoHash)
	return c.post("command/pause", params)
}

//PauseAll pauses all torrents
func (c *Client) PauseAll() (*http.Response, error) {
	return c.get("command/pauseAll")
}

//PauseMultiple pauses a list of torrents
func (c *Client) PauseMultiple(infoHashList []string) (*http.Response, error) {
	params := c.processInfoHashList(infoHashList)
	return c.post("command/pauseAll", params)
}

//SetLabel sets the labels for a list of torrents
func (c *Client) SetLabel(infoHashList []string, label string) (*http.Response, error) {
	params := c.processInfoHashList(infoHashList)
	params["label"] = label
	return c.post("command/setLabel", params)
}

//SetCategory sets the category for a list of torrents
func (c *Client) SetCategory(infoHashList []string, category string) (*http.Response, error) {
	params := c.processInfoHashList(infoHashList)
	params["category"] = category
	return c.post("command/setLabel", params)
}

//Resume resumes a specific torrent
func (c *Client) Resume(infoHash string) (*http.Response, error) {
	params := make(map[string]string)
	params["hash"] = strings.ToLower(infoHash)
	return c.post("command/resume", params)
}

//ResumeAll resumes all torrents
func (c *Client) ResumeAll(infoHashList []string) (*http.Response, error) {
	return c.get("command/resumeAll")
}

//ResumeMultiple resumes a list of torrents
func (c *Client) ResumeMultiple(infoHashList []string) (*http.Response, error) {
	params := c.processInfoHashList(infoHashList)
	return c.post("command/resumeAll", params)
}

//DeleteTemp deletes the temporary files for a list of torrents
func (c *Client) DeleteTemp(infoHashList []string) (*http.Response, error) {
	params := c.processInfoHashList(infoHashList)
	return c.post("command/delete", params)
}

//DeletePermanently deletes all files for a list of torrents
func (c *Client) DeletePermanently(infoHashList []string) (*http.Response, error) {
	params := c.processInfoHashList(infoHashList)
	return c.post("command/deletePerm", params)
}

//Recheck rechecks a list of torrents
func (c *Client) Recheck(infoHashList []string) (*http.Response, error) {
	params := c.processInfoHashList(infoHashList)
	return c.post("command/recheck", params)
}

//IncreasePriority increases the priority of a list of torrents
func (c *Client) IncreasePriority(infoHashList []string) (*http.Response, error) {
	params := c.processInfoHashList(infoHashList)
	return c.post("command/increasePrio", params)
}

//DecreasePriority decreases the priority of a list of torrents
func (c *Client) DecreasePriority(infoHashList []string) (*http.Response, error) {
	params := c.processInfoHashList(infoHashList)
	return c.post("command/decreasePrio", params)
}

//SetMaxPriority sets the max priority for a list of torrents
func (c *Client) SetMaxPriority(infoHashList []string) (*http.Response, error) {
	params := c.processInfoHashList(infoHashList)
	return c.post("command/topPrio", params)
}

//SetMinPriority sets the min priority for a list of torrents
func (c *Client) SetMinPriority(infoHashList []string) (*http.Response, error) {
	params := c.processInfoHashList(infoHashList)
	return c.post("command/bottomPrio", params)
}

//SetFilePriority sets the priority for a specific torrent file
func (c *Client) SetFilePriority(infoHash string, fileID string, priority string) (*http.Response, error) {
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
	return c.post("command/setFilePriority", params)
}

//GetGlobalDownloadLimit gets the global download limit of your qbittorrent client
func (c *Client) GetGlobalDownloadLimit() (limit int, err error) {
	var l int
	resp, err := c.get("command/getGlobalDlLimit")
	if err != nil {
		return l, err
	}
	json.NewDecoder(resp.Body).Decode(&l)
	return l, nil
}

//SetGlobalDownloadLimit sets the global download limit of your qbittorrent client
func (c *Client) SetGlobalDownloadLimit(limit string) (*http.Response, error) {
	params := make(map[string]string)
	params["limit"] = limit
	return c.post("command/setGlobalDlLimit", params)
}

//GetGlobalUploadLimit gets the global upload limit of your qbittorrent client
func (c *Client) GetGlobalUploadLimit() (limit int, err error) {
	var l int
	resp, err := c.get("command/getGlobalUpLimit")
	if err != nil {
		return l, err
	}
	json.NewDecoder(resp.Body).Decode(&l)
	return l, nil
}

//SetGlobalUploadLimit sets the global upload limit of your qbittorrent client
func (c *Client) SetGlobalUploadLimit(limit string) (*http.Response, error) {
	params := make(map[string]string)
	params["limit"] = limit
	return c.post("command/setGlobalUpLimit", params)
}

//GetTorrentDownloadLimit gets the download limit for a list of torrents
func (c *Client) GetTorrentDownloadLimit(infoHashList []string) (limits map[string]string, err error) {
	var l map[string]string
	params := c.processInfoHashList(infoHashList)
	resp, err := c.post("command/getTorrentsDlLimit", params)
	if err != nil {
		return l, err
	}
	json.NewDecoder(resp.Body).Decode(&l)
	return l, nil
}

//SetTorrentDownloadLimit sets the download limit for a list of torrents
func (c *Client) SetTorrentDownloadLimit(infoHashList []string, limit string) (*http.Response, error) {
	params := c.processInfoHashList(infoHashList)
	params["limit"] = limit
	return c.post("command/setTorrentsDlLimit", params)
}

//GetTorrentUploadLimit gets the upload limit for a list of torrents
func (c *Client) GetTorrentUploadLimit(infoHashList []string) (limits map[string]string, err error) {
	var l map[string]string
	params := c.processInfoHashList(infoHashList)
	resp, err := c.post("command/getTorrentsUpLimit", params)
	if err != nil {
		return l, err
	}
	json.NewDecoder(resp.Body).Decode(&l)
	return l, nil
}

//SetTorrentUploadLimit sets the upload limit of a list of torrents
func (c *Client) SetTorrentUploadLimit(infoHashList []string, limit string) (*http.Response, error) {
	params := c.processInfoHashList(infoHashList)
	params["limit"] = limit
	return c.post("command/setTorrentsUpLimit", params)
}

//SetPreferences sets the preferences of your qbittorrent client
func (c *Client) SetPreferences(params map[string]string) (*http.Response, error) {
	return c.postWithHeaders("command/setPreferences", params)
}

//GetAlternativeSpeedStatus gets the alternative speed status of your qbittorrent client
func (c *Client) GetAlternativeSpeedStatus() (status bool, err error) {
	var s bool
	resp, err := c.get("command/alternativeSpeedLimitsEnabled")
	if err != nil {
		return s, err
	}
	json.NewDecoder(resp.Body).Decode(&s)
	return s, nil
}

//ToggleAlternativeSpeed toggles the alternative speed of your qbittorrent client
func (c *Client) ToggleAlternativeSpeed() (*http.Response, error) {
	return c.get("command/toggleAlternativeSpeedLimits")
}

//ToggleSequentialDownload toggles the download sequence of a list of torrents
func (c *Client) ToggleSequentialDownload(infoHashList []string) (*http.Response, error) {
	params := c.processInfoHashList(infoHashList)
	return c.post("command/toggleSequentialDownload", params)
}

//ToggleFirstLastPiecePriority toggles first last piece priority of a list of torrents
func (c *Client) ToggleFirstLastPiecePriority(infoHashList []string) (*http.Response, error) {
	params := c.processInfoHashList(infoHashList)
	return c.post("command/toggleFirstLastPiecePrio", params)
}

//ForceStart force starts a list of torrents
func (c *Client) ForceStart(infoHashList []string, value bool) (*http.Response, error) {
	params := c.processInfoHashList(infoHashList)
	params["value"] = strconv.FormatBool(value)
	return c.post("command/setForceStart", params)
}
