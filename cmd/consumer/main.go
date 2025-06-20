package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/rohanreddymelachervu/ingestor/internal/config"
	"github.com/rohanreddymelachervu/ingestor/internal/events"
	"github.com/rohanreddymelachervu/ingestor/internal/kafka"
	"github.com/rohanreddymelachervu/ingestor/internal/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	log.Println("🚀 Starting Kafka Event Consumer...")

	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize repositories
	eventRepo := repository.NewEventRepository(db)
	quizRepo := repository.NewQuizRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	classroomRepo := repository.NewClassroomRepository(db)

	// Initialize event service
	eventService := events.NewService(eventRepo, quizRepo, sessionRepo, classroomRepo)

	// Kafka configuration
	kafkaBrokers := getKafkaBrokers()
	groupID := "analytics-event-processors"
	topics := []string{"quiz-events"}

	log.Printf("Connecting to Kafka brokers: %v", kafkaBrokers)
	log.Printf("Consumer group: %s", groupID)
	log.Printf("Topics: %v", topics)

	// Create Kafka consumer
	consumer, err := kafka.NewConsumer(kafkaBrokers, groupID, topics, eventService)
	if err != nil {
		log.Fatal("Failed to create Kafka consumer:", err)
	}

	// Start consuming
	ctx := context.Background()
	log.Println("📨 Starting event consumption...")

	if err := consumer.Start(ctx); err != nil {
		log.Fatal("Consumer failed:", err)
	}

	log.Println("🛑 Consumer stopped")
}

func getKafkaBrokers() []string {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "localhost:9092" // Default for local development
	}
	return strings.Split(brokers, ",")
}
