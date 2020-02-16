package main

import (
	"fmt"

	"./http"
)

func main() {
	fmt.Println("Welcome in ImpetusResel")
	api := http.NewRouter()
	api.AddRoute("/bonjour", func(w *http.Headers, r *http.Request) {
		w.SetStatusCode(200)
		w.AddEntity(http.ContentType, "text/plain; charset=utf-8")
		w.SetBody("Welcome you are on this page: " + r.URL)
	})
	api.SetDefaultRoute(func(w *http.Headers, r *http.Request) {
		w.SetStatusCode(404)
		w.AddEntity(http.ContentType, "text/plain; charset=utf-8")
		w.SetBody("Page not found")
	})
	http.ListenAndServe(8084, api)
}
