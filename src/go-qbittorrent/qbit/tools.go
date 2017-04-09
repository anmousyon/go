package qbit

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

func PrintResponse(resp *http.Response) {
	r := make([]byte, 256)
	r, _ = ioutil.ReadAll(resp.Body)
	fmt.Println("response: " + string(r))
}

func PrintRequest(req *http.Request) {
	r, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("request: " + string(r))
}
