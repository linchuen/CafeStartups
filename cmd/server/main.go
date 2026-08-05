package main

import (
	"log"
	"net/http"

	"cafestartups/internal/server"
)

func main() {
	handler := server.NewHandler()
	httpServer := &http.Server{Addr: ":8080", Handler: handler}
	log.Printf("Cafe Startups API listening on %s", httpServer.Addr)
	log.Fatal(httpServer.ListenAndServe())
}
