package main

import (
	"forum/internal/database"
	"forum/internal/server/routes"
	"log"
	"net/http"
)

func main() {
	//DB initialization
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	//Initialize server
	mux := routes.SetupRoutes()

	log.Println("Starting server on : 8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
