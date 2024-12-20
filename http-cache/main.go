package main

import (
	"log"
	"net/http"
	"strings"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /greet", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s: %s, RemoteAddr: %s\n", r.Method, r.URL.Path, r.RemoteAddr)
		e := "greet"
		w.Header().Set("ETag", e)
		w.Header().Set("Cache-Control", "max-age=300")
		if match := r.Header.Get("If-None-Match"); match != "" {
			if strings.Contains(match, e) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}
		w.Write([]byte("Hello, World"))
	})

	log.Println("Listening on 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
