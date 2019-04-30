package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type SdkRequest struct {
	AppID int `json:"app_id"`
}

var client = &http.Client{}

func request(url string) int {
	reqestStruct := SdkRequest{AppID: 123}
	reqjson, _ := json.Marshal(reqestStruct)
	req, _ := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer([]byte(reqjson)),
	)
	req.Header.Set("Content-type", "application/json")
	res, _ := client.Do(req)
	return res.StatusCode
}

func main() {
	url := "http://localhost:8888/req"
	reqestStruct := SdkRequest{AppID: 123}
	reqjson, _ := json.Marshal(reqestStruct)
	req, _ := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer([]byte(reqjson)),
	)
	req.Header.Set("Content-type", "application/json")
	res, _ := client.Do(req)
	log.Println(res.StatusCode)
}
