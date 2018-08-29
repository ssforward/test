package main

import (
	"net/http"
	"./handler"
	_ "github.com/lib/pq"
)

func main() {
	server := http.Server{
		Addr: "localhost:8080",
	}

	http.HandleFunc("/top", handler.TopHandler)

	server.ListenAndServe()
}