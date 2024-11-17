package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	handler := NewHandler()
	handler.registerRoutes(mux)

	log.Printf("Starting HTTP server at 8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("Failed to start http server")
	}
}
