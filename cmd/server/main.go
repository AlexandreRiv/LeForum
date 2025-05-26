package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/templates/index.html")
	})

	log.Println("LeForum running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
