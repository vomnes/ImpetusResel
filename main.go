package main

import (
	"fmt"

	"./http"
)

func main() {
	fmt.Println("Welcome in ImpetusResel")
	r := http.NewRequest()
	r.RequestParse("GET / HTTP/1.1\r\nHost: localhost:8084\r\nUser-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:74.0) Gecko/20100101 Firefox/74.0\r\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8\r\nAccept-Language: en-US,en;q=0.5\r\nAccept-Encoding: gzip, deflate\r\nConnection: keep-alive\r\nUpgrade-Insecure-Requests: 1\r\n\r\n")
	r.Print()
	// http.ListenAndServe(8084)
}
