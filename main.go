package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

// SspResponse is convert to json
type SspResponse struct {
	URL string `json:"url"`
}

func handler(w http.ResponseWriter, r *http.Request) {

	id, _ := uuid.NewUUID()
	fmt.Println(id)
	sspjson := SspResponse{}
	sspjson.URL = "http://hoge.example.com"
	out, _ := json.Marshal(sspjson)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(out))
}

func main() {
	fmt.Println("server start")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
