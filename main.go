package main

import (
	"forum/internal/database"
	"forum/internal/server/servsetup"
	"forum/internal/server/templates"
	"log"
)

func main() {
	//DB initialization
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	//Initialize server
	err := templates.LoadTemplates()
	if err != nil {
		log.Fatalf("Failed to load templates: %v", err)
	}
	servsetup.LaunchServer()
}
