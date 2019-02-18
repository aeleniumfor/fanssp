package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
)

// SspResponse is convert to json
type SspResponse struct {
	URL string `json:"url"`
}

// DspRequest is convert to json
type DspRequest struct {
	SspName     string `json:"ssp_name"`
	RequestTime string `json:"request_time"`
	RequestID   string `json:"request_id"`
}

func handler(w http.ResponseWriter, r *http.Request) {

	id, _ := uuid.NewUUID()
	sspjson := SspResponse{"http://hoge.example.com"}

	dsprequest := DspRequest{
		SspName:     "hoge",
		RequestTime: "time",
		RequestID:   id.String(),
	}
	request(dsprequest)

	out, _ := json.Marshal(sspjson)
	outjson := string(out)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, outjson)
}

func request(dsprequest DspRequest) {
	url := "http://localhost:8080"

	json, _ := json.Marshal(dsprequest)
	req, _ := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer([]byte(string(json))),
	)

	client := &http.Client{}
	res, _ := client.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))

}

func main() {
	fmt.Println("server start")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8081", nil)
}
