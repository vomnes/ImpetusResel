package main

import (
	"fmt"

	"./http"
)

func main() {
	fmt.Println("Welcome in ImpetusResel")
	r := http.NewRequest()
	// r.RequestParse("GET / HTTP/1.1\r\nHost: localhost:8084\r\nUser-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:74.0) Gecko/20100101 Firefox/74.0\r\nContent-Length: 42\r\nContent-Type: application/x-www-form-urlencoded\r\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8\r\nAccept-Language: en-US,en;q=0.5\r\nAccept-Encoding: gzip, deflate\r\nConnection: keep-alive\r\nUpgrade-Insecure-Requests: 1\r\n\r\nfield1=value1&field2=value2&field2=Hello+world !")
	r.RequestParse("POST / HTTP/1.1\r\nHost: localhost:8084\r\nUser-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:74.0) Gecko/20100101 Firefox/74.0\r\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8\r\nAccept-Language: en-US,en;q=0.5\r\nAccept-Encoding: gzip, deflate\r\nContent-Type: multipart/form-data; boundary=\"---------------------------20762440193078419611623191500\"\r\nContent-Length: 772\r\nConnection: keep-alive\r\nUpgrade-Insecure-Requests: 1\r\n\r\n-----------------------------20762440193078419611623191500\r\nContent-Disposition: form-data; name=\"text\"\r\n\r\nvalentin omnes\r\n-----------------------------20762440193078419611623191500\r\nContent-Disposition: form-data; name=\"file1\"; filename=\"\"\r\nContent-Type: application/octet-stream\r\n\r\n\r\n-----------------------------20762440193078419611623191500\r\nContent-Disposition: form-data; name=\"file2\"; filename=\"http.html\"\r\nContent-Type: text/html\r\n\r\n<form action=\"http://localhost:8084\" method=\"post\" enctype=\"multipart/form-data\">\n  <p><input type=\"text\" name=\"text\" value=\"text default\">\n  <p><input type=\"file\" name=\"file1\">\n  <p><input type=\"file\" name=\"file2\">\n  <p><button type=\"submit\">Submit</button>\n</form>\n\r\n-----------------------------20762440193078419611623191500--\r\n")
	r.Print()
	// http.ListenAndServe(8084)
}
