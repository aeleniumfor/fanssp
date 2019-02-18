package main

import (
	"fmt"
	"net/http"
	"github.com/google/uuid"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World")
}

func main() {
	fmt.Println(uuid.NewUUID())
	fmt.Println("server start")
	http.HandleFunc("/", handler) // ハンドラを登録してウェブページを表示させる
	http.ListenAndServe(":8080", nil)
}
