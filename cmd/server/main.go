package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/rohanreddymelachervu/ingestor/internal/auth"
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

	// Initialize Gin router
	r := gin.Default()

	// Setup routes
	authService := auth.NewService(db, cfg.JWTSecret)
	authHandler := auth.NewHandler(authService)
	server.RegisterRoutes(r, authHandler, cfg.JWTSecret)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
