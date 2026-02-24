package main

import (
	"log"
	"net/http"
)

func main() {
	addr := ":8080"
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, NewHTTPHandler()); err != nil {
		log.Fatal(err)
	}
}
