package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
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
	AppID       int    `json:"app_id"`
}

// WinNotice is convert to json
type WinNotice struct {
	RequestID string `json:"request_id"`
	Price     int    `json:"price"`
}

// PriceInfo is convert to json
type PriceInfo struct {
	DspHost     string
	DspResponse DspResponse
	Status      bool
}

type SdkRequest struct {
	AppID int `json:"app_id"`
}

var hosts string = os.Getenv("DSPHOSTS")

// HostArray is Split
var HostArray []string = strings.Split(hosts, " ")

func er(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {

	sdkreq := SdkRequest{}
	if r.Method != "POST" {
		data, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(data, &sdkreq)
		log.Panicln(sdkreq)
	}else{
		sdkreq.AppID = 123
	}
	fmt.Println(sdkreq)

	count := len(HostArray) // hostの数に依存する
	id, _ := uuid.NewUUID()
	ids := id.String()

	dsprequest := DspRequest{
		SspName:     "r_ryusei",
		RequestTime: now(),
		RequestID:   ids,
		AppID:       sdkreq.AppID,
	}

	auction := []PriceInfo{}
	// DSPに対してリクエスを行う
	ch := make(chan PriceInfo, count)
	for _, host := range HostArray {
		go func(host string) {
			// HostArray[i]はurlの配列を一つ一つに分解したもの

			ch <- request(dsprequest, host)
		}(host)
	}

	for range HostArray {
		data := <-ch
		if data.Status == true {
			// レスポンスがきちんと帰ってきてる時
			auction = append(auction, data)

		}
	}
	if len(auction) == 0 {
		// dspのレスポンスが全てなかった場合
		dsp := DspResponse{
			RequestID: ids,
			URL:       "http://自社広告.コム:8080/ごめんね",
			Price:     0,
		}
		data := PriceInfo{
			DspResponse: dsp,
		}
		auction = append(auction, data)
	} else if len(auction) == 1 {
		// レスポンスが1つの場合
		win := WinNotice{
			RequestID: ids,
			Price:     1,
		}
		winrequest(win, HostArray[0])

	} else {
		// ソートするやつ 数値以外が来たら終わる
		sort.Slice(auction, func(i, j int) bool { return auction[i].DspResponse.Price > auction[j].DspResponse.Price })
		// とりあえず一つに対して送る処理
		win := WinNotice{
			RequestID: ids,
			Price:     auction[1].DspResponse.Price,
		}

		// TODO これを修正したい
		winrequest(win, HostArray[0])
	}

	sspjson := SspResponse{URL: auction[0].DspResponse.URL}
	out, _ := json.Marshal(sspjson)
	outjson := string(out)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, outjson)
}

func request(dsprequest DspRequest, url string) PriceInfo {
	reqjson, _ := json.Marshal(dsprequest)
	req, _ := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer([]byte(reqjson)),
	)
	req.Header.Set("Content-type", "application/json")

	client := &http.Client{Timeout: time.Duration(100) * time.Millisecond}
	res, err := client.Do(req)

	if res == nil || err != nil {
		//変に値が帰ってきても困るので
		er(err)
		return PriceInfo{Status: false}
	}

	dsp := DspResponse{}
	data, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(data, &dsp)
	res.Body.Close() // メッソドを見つけたからCloseしとくけどやらないと行けないかは謎

	priceinfo := PriceInfo{
		DspHost:     url,
		DspResponse: dsp,
		Status:      true,
	}
	return priceinfo
}

func winrequest(win WinNotice, url string) {
	json, _ := json.Marshal(win)
	req, _ := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer([]byte(json)),
	)
	req.Header.Set("Content-type", "application/json")

	client := &http.Client{}
	_, err := client.Do(req)
	if err != nil {
		log.Println("にゃーん")
	}
	// body, _ := ioutil.ReadAll(res.Body)
	// res.Body.Close()
	//fmt.Println(string(body))
}

func now() string {
	t := time.Now()
	str := fmt.Sprintf("%d%02d%02d-%02d%02d%02d.%04d", t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second(), time.Millisecond)
	return str
}

func main() {
	fmt.Println("server start")
	fmt.Println(HostArray)
	http.HandleFunc("/req", handler)
	http.ListenAndServe(":8888", nil)
}
