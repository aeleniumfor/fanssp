package main

import (
	"github.com/fanssp/common"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	
	data, _ := ioutil.ReadAll(r.Body)
	sspReq := common.DspRequest{}
	json.Unmarshal(data, &sspReq)
	price := randTOint()
	time.Sleep(time.Duration(100) * time.Second)
	dspjson := common.DspResponse{}
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
