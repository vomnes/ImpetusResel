package main

import (
	"fmt"
	"time"

	"./net/http"
)

func main() {
	fmt.Println("Welcome in ImpetusResel")
	api := http.NewRouter()
	api.AddRoute("/bonjour", func(w *http.Headers, r *http.Request) {
		w.SetStatusCode(200)
		w.AddEntity(http.ContentType, "text/plain; charset=utf-8")
		w.SetBody("Welcome you are on this page: " + r.URL)
		fmt.Println("1 sec sleep")
		time.Sleep(10 * time.Second)
	})
	api.SetDefaultRoute(func(w *http.Headers, r *http.Request) {
		w.SetStatusCode(404)
		w.AddEntity(http.ContentType, "text/plain; charset=utf-8")
		w.SetBody("Page not found\n")
	})
	http.ListenAndServe(8085, api)
}
