package models

import "time"

// EventPayload represents an event in the analytics system
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
