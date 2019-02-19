package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"sort"
)

// SspResponse is convert to json
type SspResponse struct {
	URL string `json:"url"`
}

// DspResponse is convert to json
type DspResponse struct {
	RequestID string `json:"request_id"`
	URL       string `json:"url"`
	Price     int    `json:"price"`
}

// DspRequest is convert to json
type DspRequest struct {
	SspName     string `json:"ssp_name"`
	RequestTime string `json:"request_time"`
	RequestID   string `json:"request_id"`
}

func handler(w http.ResponseWriter, r *http.Request) {

	count := 10
	var dspres [10]DspResponse
	id, _ := uuid.NewUUID()

	dsprequest := DspRequest{
		SspName:     "hoge",
		RequestTime: "time",
		RequestID:   id.String(),
	}

	// DSPに対してリクエスを行う
	ch := make(chan []byte)
	for i := 0; i < count; i++ {
		go func() {
			ch <- request(dsprequest)
		}()
	}

	for i := 0; i < count; i++ {

		dsp := DspResponse{}
		json.Unmarshal(<-ch, &dsp)
		dspres[i] = dsp
	}

	dspresslice := dspres[:]
	// ソートするやつ 数値以外が来たら終わる
	sort.Slice(dspresslice, func(i, j int) bool { return dspresslice[i].Price > dspresslice[j].Price })
	
	sspjson := SspResponse{"http://hoge.example.com"}
	out, _ := json.Marshal(sspjson)
	outjson := string(out)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, outjson)
}

func request(dsprequest DspRequest) []byte {
	url := "http://localhost:8085"

	json, _ := json.Marshal(dsprequest)
	req, _ := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer([]byte(string(json))),
	)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	body, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	return body
}

func main() {
	fmt.Println("server start")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8081", nil)
}
