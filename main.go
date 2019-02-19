package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("request start")

	ch := make(chan bool)
	for i := 0; i < 5; i++ {
		go func() {
			url := "http://localhost:8081"
			res, err := http.Get(url)
			if err != nil {
				fmt.Println(err)
			}
			res.Body.Close() // メッソドを見つけたからCloseしとくけどやらないと行けないかは謎
			ch <- true
		}()
	}

	for i := 0; i < 10; i++ {
		fmt.Println(<-ch)
	}
}
