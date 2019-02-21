package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/fanssp/common"
	"github.com/google/uuid"
)

var hosts string = os.Getenv("DSPHOSTS")

// HostArray is Split
var HostArray []string = strings.Split(hosts, " ")


func handler(w http.ResponseWriter, r *http.Request) {
	count := len(HostArray) // hostの数に依存する
	dspres := []common.DspResponse{}
	id, _ := uuid.NewUUID()
	ids := id.String()

	dsprequest := common.DspRequest{
		SspName:     "r_ryusei",
		RequestTime: now(),
		RequestID:   ids,
		AppID:       123,
	}

	// DSPに対してリクエスを行う
	ch := make(chan []byte, count)
	for _ ,host := range HostArray {
		go func(host string) {
			// HostArray[i]はurlの配列を一つ一つに分解したもの

			ch <- request(dsprequest, host)
		}(host)
	}

	for range HostArray {
		dsp := common.DspResponse{}
		data := <-ch
		if len(data) != 0 {
			json.Unmarshal(data, &dsp)
			dspres = append(dspres, dsp)

		}
	}
	if len(dspres) == 0 {
		// dspのレスポンスが全てなかった場合
		dsp := common.DspResponse{
			RequestID: ids,
			URL:       "http://自社広告.コム:8080/ごめんね",
			Price:     0,
		}
		dspres = append(dspres, dsp)
	} else if len(dspres) == 1 {
		// レスポンスが1つの場合
		win := common.WinNotice{
			RequestID: ids,
			Price:     1,
		}
		winrequest(win, HostArray[0])

	} else {
		// ソートするやつ 数値以外が来たら終わる
		sort.Slice(dspres, func(i, j int) bool { return dspres[i].Price > dspres[j].Price })
		// とりあえず一つに対して送る処理
		win := common.WinNotice{
			RequestID: ids,
			Price:     dspres[1].Price,
		}
		winrequest(win, HostArray[0])
	}

	sspjson := common.SspResponse{dspres[0].URL}
	out, _ := json.Marshal(sspjson)
	outjson := string(out)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, outjson)
}

func request(dsprequest common.DspRequest, url string) []byte {
	json, _ := json.Marshal(dsprequest)
	req, _ := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer([]byte(json)),
	)
	req.Header.Set("Content-type", "application/json")

	client := &http.Client{Timeout: time.Duration(100) * time.Millisecond}
	res, _ := client.Do(req)
	if res == nil {
		return []byte{}
	}
	body, _ := ioutil.ReadAll(res.Body)
	res.Body.Close() // メッソドを見つけたからCloseしとくけどやらないと行けないかは謎
	return body
}

func winrequest(win common.WinNotice, url string) {
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
