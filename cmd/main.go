package main

import (
	"log"
	"net/http"

	pkg "forum-v1/pkg"
)

func main() {
	// mux := http.NewServeMux()
	// mux.Handle("/", home)
	mux := http.NewServeMux()
	mux.HandleFunc("/", pkg.Home)
	mux.Handle("/image/", http.StripPrefix("/image", http.FileServer(http.Dir("./web/image"))))
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./web/static"))))
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
