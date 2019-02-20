package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

var hosts string = os.Getenv("DSPHOSTS")

// HostArray is Split
var HostArray []string = strings.Split(hosts, " ")

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
	AppID       int `json:"app_id"`
}

// WinNotice is convert to json
type WinNotice struct {
	RequestID string `json:"request_id"`
	Price     int    `json:"price"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get request")
	count := len(HostArray) // hostの数に依存する
	var dspres []DspResponse
	id, _ := uuid.NewUUID()

	dsprequest := DspRequest{
		SspName:     "hoge",
		RequestTime: "time",
		RequestID:   id.String(),
		AppID:       123,
	}

	// DSPに対してリクエスを行う
	ch := make(chan []byte, 5)
	for i := 0; i < count; i++ {
		go func(i int) {
			// HostArray[i]はurlの配列を一つ一つに分解したもの
			ch <- request(dsprequest, HostArray[i])
		}(i)
	}

	for i := 0; i < count; i++ {
		dsp := DspResponse{}
		data := <-ch
		if len(data) != 0 {
			json.Unmarshal(data, &dsp)
			dspres = append(dspres, dsp)

		}
	}
	if len(dspres) == 0 {
		// dspのレスポンスが全てなかった場合
		dsp := DspResponse{
			RequestID: id.String(),
			URL:       "http://自社広告.コム:8080/ごめんね",
			Price:     0,
		}
		dspres = append(dspres, dsp)
	} else if len(dspres) == 1 {
		// レスポンスが1つの場合
		win := WinNotice{
			RequestID: id.String(),
			Price:     1,
		}
		winrequest(win, HostArray[0])

	} else {
		// ソートするやつ 数値以外が来たら終わる
		sort.Slice(dspres, func(i, j int) bool { return dspres[i].Price > dspres[j].Price })
		// とりあえず一つに対して送る処理
		win := WinNotice{
			RequestID: id.String(),
			Price:     dspres[1].Price,
		}
		winrequest(win, HostArray[0])
	}

	sspjson := SspResponse{dspres[0].URL}
	out, _ := json.Marshal(sspjson)
	outjson := string(out)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, outjson)
}

func request(dsprequest DspRequest, url string) []byte {
	json, _ := json.Marshal(dsprequest)
	req, _ := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer([]byte(json)),
	)

	client := &http.Client{Timeout: time.Duration(100) * time.Millisecond}
	res, _ := client.Do(req)
	fmt.Println(res)

	if res == nil {
		return []byte{}
	}
	body, _ := ioutil.ReadAll(res.Body)
	res.Body.Close() // メッソドを見つけたからCloseしとくけどやらないと行けないかは謎
	return body
}

func winrequest(win WinNotice, url string) {
	json, _ := json.Marshal(win)
	req, _ := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer([]byte(string(json))),
	)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("にゃーん")
	}
	body, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	fmt.Println(string(body))

}

func main() {
	fmt.Println("server start")
	fmt.Println(HostArray)
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8081", nil)
}
