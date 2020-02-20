package main

import (
	"fmt"
	"time"

	"./net/http"
)

func main() {
	fmt.Println("Welcome in ImpetusResel")
	api := http.NewRouter()
	api.AddRoute("/bonjour", func(w http.Headers, r *http.Request) {
		w.SetStatusCode(200)
		w.AddEntity(http.ContentType, "text/plain; charset=utf-8")
		w.SetBody("Welcome you are on this page: " + r.URL)
		fmt.Println("1 sec sleep ->", string(w.Bytes()))
		time.Sleep(10 * time.Second)
	})
	api.SetDefaultRoute(func(w http.Headers, r *http.Request) {
		w.SetStatusCode(404)
		w.AddEntity(http.ContentType, "text/plain; charset=utf-8")
		w.SetBody("Page not found\n")
	})
	http.ListenAndServe(8085, api)

	// req, err := http.NewRequest("GET", "localhost/bonjour", []byte{}) // http://www.googleapis.com/books/v1/volumes?q=isbn:0747532699
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(string(req.Bytes()))
	// resp, err := http.Do(&req)
	// if err != nil {
	// 	return
	// }
	// resp.Print()
}
