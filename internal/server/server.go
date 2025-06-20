package server

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/rohanreddymelachervu/ingestor/internal/auth"
	"github.com/rohanreddymelachervu/ingestor/internal/events"
	"github.com/rohanreddymelachervu/ingestor/internal/kafka"
	"github.com/rohanreddymelachervu/ingestor/internal/reports"
	"github.com/rohanreddymelachervu/ingestor/internal/repository"
)

// RegisterRoutes sets up all endpoints with proper clean architecture
func RegisterRoutes(r *gin.Engine, jwtSecret string, db *gorm.DB) {
	// Initialize repositories
	eventRepo := repository.NewEventRepository(db)
	quizRepo := repository.NewQuizRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	classroomRepo := repository.NewClassroomRepository(db)

	// Initialize services
	eventsService := events.NewService(eventRepo, quizRepo, sessionRepo, classroomRepo)
	reportsService := reports.NewService(eventRepo, classroomRepo)
	authService := auth.NewService(db, jwtSecret)

	// Initialize events handler (with or without Kafka)
	var eventsHandler *events.Handler

	// Check if Kafka mode is enabled
	useKafka := os.Getenv("KAFKA_ENABLED") == "true"

	if useKafka {
		log.Println("🚀 Kafka mode enabled - events will be published to Kafka")

		// Get Kafka configuration
		kafkaBrokers := getKafkaBrokers()
		topicName := getKafkaTopic()

		log.Printf("Kafka brokers: %v", kafkaBrokers)
		log.Printf("Kafka topic: %s", topicName)

		// Initialize Kafka producer
		producer, err := kafka.NewProducer(kafkaBrokers, topicName)
		if err != nil {
			log.Printf("⚠️  Failed to initialize Kafka producer: %v", err)
			log.Println("🔄 Falling back to direct database mode")
			eventsHandler = events.NewHandler(eventsService)
		} else {
			log.Println("✅ Kafka producer initialized successfully")
			eventsHandler = events.NewHandlerWithKafka(eventsService, producer)
		}
	} else {
		log.Println("📊 Direct database mode - events will be processed immediately")
		eventsHandler = events.NewHandler(eventsService)
	}

	// Initialize other handlers
	reportsHandler := reports.NewHandler(reportsService)
	authHandler := auth.NewHandler(authService)

	api := r.Group("/api")
	{
		// Public auth endpoints (no JWT required)
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/signup", authHandler.SignUp)
			authGroup.POST("/login", authHandler.Login)
		}

		// Secured routes require a valid JWT
		secured := api.Group("")
		secured.Use(auth.AuthMiddleware(jwtSecret))
		{
			// Event ingestion: WRITE scope required (for Whiteboard & Notebook apps)
			eventsGroup := secured.Group("")
			eventsGroup.Use(auth.RequireScope("WRITE"))
			{
				eventsGroup.POST("/events", eventsHandler.CreateEvent)
				eventsGroup.POST("/events/batch", eventsHandler.CreateBatchEvents)
			}

			// Reporting: READ scope required (for Analytics Dashboard)
			reportsGroup := secured.Group("/reports")
			reportsGroup.Use(auth.RequireScope("READ"))
			{
				reportsGroup.GET("/active-participants", reportsHandler.GetActiveParticipants)
				reportsGroup.GET("/questions-per-minute", reportsHandler.GetQuestionsPerMinute)
				reportsGroup.GET("/student-performance", reportsHandler.GetStudentPerformance)
				reportsGroup.GET("/classroom-engagement", reportsHandler.GetClassroomEngagement)
				reportsGroup.GET("/content-effectiveness", reportsHandler.GetContentEffectiveness)
				reportsGroup.GET("/response-rate", reportsHandler.GetResponseRate)
				reportsGroup.GET("/latency-analysis", reportsHandler.GetLatencyAnalysis)
				reportsGroup.GET("/timeout-analysis", reportsHandler.GetTimeoutAnalysis)
				reportsGroup.GET("/completion-rate", reportsHandler.GetCompletionRate)
				reportsGroup.GET("/dropoff-analysis", reportsHandler.GetDropoffAnalysis)
				reportsGroup.GET("/student-performance-list", reportsHandler.GetStudentPerformanceList)
				reportsGroup.GET("/classroom-engagement-history", reportsHandler.GetClassroomEngagementHistory)
			}
		}
	}
}

func getKafkaBrokers() []string {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "localhost:9092" // Default for local development
	}
	return strings.Split(brokers, ",")
}

func getKafkaTopic() string {
	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		topic = "quiz-events" // Default topic name
	}
	return topic
}
