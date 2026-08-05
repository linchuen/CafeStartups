package main

import (
	"log"
	"net/http"
)

func main() {
	store := &gameStore{games: make(map[string]*gameRoom)}
	server := &http.Server{Addr: ":8080", Handler: newHandler(store)}
	log.Printf("Cafe Startups API listening on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
