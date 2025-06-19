package server

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/rohanreddymelachervu/ingestor/internal/auth"
	"github.com/rohanreddymelachervu/ingestor/internal/events"
	"github.com/rohanreddymelachervu/ingestor/internal/reports"
	"github.com/rohanreddymelachervu/ingestor/internal/repository"
)

// RegisterRoutes sets up all endpoints with proper clean architecture
func RegisterRoutes(r *gin.Engine, authHandler *auth.Handler, jwtSecret string, db *gorm.DB) {
	// Initialize repositories
	eventRepo := repository.NewEventRepository(db)
	quizRepo := repository.NewQuizRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	classroomRepo := repository.NewClassroomRepository(db)

	// Initialize services
	eventsService := events.NewService(eventRepo, quizRepo, sessionRepo, classroomRepo)
	reportsService := reports.NewService(eventRepo)

	// Initialize handlers
	eventsHandler := events.NewHandler(eventsService)
	reportsHandler := reports.NewHandler(reportsService)

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
			}
		}
	}
}
