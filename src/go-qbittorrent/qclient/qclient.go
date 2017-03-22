package qclient

import (
	"net/http"
	"fmt"
	"net/url"
	"go-qbittorrent/jsonHelpers"
	"strings"
	"strconv"
)

type Client struct {
	URL string
	Authenticated bool
	Session string //replace with session type
}

func NewClient(url string) *Client {
	c := &Client{}

	if url[len(url)-1:] != "/" {
		url = url + "/"
	}
	c.URL = url

	//find a way to check if authenticated

	return c
}

func (c Client) get(endpoint string) interface{} {
	httpClient := &http.Client{}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Println("error on creating request")
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("error on performing request")
	}

	body := jsonHelpers.RespToJson(resp)

	return jsonHelpers.JsonToStruct(body)
}

func (c Client) getWithParams(endpoint string, params map[string]string) interface {} {
	httpClient := &http.Client{}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Println("error on creating request")
	}

	//add parameters to url
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("error on performing request")
	}

	body := jsonHelpers.RespToJson(resp)

	return jsonHelpers.JsonToStruct(body)
}

func (c Client) post(endpoint string, data map[string]string) interface{} {
	httpClient := &http.Client{}

	req, err := http.NewRequest("POST", endpoint, nil)
	form := url.Values{}
	for k, v := range data {
		form.Add(k, v)
	}
	if err != nil {
		fmt.Println("error on creating request")
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("error on performing request")
	}

	body := jsonHelpers.RespToJson(resp)

	return jsonHelpers.JsonToStruct(body)
}

func (c Client) postWithHeaders(endpoint string, data map[string]string) interface{} {
	httpClient := &http.Client{}

	req, err := http.NewRequest("POST", endpoint, nil)

	req.Header.Set("content-type", "application/x-www-form-urlencoded")

	form := url.Values{}
	for k, v := range data {
		form.Add(k, v)
	}
	if err != nil {
		fmt.Println("error on creating request")
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("error on performing request")
	}

	body := jsonHelpers.RespToJson(resp)

	return jsonHelpers.JsonToStruct(body)
}


func (c Client) Login(username string, password string) {
	var creds map[string] string
	creds["username"] = username
	creds["password"] = password

	resp := c.post(c.URL+"login", creds)

	if resp["Response"] == "Ok" {
		c.Authenticated = true
	} else {
		fmt.Println("Not Authenticated")
	}
}


func (c Client) Logout() {
	resp := c.get("logout")

	if resp["Response"] == "Ok" {
		c.Authenticated = false
	} else {
		fmt.Println("Not Authenticated")
	}
}

func (c Client) Shutdown() interface{} {
	return c.get("command/shutdown")
}

func (c Client) Torrents(filters map[string]string ) interface{} {
	params := map[string]string{}
	for k, v := range filters {
		if k == "status" {
			k = "filter"
		}
		params[k] = v
	}
	return c.getWithParams("query/torrents", params)
}

func (c Client) Torrent(infoHash string) interface{} {
	return c.get("query/propertiesGeneral/" + strings.ToLower(infoHash))
}

func (c Client) TorrentTrackers(infoHash string) interface{} {
	return c.get("query/propertiesTrackers/" + strings.ToLower(infoHash))
}

func (c Client) TorrentWebSeeds(infoHash string) interface{} {
	return c.get("query/propertiesWebSeeds/" + strings.ToLower(infoHash))
}

func (c Client) TorrentFiles(infoHash string) interface{} {
	return c.get("query/propertiesFiles" + strings.ToLower(infoHash))
}

func (c Client) Sync(rid string) interface{} {
	params := map[string]string{}
	params["rid"] = rid
	return c.getWithParams("sync/maindata", params)
}

func (c Client) DownloadFromLink(link string, options map[string]string) interface{} {
	options["urls"] = link
	return c.post("command/download", options)
}

//TODO: implement a way to download using file buffer and new post method?
func (c Client) DownloadFromFile(fileName string, options map[string]string) interface{} {
	return c.post("command/download", options)
}

func (c Client) AddTrackers(infoHash string, trackers string) interface{} {
	params := map[string]string{}
	params["hash"] = strings.ToLower(infoHash)
	params["urls"] = trackers
	return c.post("command/addTrackers", params)
}

func (Client) processInfoHashList(infoHashList []string) map[string]string {
	d := map[string]string{}
	infoHash := ""
	for _, v := range infoHashList{
		infoHash = infoHash + "|" + v
	}
	d["hashes"] = infoHash
	return d
}

func (c Client) Pause(infoHash string) interface{} {
	params := map[string]string{}
	params["hash"] = strings.ToLower(infoHash)
	return c.post("command/pause", params)
}

func (c Client) PauseAll() interface{} {
	return c.get("command/pauseAll")
}

func (c Client) PauseMultiple(infoHashList []string) interface{} {
	params := c.processInfoHashList(infoHashList)
	return c.post("command/pauseAll", params)
}

func (c Client) SetLabel(infoHashList []string, label string) interface{} {
	params := c.processInfoHashList(infoHashList)
	params["label"] = label
	return c.post("command/setLabel", params)
}

func (c Client) SetCategory(infoHashList []string, category string) interface{} {
	params := c.processInfoHashList(infoHashList)
	params["category"] = category
	return c.post("command/setLabel", params)
}

func (c Client) Resume(infoHash string) interface{} {
	params := map[string]string{}
	params["hash"] = strings.ToLower(infoHash)
	return c.post("command/resume", params)
}

func (c Client) ResumeAll(infoHashList []string) interface{} {
	return c.get("command/resumeAll")
}

func (c Client) ResumeMultiple(infoHashList []string) interface{} {
	params := c.processInfoHashList(infoHashList)
	return c.post("command/resumeAll", params)
}

func (c Client) DeleteTemp(infoHashList []string) interface{} {
	params := c.processInfoHashList(infoHashList)
	return c.post("command/delete", params)
}

func (c Client) DeletePermanently(infoHashList []string) interface{} {
	params := c.processInfoHashList(infoHashList)
	return c.post("command/deletePerm", params)
}

func (c Client) Recheck(infoHashList []string) interface{} {
	params := c.processInfoHashList(infoHashList)
	return c.post("command/recheck", params)
}

func (c Client) IncreasePriority(infoHashList []string) interface{} {
	params := c.processInfoHashList(infoHashList)
	return c.post("command/increasePrio", params)
}

func (c Client) DecreasePriority(infoHashList []string) interface{} {
	params := c.processInfoHashList(infoHashList)
	return c.post("command/decreasePrio", params)
}

func (c Client) SetMaxPriority(infoHashList []string) interface{} {
	params := c.processInfoHashList(infoHashList)
	return c.post("command/topPrio", params)
}

func (c Client) SetMinPriority(infoHashList []string) interface{} {
	params := c.processInfoHashList(infoHashList)
	return c.post("command/bottomPrio", params)
}

func (c Client) SetFilePriority(infoHash string, file_id string, priority string) interface{} {
	priorities := [...]string{"0", "1", "2", "7"}
	for _, v := range priorities {
		if v == priority {
			fmt.Println("error, priority no tavailable")
		}
	}
	params := map[string]string{}
	params["hash"] = infoHash
	params["id"] = file_id
	params["priority"] = priority
	return c.post("command/setFilePriority", params)
}

func (c Client) GetGlobalDownloadLimit() interface{} {
	return c.get("command/getGlobalDlLimit")
}

func (c Client) SetGlobalDownloadLimit(limit string) interface{} {
	params := map[string]string{}
	params["limit"] = limit
	return c.post("command/setGlobalDlLimit", params)
}

func (c Client) GetGlobalUploadLimit() interface{} {
	return c.get("command/getGlobalUpLimit")
}

func (c Client) SetGlobalUploadLimit(limit string) interface{} {
	params := map[string]string{}
	params["limit"] = limit
	return c.post("command/setGlobalUpLimit", params)
}

func (c Client) GetTorrentDownloadLimit(infoHashList []string) interface{} {
	params := c.processInfoHashList(infoHashList)
	return c.post("command/getTorrentsDlLimit", params)
}

func (c Client) SetTorrentDownloadLimit(infoHashList []string, limit string) interface{} {
	params := c.processInfoHashList(infoHashList)
	params["limit"] = limit
	return c.post("command/setTorrentsDlLimit", params)
}

func (c Client) GetTorrentUploadLimit(infoHashList []string) interface{} {
	params := c.processInfoHashList(infoHashList)
	return c.post("command/getTorrentsUpLimit", params)
}

func (c Client) SetTorrentUploadLimit(infoHashList []string, limit string) interface{} {
	params := c.processInfoHashList(infoHashList)
	params["limit"] = limit
	return c.post("command/setTorrentsUpLimit", params)
}

func (c Client) SetPreferences(params map[string]string) interface{} {
	return c.postWithHeaders("command/setPreferences", params)
}

func (c Client) GetAlternativeSpeedStatus() interface{} {
	return c.get("command/alternativeSpeedLimitsEnabled")
}

func (c Client) ToggleAlternativeSpeed() interface{} {
	return c.get("command/toggleAlternativeSpeedLimits")
}

func (c Client) ToggleSequentialDownload(infoHashList []string) interface{} {
	params := c.processInfoHashList(infoHashList)
	return c.post("command/toggleSequentialDownload", params)
}

func (c Client) ToggleFirstLastPiecePriority(infoHashList []string) interface{} {
	params := c.processInfoHashList(infoHashList)
	return c.post("command/toggleFirstLastPiecePrio", params)
}

func (c Client) ForceStart(infoHashList []string, value bool) interface{} {
	params := c.processInfoHashList(infoHashList)
	params["value"] = strconv.FormatBool(value)
	return c.post("command/setForceStart", params)
}