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

// WinNotice is convert to json
type WinNotice struct {
	RequestID string `json:"request_id"`
	Price     int    `json:"price"`
}

func handler(w http.ResponseWriter, r *http.Request) {

	count := 10
	var dspres []DspResponse
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
		dspres = append(dspres, dsp)
	}
	// ソートするやつ 数値以外が来たら終わる
	sort.Slice(dspres, func(i, j int) bool { return dspres[i].Price > dspres[j].Price })

	// とりあえず一つに対して送る処理
	win := WinNotice{
		RequestID: id.String(),
		Price:     dspres[1].Price,
	}

	winrequest(win)

	sspjson := SspResponse{dspres[0].URL}
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
	res.Body.Close() // メッソドを見つけたからCloseしとくけどやらないと行けないかは謎
	return body
}

func winrequest(win WinNotice) {
	url := "http://localhost:8085/win"

	json, _ := json.Marshal(win)
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
	fmt.Println(string(body))
	
}

func main() {
	fmt.Println("server start")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8081", nil)
}
