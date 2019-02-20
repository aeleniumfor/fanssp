package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("request start")
	var hosts string = os.Getenv("DSPHOSTS")
	var HostArray []string = strings.Split(hosts, " ")
	fmt.Println(HostArray)
	for i := 0; i < len(HostArray); i++ {
		fmt.Println(HostArray[i])
	}
	// ch := make(chan bool)
	// for i := 0; i < 5; i++ {
	// 	go func() {
	// 		url := "http://localhost:8081"
	// 		res, err := http.Get(url)
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	// 		res.Body.Close() // メッソドを見つけたからCloseしとくけどやらないと行けないかは謎
	// 		ch <- true
	// 	}()
	// }

	// for i := 0; i < 10; i++ {
	// 	fmt.Println(<-ch)
	// }
}
