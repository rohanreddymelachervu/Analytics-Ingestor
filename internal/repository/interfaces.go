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

	// Analytics methods with pagination support
	GetActiveParticipants(sessionID uuid.UUID, timeRange time.Duration, pagination PaginationParams) (*PaginatedResponse[ParticipantMetrics], error)
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

	// Paginated methods for large result sets
	GetStudentPerformanceList(classroomID uuid.UUID, pagination PaginationParams) (*PaginatedResponse[StudentPerformanceData], error)
	GetClassroomEngagementHistory(classroomID uuid.UUID, dateRange time.Duration, pagination PaginationParams) (*PaginatedResponse[ClassroomEngagementData], error)

	// NEW: Missing basic metrics methods
	// Quiz-level analytics across all sessions and classrooms
	GetQuizSummary(quizID uuid.UUID) (*QuizSummaryData, error)

	// Question-level analytics across all sessions
	GetQuestionAnalysis(questionID uuid.UUID) (*QuestionAnalysisData, error)
	GetQuizQuestionsList(quizID uuid.UUID, pagination PaginationParams) (*PaginatedResponse[QuestionAnalysisData], error)

	// Session comparison and history
	GetClassroomSessions(classroomID uuid.UUID, pagination PaginationParams) (*PaginatedResponse[SessionComparisonData], error)
	GetQuizSessions(quizID uuid.UUID, pagination PaginationParams) (*PaginatedResponse[SessionComparisonData], error)

	// Student rankings and leaderboards
	GetClassroomStudentRankings(classroomID uuid.UUID, pagination PaginationParams) (*PaginatedResponse[StudentRankingData], error)
	GetSessionStudentRankings(sessionID uuid.UUID, pagination PaginationParams) (*PaginatedResponse[StudentRankingData], error)

	// NEW: Basic Overview Methods
	// Classroom overview dashboard - basic stats
	GetClassroomOverview(classroomID uuid.UUID) (*ClassroomOverviewData, error)
	// Class performance summary - overall averages and participation
	GetClassPerformanceSummary(classroomID uuid.UUID) (*ClassPerformanceSummaryData, error)
	// Student activity summary - participation and quiz history
	GetStudentActivitySummary(studentID, classroomID uuid.UUID) (*StudentActivitySummaryData, error)
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

	// Paginated classroom methods
	GetClassroomStudentsPaginated(classroomID uuid.UUID, pagination PaginationParams) (*PaginatedResponse[models.Student], error)
}
