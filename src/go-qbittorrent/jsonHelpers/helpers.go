package jsonHelpers

import (
	"encoding/json"
	"net/http"
	"io/ioutil"
	"fmt"
)

func JsonToStruct(body []byte) interface{} {
	var u interface{}

	if len(body) == 0 {
		json.Unmarshal([]byte{}, &u)
	} else {
		json.Unmarshal(body, &u)
	}

	return u
}

func RespToJson(resp *http.Response) []byte {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error on response read")
	}
	return body
}
