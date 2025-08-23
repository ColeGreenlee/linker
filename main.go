package main

import (
	"log"

	"linker/internal/api"
	"linker/internal/config"
	"linker/internal/database"
)

func main() {
	cfg := config.Load()
	
	db, err := database.Init(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	server := api.NewServer(cfg, db)
	
	log.Printf("Starting server on port %s", cfg.Port)
	if err := server.Start(); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}