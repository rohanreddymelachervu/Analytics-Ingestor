package repository

import (
	"time"

	"github.com/google/uuid"
)

// Report data structures
type ParticipantMetrics struct {
	StudentID uuid.UUID `json:"student_id"`
	Accuracy  float64   `json:"accuracy"`
}

type QuestionsPerMinuteStats struct {
	TotalQuestions int     `json:"total_questions"`
	AverageQPM     float64 `json:"average_qpm"`
	PeakQPM        float64 `json:"peak_qpm"`
}

type StudentPerformanceData struct {
	QuestionsAttempted  int     `json:"questions_attempted"`
	CorrectAnswers      int     `json:"correct_answers"`
	OverallAccuracy     float64 `json:"overall_accuracy"`
	AverageResponseTime string  `json:"average_response_time"`
}

type ClassroomEngagementData struct {
	TotalStudents   int     `json:"total_students"`
	ActiveStudents  int     `json:"active_students"`
	EngagementRate  float64 `json:"engagement_rate"`
	AverageAccuracy float64 `json:"average_accuracy"`
	TotalQuestions  int     `json:"total_questions"`
	ResponseRate    float64 `json:"response_rate"`
}

type ContentEffectivenessData struct {
	QuizID             uuid.UUID `json:"quiz_id"`
	TotalQuestions     int       `json:"total_questions"`
	AverageAccuracy    float64   `json:"average_accuracy"`
	OverallEngagement  float64   `json:"overall_engagement"`
	EffectivenessScore float64   `json:"effectiveness_score"`
	Recommendations    string    `json:"recommendations"`
}

type QuestionAnalysis struct {
	QuestionID         uuid.UUID `json:"question_id"`
	Accuracy           float64   `json:"accuracy"`
	AvgResponseTime    string    `json:"avg_response_time"`
	Difficulty         string    `json:"difficulty"`
	EffectivenessScore float64   `json:"effectiveness_score"`
}

// Critical data structures for missing metrics
type ResponseRateData struct {
	QuestionID       uuid.UUID `json:"question_id"`
	SessionID        uuid.UUID `json:"session_id"`
	StudentsReceived int       `json:"students_received"`
	StudentsAnswered int       `json:"students_answered"`
	ResponseRate     float64   `json:"response_rate"`
}

type LatencyData struct {
	QuestionID         uuid.UUID     `json:"question_id"`
	SessionID          uuid.UUID     `json:"session_id"`
	FirstAnswerLatency time.Duration `json:"first_answer_latency"`
	AverageLatency     time.Duration `json:"average_latency"`
	MedianLatency      time.Duration `json:"median_latency"`
}

type TimeoutData struct {
	QuestionID    uuid.UUID `json:"question_id"`
	SessionID     uuid.UUID `json:"session_id"`
	TotalStudents int       `json:"total_students"`
	TimeoutCount  int       `json:"timeout_count"`
	SkippedCount  int       `json:"skipped_count"`
	TimeoutRate   float64   `json:"timeout_rate"`
	SkippedRate   float64   `json:"skipped_rate"`
}

type CompletionRateData struct {
	SessionID         uuid.UUID `json:"session_id"`
	TotalStudents     int       `json:"total_students"`
	CompletedStudents int       `json:"completed_students"`
	CompletionRate    float64   `json:"completion_rate"`
	TotalQuestions    int       `json:"total_questions"`
	AverageCompletion float64   `json:"average_completion"`
}

type DropoffPoint struct {
	QuestionID      uuid.UUID `json:"question_id"`
	QuestionOrder   int       `json:"question_order"`
	StudentsAtStart int       `json:"students_at_start"`
	StudentsAtEnd   int       `json:"students_at_end"`
	DropoffRate     float64   `json:"dropoff_rate"`
}
