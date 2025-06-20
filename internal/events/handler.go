package events

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rohanreddymelachervu/ingestor/internal/kafka"
	"github.com/rohanreddymelachervu/ingestor/internal/models"
)

type Handler struct {
	service   *Service
	kafkaMode bool
	producer  *kafka.Producer
}

// NewHandler creates a handler with direct database processing (default mode)
func NewHandler(s *Service) *Handler {
	return &Handler{
		service:   s,
		kafkaMode: false,
	}
}

// NewHandlerWithKafka creates a handler with Kafka producer for event publishing
func NewHandlerWithKafka(s *Service, producer *kafka.Producer) *Handler {
	return &Handler{
		service:   s,
		kafkaMode: true,
		producer:  producer,
	}
}

func (h *Handler) CreateEvent(c *gin.Context) {
	var event models.EventPayload
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	if h.kafkaMode && h.producer != nil {
		// Kafka mode: publish to topic
		err := h.producer.PublishEvent(event.EventID, event.EventType, event.SessionID, event)
		if err != nil {
			log.Printf("Failed to publish to Kafka, falling back to direct processing: %v", err)
			// Fallback to direct processing if Kafka fails
			err = h.service.ProcessEvent(event, userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		c.JSON(http.StatusCreated, gin.H{
			"message":   "Event queued successfully",
			"event_id":  event.EventID,
			"timestamp": event.Timestamp,
			"mode":      "kafka",
		})
	} else {
		// Direct mode: process immediately
		err := h.service.ProcessEvent(event, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message":   "Event processed successfully",
			"event_id":  event.EventID,
			"timestamp": event.Timestamp,
			"mode":      "direct",
		})
	}
}

func (h *Handler) CreateBatchEvents(c *gin.Context) {
	var events []models.EventPayload
	if err := c.ShouldBindJSON(&events); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	processedCount := 0
	errors := []string{}

	for _, event := range events {
		if h.kafkaMode && h.producer != nil {
			// Kafka mode: publish each event
			err := h.producer.PublishEvent(event.EventID, event.EventType, event.SessionID, event)
			if err != nil {
				log.Printf("Failed to publish event %s to Kafka: %v", event.EventID, err)
				// Try direct processing as fallback
				err = h.service.ProcessEvent(event, userID)
				if err != nil {
					errors = append(errors, fmt.Sprintf("Event %s: %s", event.EventID, err.Error()))
				} else {
					processedCount++
				}
			} else {
				processedCount++
			}
		} else {
			// Direct mode: process immediately
			err := h.service.ProcessEvent(event, userID)
			if err != nil {
				errors = append(errors, fmt.Sprintf("Event %s: %s", event.EventID, err.Error()))
			} else {
				processedCount++
			}
		}
	}

	mode := "direct"
	if h.kafkaMode {
		mode = "kafka"
	}

	response := gin.H{
		"message":         "Batch events processed",
		"user_id":         userID,
		"processed_count": processedCount,
		"total_events":    len(events),
		"mode":            mode,
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	c.JSON(http.StatusCreated, response)
}
