package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

// SspResponse is convert to json
type SSPResponse struct {
	URL string `json:"url"`
}

// DspResponse is convert to json
type DSPResponse struct {
	RequestID string `json:"request_id"`
	URL       string `json:"url"`
	Price     int    `json:"price"`
}

// DspRequest is convert to json
type DSPpRequest struct {
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

// SdkRequest is convert to json
type SDKRequest struct {
	AppID int `json:"app_id"`
}

var hosts = os.Getenv("DSPHOSTS")

// HostArray is Split
var HostArray []string = strings.Split(hosts, " ")

var client = &http.Client{Timeout: time.Duration(100) * time.Millisecond}
var clientWin = &http.Client{Timeout: time.Duration(1000) * time.Millisecond}

func er(e error, errPoint string) {
	if e != nil {
		log.Println(errPoint, e)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	sdkreq := SDKRequest{}
	if r.Method == "POST" {
		data, err := ioutil.ReadAll(r.Body)
		er(err, "Post Request")
		json.Unmarshal(data, &sdkreq)
	} else {
		sdkreq.AppID = 123
	}
	count := len(HostArray) // hostの数に依存する
	id, _ := uuid.NewUUID()
	ids := id.String()

	dsprequest := DSPRequest{
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

			ch <- SendRequest(dsprequest, host)
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
		SendWinRequest(win, HostArray[0])

	} else {
		// ソートするやつ 数値以外が来たら終わる
		sort.Slice(auction, func(i, j int) bool { return auction[i].DspResponse.Price > auction[j].DspResponse.Price })
		// とりあえず一つに対して送る処理
		win := WinNotice{
			RequestID: ids,
			Price:     auction[1].DspResponse.Price,
		}
		
		// TODO これを修正したい
		SendWinRequest(win, auction[0].DspHost)
	}

	sspjson := SspResponse{URL: auction[0].DspResponse.URL}
	out, _ := json.Marshal(sspjson)
	outjson := string(out)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, outjson)
}

// SendRequest is dsp request
func SendRequest(dsprequest DspRequest, url string) PriceInfo {
	url = url + "/req"
	reqjson, _ := json.Marshal(dsprequest)
	req, _ := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer([]byte(reqjson)),
	)
	req.Header.Set("Content-type", "application/json")

	res, err := clientWin.Do(req)

	if res == nil || err != nil {
		//変に値が帰ってきても困るので
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
	return priceinfo,
}

// SendWinRequest is winnotice request
func SendWinRequest(win WinNotice, url string) {
	url = url + "/win"
	json, _ := json.Marshal(win)
	req, _ := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer([]byte(json)),
	)
	req.Header.Set("Content-type", "application/json")
	res, err := clientWin.Do(req)
	if err != nil {
		log.Println("にゃーん", err)
		return
	}

	// body, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	//fmt.Println(string(body))
}

func now() string {
	t := time.Now()
	str := fmt.Sprintf("%d%02d%02d-%02d%02d%02d.%04d", t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second(), time.Millisecond)
	return str
}

func main() {
	fmt.Println("unix server start")
	fmt.Println(HostArray)
	mux := http.NewServeMux()
	mux.HandleFunc("/req", handler)
	li, err := net.Listen("unix","/var/run/go/go.socket")
	if err != nil {
		panic(err)
	}

	err = http.Serve(li,mux)
	if err != nil {
		panic(err)
	}
	li.Close()

	// http.HandleFunc("/req", handler)
	// log.Fatalln(http.ListenAndServe(":8888", nil))
}
