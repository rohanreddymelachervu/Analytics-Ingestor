package events

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

// Event payload structure matching the user's specification
type EventPayload struct {
	EventID     string    `json:"event_id" binding:"required"`
	EventType   string    `json:"event_type" binding:"required"`
	Timestamp   time.Time `json:"timestamp" binding:"required"`
	SessionID   string    `json:"session_id" binding:"required"`
	QuizID      string    `json:"quiz_id" binding:"required"`
	ClassroomID string    `json:"classroom_id" binding:"required"`
	QuestionID  string    `json:"question_id" binding:"required"`
	TeacherID   *string   `json:"teacher_id,omitempty"`
	TimerSec    *int      `json:"timer_sec,omitempty"`
	// Additional fields for ANSWER_SUBMITTED events
	StudentID      *string `json:"student_id,omitempty"`
	Answer         *string `json:"answer,omitempty"`
	ResponseTimeMs *int    `json:"response_time_ms,omitempty"`
}

func (h *Handler) CreateEvent(c *gin.Context) {
	var event EventPayload
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	err := h.service.ProcessEvent(event, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Event processed successfully",
		"event_id":  event.EventID,
		"timestamp": event.Timestamp,
	})
}

func (h *Handler) CreateBatchEvents(c *gin.Context) {
	var events []EventPayload
	if err := c.ShouldBindJSON(&events); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	processedCount := 0
	errors := []string{}

	for _, event := range events {
		err := h.service.ProcessEvent(event, userID)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Event %s: %s", event.EventID, err.Error()))
		} else {
			processedCount++
		}
	}

	response := gin.H{
		"message":         "Batch events processed",
		"user_id":         userID,
		"processed_count": processedCount,
		"total_events":    len(events),
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	c.JSON(http.StatusCreated, response)
}
