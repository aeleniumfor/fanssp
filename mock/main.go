package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

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
	
	data, _ := ioutil.ReadAll(r.Body)
	sspReq := DspRequest{}
	json.Unmarshal(data, &sspReq)
	price := randTOint()
	time.Sleep(time.Duration(0) * time.Millisecond)
	dspjson := DspResponse{}
	dspjson.RequestID = sspReq.RequestID
	dspjson.URL = "http://hoge.com/" + strconv.Itoa(price)
	dspjson.Price = price

	fmt.Println(dspjson)
	out, _ := json.Marshal(dspjson)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(out))
}

func winNotice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	out := `{"result": "ok"}`

	fmt.Fprintf(w,out)
}

func randTOint() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(150)
}

func now() string {
	t := time.Now()
	return t.String()
}

func main() {
	fmt.Println("start mock server")
	http.HandleFunc("/", handler)
	http.HandleFunc("/win", winNotice)
	http.ListenAndServe(":8080", nil)
}
