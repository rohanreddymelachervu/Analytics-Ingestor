package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/rohanreddymelachervu/ingestor/internal/auth"
)

// RegisterRoutes sets up all endpoints with auth middleware
// RegisterRoutes sets up all endpoints with proper middleware placement
func RegisterRoutes(r *gin.Engine, authHandler *auth.Handler, jwtSecret string) {
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
			// Event ingestion: WRITE scope required
			events := secured.Group("")
			events.Use(auth.RequireScope("WRITE"))
			{
				events.POST("/events", handleCreateEvent)
				events.POST("/events/batch", handleCreateBatchEvents)
			}

			// Reporting: READ scope required
			reports := secured.Group("/reports")
			reports.Use(auth.RequireScope("READ"))
			{
				reports.GET("/active-participants", handleActiveParticipants)
				reports.GET("/questions-per-minute", handleQuestionsPerMinute)
			}
		}
	}
}

// Placeholder handlers - implement your business logic here
func handleCreateEvent(c *gin.Context) {
	// Extract user ID from context (set by AuthMiddleware)
	userID, _ := c.Get("userID")
	
	c.JSON(http.StatusCreated, gin.H{
		"message": "Event created successfully",
		"user_id": userID,
	})
}

func handleCreateBatchEvents(c *gin.Context) {
	userID, _ := c.Get("userID")
	
	c.JSON(http.StatusCreated, gin.H{
		"message": "Batch events created successfully", 
		"user_id": userID,
	})
}

func handleActiveParticipants(c *gin.Context) {
	userID, _ := c.Get("userID")
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Active participants report",
		"user_id": userID,
		"data": []string{"participant1", "participant2"},
	})
}

func handleQuestionsPerMinute(c *gin.Context) {
	userID, _ := c.Get("userID")
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Questions per minute report",
		"user_id": userID, 
		"data": map[string]int{"questions_per_minute": 5},
	})
}
