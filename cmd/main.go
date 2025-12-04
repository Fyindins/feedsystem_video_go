package main

import (
	"feedsystem_video_go/internal/config"
	"feedsystem_video_go/internal/database"
	"feedsystem_video_go/internal/router"
	"log"
)

func main() {
	// Load config
	log.Printf("Loading config from configs/config.yaml")
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	// Connect database
	log.Printf("Database config: %v", cfg.Database)

	db, err := database.NewDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
	}
	defer database.CloseDB(db)
	// Set router
	r := router.SetRouter(db)
	log.Printf("Server is running on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
