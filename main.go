package main

import (
	"Forum/server"
	"log"
	"net/http"
)

func main() {
	mux := server.InitializeServer()

	log.Println("Starting server on : 8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
