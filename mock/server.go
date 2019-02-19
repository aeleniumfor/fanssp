package main

import (
	"io/ioutil"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
)

// DspResponse is convert to json
type DspResponse struct {
	RequestID string `json:"request_id"`
	URL       string `json:"url"`
	Price     string `json:"price"`
}


// DspRequest is convert to json
type DspRequest struct {
	SspName     string `json:"ssp_name"`
	RequestTime string `json:"request_time"`
	RequestID   string `json:"request_id"`
}

func now() string {
	t := time.Now()
	return t.String()

}

func handler(w http.ResponseWriter, r *http.Request) {

	data, _ := ioutil.ReadAll(r.Body)
	
	sspReq := DspRequest{}
	json.Unmarshal(data,&sspReq)
	fmt.Println(sspReq)

	dspjson := DspResponse{}
	dspjson.RequestID = sspReq.RequestID
	dspjson.URL = "http://hoge.com"
	dspjson.Price = "50"
	out, _ := json.Marshal(dspjson)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(out))
}

func main() {
	fmt.Println("server start")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
