package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

// DspResponse is convert to json
type DspResponse struct {
	RequestID string `json:"request_id"`
	URL       string `json:"url"`
	Price     string `json:"price"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.NewUUID()

	dspjson := DspResponse{}
	dspjson.RequestID = id.String()
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
