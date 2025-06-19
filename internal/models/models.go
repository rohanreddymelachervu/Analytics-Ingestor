package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Quiz represents a quiz entity - matches 000001_init_schema.up.sql
type Quiz struct {
	QuizID      uuid.UUID `gorm:"type:uuid;primary_key" json:"quiz_id"`
	Title       string    `gorm:"not null" json:"title"`
	Description *string   `json:"description"`
}

func (Quiz) TableName() string {
	return "quizzes"
}

// Classroom represents a classroom entity - matches 000002_init_schema.up.sql
type Classroom struct {
	ClassroomID uuid.UUID `gorm:"type:uuid;primary_key" json:"classroom_id"`
	Name        string    `gorm:"not null" json:"name"`
}

func (Classroom) TableName() string {
	return "classrooms"
}

// Student represents a student entity - matches 000003_init_schema.up.sql
type Student struct {
	StudentID uuid.UUID `gorm:"type:uuid;primary_key" json:"student_id"`
	Name      *string   `json:"name"` // nullable in migration
}

func (Student) TableName() string {
	return "students"
}

// Question represents a question entity - matches 000006_init_schema.up.sql
type Question struct {
	QuestionID uuid.UUID `gorm:"type:uuid;primary_key" json:"question_id"`
	QuizID     uuid.UUID `gorm:"type:uuid;not null" json:"quiz_id"`
}

func (Question) TableName() string {
	return "questions"
}

// QuizSession represents an active quiz session - matches 000005_init_schema.up.sql
type QuizSession struct {
	SessionID   uuid.UUID  `gorm:"type:uuid;primary_key" json:"session_id"`
	QuizID      uuid.UUID  `gorm:"type:uuid;not null" json:"quiz_id"`
	ClassroomID uuid.UUID  `gorm:"type:uuid;not null" json:"classroom_id"`
	StartedAt   time.Time  `gorm:"not null" json:"started_at"`
	EndedAt     *time.Time `json:"ended_at"`
}

func (QuizSession) TableName() string {
	return "quiz_sessions"
}

// ClassroomStudent represents the many-to-many relationship - matches 000004_init_schema.up.sql
type ClassroomStudent struct {
	ClassroomID uuid.UUID `gorm:"type:uuid;primary_key" json:"classroom_id"`
	StudentID   uuid.UUID `gorm:"type:uuid;primary_key" json:"student_id"`
}

func (ClassroomStudent) TableName() string {
	return "classroom_students"
}

// QuestionPublishedEvent represents teacher publishing a question - matches 000007_init_schema.up.sql
type QuestionPublishedEvent struct {
	EventID          uuid.UUID  `gorm:"type:uuid;primary_key" json:"event_id"`
	SessionID        uuid.UUID  `gorm:"type:uuid;not null" json:"session_id"`
	QuestionID       uuid.UUID  `gorm:"type:uuid;not null" json:"question_id"`
	TeacherID        *uuid.UUID `gorm:"type:uuid" json:"teacher_id"` // nullable
	PublishedAt      time.Time  `gorm:"not null" json:"published_at"`
	TimerDurationSec int        `gorm:"not null" json:"timer_duration_sec"`
}

func (QuestionPublishedEvent) TableName() string {
	return "question_published_events"
}

// AnswerSubmittedEvent represents student submitting an answer - matches 000008_init_schema.up.sql
type AnswerSubmittedEvent struct {
	EventID     uuid.UUID `gorm:"type:uuid;primary_key" json:"event_id"`
	SessionID   uuid.UUID `gorm:"type:uuid;not null" json:"session_id"`
	QuestionID  uuid.UUID `gorm:"type:uuid;not null" json:"question_id"`
	StudentID   uuid.UUID `gorm:"type:uuid;not null" json:"student_id"`
	Answer      string    `gorm:"not null" json:"answer"`
	IsCorrect   bool      `gorm:"not null" json:"is_correct"`
	SubmittedAt time.Time `gorm:"not null" json:"submitted_at"`
}

func (AnswerSubmittedEvent) TableName() string {
	return "answer_submitted_events"
}

// User represents authentication users - matches 000009_init_schema.up.sql
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Email     string    `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	Role      string    `gorm:"size:20;not null;default:'writer'" json:"role"`
	CreatedAt time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:now()" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

// AutoMigrate runs all migrations for the models
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Quiz{},
		&Classroom{},
		&Student{},
		&Question{},
		&QuizSession{},
		&ClassroomStudent{},
		&QuestionPublishedEvent{},
		&AnswerSubmittedEvent{},
		&User{},
	)
}
