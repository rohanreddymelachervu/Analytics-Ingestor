package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/rohanreddymelachervu/ingestor/internal/models"
)

// EventRepository handles event-related database operations
type EventRepository interface {
	SaveQuestionPublishedEvent(event *models.QuestionPublishedEvent) error
	SaveAnswerSubmittedEvent(event *models.AnswerSubmittedEvent) error
	GetActiveParticipants(sessionID uuid.UUID, timeRange time.Duration) ([]ParticipantMetrics, error)
	GetQuestionsPerMinuteStats(sessionID uuid.UUID) (*QuestionsPerMinuteStats, error)
	GetStudentPerformance(studentID, classroomID uuid.UUID) (*StudentPerformanceData, error)
	GetClassroomEngagement(classroomID uuid.UUID, dateRange time.Duration) (*ClassroomEngagementData, error)
	GetContentEffectiveness(quizID uuid.UUID) (*ContentEffectivenessData, error)

	// Critical methods for missing metrics
	ValidateAnswerTiming(sessionID, questionID uuid.UUID, answerTimestamp time.Time) error
	GetResponseRate(sessionID, questionID uuid.UUID) (*ResponseRateData, error)
	GetLatencyToFirstAnswer(sessionID, questionID uuid.UUID) (*LatencyData, error)
	GetTimeoutAndSkippedRate(sessionID, questionID uuid.UUID) (*TimeoutData, error)
	GetCompletionRate(sessionID uuid.UUID) (*CompletionRateData, error)
	GetDropoffPoints(sessionID uuid.UUID) ([]DropoffPoint, error)
}

// QuizRepository handles quiz-related operations
type QuizRepository interface {
	CreateQuiz(quiz *models.Quiz) error
	GetQuizByID(quizID uuid.UUID) (*models.Quiz, error)
	CreateQuestion(question *models.Question) error
	GetQuestionByID(questionID uuid.UUID) (*models.Question, error)
}

// SessionRepository handles session-related operations
type SessionRepository interface {
	CreateSession(session *models.QuizSession) error
	GetSessionByID(sessionID uuid.UUID) (*models.QuizSession, error)
	UpdateSession(session *models.QuizSession) error
}

// ClassroomRepository handles classroom and student operations
type ClassroomRepository interface {
	CreateClassroom(classroom *models.Classroom) error
	CreateStudent(student *models.Student) error
	AddStudentToClassroom(classroomID, studentID uuid.UUID) error
	GetClassroomStudents(classroomID uuid.UUID) ([]models.Student, error)
}
