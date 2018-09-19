package main

import (
	"net"
    	"net/http"
	"./handler"
    	"net/http/fcgi"

)

func main() {
	l, err := net.Listen("tcp", "127.0.0.1:9000")
    	if err != nil {
        	return
    	}

	http.HandleFunc("/top", handler.TopHandler)
	fcgi.Serve(l, nil)
}
