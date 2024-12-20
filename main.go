package main

import (
	"forum/internal/config"
	"forum/internal/database"
	"forum/internal/server/servsetup"
	"forum/internal/server/templates"
	"forum/internal/utils"
	"log"
)

func main() {
	utils.LoadEnvFile(config.EnvFilePath)
	config.EnvVar.Set()
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	err := templates.LoadTemplates()
	if err != nil {
		log.Fatalf("Failed to load templates: %v", err)
	}
	servsetup.LaunchServer()
}
