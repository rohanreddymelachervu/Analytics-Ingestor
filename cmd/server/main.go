package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/rohanreddymelachervu/ingestor/internal/config"
	"github.com/rohanreddymelachervu/ingestor/internal/server"
)

func main() {
	// Load configuration from environment
	cfg := config.Load()

	// Connect to Postgres
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Skip auto-migration since we have explicit migration files
	// The database schema is managed by the migration files in migrations/postgres/
	log.Println("Using existing database schema from migration files")

	// Initialize Gin router
	r := gin.Default()

	server.RegisterRoutes(r, cfg.JWTSecret, db)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}
