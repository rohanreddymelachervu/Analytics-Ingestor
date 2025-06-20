package repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rohanreddymelachervu/ingestor/internal/models"
	"gorm.io/gorm"
)

// Repository implementations
type eventRepository struct {
	db *gorm.DB
}

type quizRepository struct {
	db *gorm.DB
}

type sessionRepository struct {
	db *gorm.DB
}

type classroomRepository struct {
	db *gorm.DB
}

// Constructor functions
func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}

func NewQuizRepository(db *gorm.DB) QuizRepository {
	return &quizRepository{db: db}
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func NewClassroomRepository(db *gorm.DB) ClassroomRepository {
	return &classroomRepository{db: db}
}

// EventRepository implementations
func (r *eventRepository) SaveQuestionPublishedEvent(event *models.QuestionPublishedEvent) error {
	return r.db.Create(event).Error
}

func (r *eventRepository) SaveAnswerSubmittedEvent(event *models.AnswerSubmittedEvent) error {
	return r.db.Create(event).Error
}

func (r *eventRepository) GetActiveParticipants(sessionID uuid.UUID, timeRange time.Duration, pagination PaginationParams) (*PaginatedResponse[ParticipantMetrics], error) {
	var results []ParticipantMetrics
	var totalCount int64

	cutoffTime := time.Now().Add(-timeRange)

	// First, get the total count for pagination
	err := r.db.Raw(`
		SELECT COUNT(DISTINCT s.student_id)
		FROM answer_submitted_events ase
		JOIN students s ON ase.student_id = s.student_id
		WHERE ase.session_id = ? AND ase.submitted_at >= ?
	`, sessionID, cutoffTime).Scan(&totalCount).Error

	if err != nil {
		return nil, err
	}

	// Then get the paginated results
	err = r.db.Raw(`
		SELECT 
			s.student_id,
			s.name,
			MAX(ase.submitted_at) as last_activity,
			COUNT(ase.event_id) as answers_submitted,
			SUM(CASE WHEN ase.is_correct THEN 1 ELSE 0 END) as correct_answers,
			ROUND(
				AVG(CASE WHEN ase.is_correct THEN 1.0 ELSE 0.0 END) * 100, 2
			) as accuracy
		FROM answer_submitted_events ase
		JOIN students s ON ase.student_id = s.student_id
		WHERE ase.session_id = ? AND ase.submitted_at >= ?
		GROUP BY s.student_id, s.name
		ORDER BY last_activity DESC
		LIMIT ? OFFSET ?
	`, sessionID, cutoffTime, pagination.PageSize, pagination.Offset).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	response := NewPaginatedResponse(results, pagination, int(totalCount))
	return &response, nil
}

func (r *eventRepository) GetQuestionsPerMinuteStats(sessionID uuid.UUID) (*QuestionsPerMinuteStats, error) {
	var stats QuestionsPerMinuteStats

	err := r.db.Raw(`
		SELECT 
			COUNT(*) as total_questions,
			ROUND(
				COUNT(*) / GREATEST(EXTRACT(EPOCH FROM (MAX(published_at) - MIN(published_at))) / 60, 1), 2
			) as average_qpm
		FROM question_published_events 
		WHERE session_id = ?
	`, sessionID).Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	stats.PeakQPM = stats.AverageQPM
	return &stats, nil
}

func (r *eventRepository) GetStudentPerformance(studentID, classroomID uuid.UUID) (*StudentPerformanceData, error) {
	var performance StudentPerformanceData

	err := r.db.Raw(`
		SELECT 
			COUNT(*) as questions_attempted,
			SUM(CASE WHEN is_correct THEN 1 ELSE 0 END) as correct_answers,
			ROUND(AVG(CASE WHEN is_correct THEN 1.0 ELSE 0.0 END) * 100, 2) as overall_accuracy,
			'N/A' as average_response_time
		FROM answer_submitted_events ase
		JOIN quiz_sessions qs ON ase.session_id = qs.session_id
		WHERE ase.student_id = ? AND qs.classroom_id = ?
	`, studentID, classroomID).Scan(&performance).Error

	return &performance, err
}

func (r *eventRepository) GetClassroomEngagement(classroomID uuid.UUID, dateRange time.Duration) (*ClassroomEngagementData, error) {
	var engagement ClassroomEngagementData
	cutoffTime := time.Now().Add(-dateRange)

	// Get total students in classroom
	err := r.db.Raw(`
		SELECT COUNT(*) as total_students
		FROM classroom_students 
		WHERE classroom_id = ?
	`, classroomID).Scan(&engagement.TotalStudents).Error

	if err != nil {
		return nil, err
	}

	// Get engagement metrics
	err = r.db.Raw(`
		SELECT 
			COUNT(DISTINCT ase.student_id) as active_students,
			COUNT(DISTINCT qpe.question_id) as total_questions,
			ROUND(
				COUNT(DISTINCT ase.student_id) * 100.0 / ?, 2
			) as engagement_rate,
			ROUND(
				AVG(CASE WHEN ase.is_correct THEN 1.0 ELSE 0.0 END) * 100, 2
			) as average_accuracy,
			ROUND(
				COUNT(ase.event_id) * 100.0 / 
				(COUNT(DISTINCT ase.student_id) * COUNT(DISTINCT qpe.question_id)), 2
			) as response_rate
		FROM quiz_sessions qs
		LEFT JOIN question_published_events qpe ON qpe.session_id = qs.session_id
		LEFT JOIN answer_submitted_events ase ON ase.session_id = qs.session_id 
			AND ase.submitted_at >= ?
		WHERE qs.classroom_id = ? AND qs.started_at >= ?
	`, engagement.TotalStudents, cutoffTime, classroomID, cutoffTime).Scan(&engagement).Error

	return &engagement, err
}

func (r *eventRepository) GetContentEffectiveness(quizID uuid.UUID) (*ContentEffectivenessData, error) {
	var effectiveness ContentEffectivenessData

	err := r.db.Raw(`
		SELECT 
			? as quiz_id,
			COUNT(DISTINCT qpe.question_id) as total_questions,
			ROUND(AVG(CASE WHEN ase.is_correct THEN 1.0 ELSE 0.0 END) * 100, 2) as average_accuracy,
			ROUND(
				COUNT(DISTINCT ase.student_id) * 100.0 / 
				GREATEST(COUNT(DISTINCT cs.student_id), 1), 2
			) as overall_engagement
		FROM quiz_sessions qs
		JOIN question_published_events qpe ON qpe.session_id = qs.session_id
		LEFT JOIN answer_submitted_events ase ON ase.session_id = qs.session_id
		LEFT JOIN classroom_students cs ON cs.classroom_id = qs.classroom_id
		WHERE qs.quiz_id = ?
	`, quizID, quizID).Scan(&effectiveness).Error

	if err != nil {
		return nil, err
	}

	effectiveness.QuizID = quizID
	effectiveness.EffectivenessScore = (effectiveness.AverageAccuracy + effectiveness.OverallEngagement) / 2

	if effectiveness.EffectivenessScore >= 80 {
		effectiveness.Recommendations = "Quiz is performing well. Consider similar content."
	} else if effectiveness.EffectivenessScore >= 60 {
		effectiveness.Recommendations = "Quiz has moderate effectiveness. Review difficult questions."
	} else {
		effectiveness.Recommendations = "Quiz needs improvement. Consider revising content and delivery."
	}

	return &effectiveness, nil
}

func (r *eventRepository) ValidateAnswerTiming(sessionID, questionID uuid.UUID, answerTimestamp time.Time) error {
	var result struct {
		PublishedAt      time.Time `db:"published_at"`
		TimerDurationSec int       `db:"timer_duration_sec"`
	}

	err := r.db.Raw(`
		SELECT published_at, timer_duration_sec
		FROM question_published_events 
		WHERE session_id = ? AND question_id = ?
		ORDER BY published_at DESC 
		LIMIT 1
	`, sessionID, questionID).Scan(&result).Error

	if err != nil {
		return fmt.Errorf("failed to find question publish event: %v", err)
	}

	deadline := result.PublishedAt.Add(time.Duration(result.TimerDurationSec) * time.Second)
	if answerTimestamp.After(deadline) {
		return fmt.Errorf("answer submitted at %v is after deadline %v", answerTimestamp, deadline)
	}

	return nil
}

func (r *eventRepository) GetResponseRate(sessionID, questionID uuid.UUID) (*ResponseRateData, error) {
	var data ResponseRateData

	err := r.db.Raw(`
		SELECT 
			? as question_id,
			? as session_id,
			COUNT(DISTINCT cs.student_id) as students_received,
			COUNT(DISTINCT ase.student_id) as students_answered,
			ROUND(
				COUNT(DISTINCT ase.student_id) * 100.0 / 
				GREATEST(COUNT(DISTINCT cs.student_id), 1), 2
			) as response_rate
		FROM quiz_sessions qs
		JOIN classroom_students cs ON cs.classroom_id = qs.classroom_id
		LEFT JOIN answer_submitted_events ase ON ase.session_id = qs.session_id 
			AND ase.question_id = ?
		WHERE qs.session_id = ?
	`, questionID, sessionID, questionID, sessionID).Scan(&data).Error

	return &data, err
}

func (r *eventRepository) GetLatencyToFirstAnswer(sessionID, questionID uuid.UUID) (*LatencyData, error) {
	var data LatencyData

	// Get question publish time and answer times
	var result struct {
		PublishedAt     time.Time  `db:"published_at"`
		FirstAnswerAt   *time.Time `db:"first_answer_at"`
		AvgResponseTime *int       `db:"avg_response_time"`
		MedianTime      *int       `db:"median_time"`
	}

	err := r.db.Raw(`
		WITH answer_latencies AS (
			SELECT 
				EXTRACT(EPOCH FROM (ase.submitted_at - qpe.published_at)) as latency_seconds
			FROM question_published_events qpe
			JOIN answer_submitted_events ase ON ase.session_id = qpe.session_id 
				AND ase.question_id = qpe.question_id
			WHERE qpe.session_id = ? AND qpe.question_id = ?
		)
		SELECT 
			qpe.published_at,
			MIN(ase.submitted_at) as first_answer_at,
			ROUND(AVG(EXTRACT(EPOCH FROM (ase.submitted_at - qpe.published_at)))) as avg_response_time,
			ROUND(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY EXTRACT(EPOCH FROM (ase.submitted_at - qpe.published_at)))) as median_time
		FROM question_published_events qpe
		LEFT JOIN answer_submitted_events ase ON ase.session_id = qpe.session_id 
			AND ase.question_id = qpe.question_id
		WHERE qpe.session_id = ? AND qpe.question_id = ?
		GROUP BY qpe.published_at
	`, sessionID, questionID, sessionID, questionID).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	data.QuestionID = questionID
	data.SessionID = sessionID

	if result.FirstAnswerAt != nil {
		data.FirstAnswerLatency = result.FirstAnswerAt.Sub(result.PublishedAt)
	}

	if result.AvgResponseTime != nil {
		data.AverageLatency = time.Duration(*result.AvgResponseTime) * time.Second
	}

	if result.MedianTime != nil {
		data.MedianLatency = time.Duration(*result.MedianTime) * time.Second
	}

	return &data, nil
}

func (r *eventRepository) GetTimeoutAndSkippedRate(sessionID, questionID uuid.UUID) (*TimeoutData, error) {
	var data TimeoutData

	err := r.db.Raw(`
		SELECT 
			? as question_id,
			? as session_id,
			COUNT(DISTINCT cs.student_id) as total_students,
			COUNT(DISTINCT cs.student_id) - COUNT(DISTINCT ase.student_id) as timeout_count,
			0 as skipped_count,
			ROUND(
				(COUNT(DISTINCT cs.student_id) - COUNT(DISTINCT ase.student_id)) * 100.0 / 
				GREATEST(COUNT(DISTINCT cs.student_id), 1), 2
			) as timeout_rate,
			0.0 as skipped_rate
		FROM quiz_sessions qs
		JOIN classroom_students cs ON cs.classroom_id = qs.classroom_id
		LEFT JOIN answer_submitted_events ase ON ase.session_id = qs.session_id 
			AND ase.question_id = ?
		WHERE qs.session_id = ?
	`, questionID, sessionID, questionID, sessionID).Scan(&data).Error

	return &data, err
}

func (r *eventRepository) GetCompletionRate(sessionID uuid.UUID) (*CompletionRateData, error) {
	var data CompletionRateData

	err := r.db.Raw(`
		WITH session_stats AS (
			SELECT 
				COUNT(DISTINCT cs.student_id) as total_students,
				COUNT(DISTINCT qpe.question_id) as total_questions
			FROM quiz_sessions qs
			JOIN classroom_students cs ON cs.classroom_id = qs.classroom_id
			LEFT JOIN question_published_events qpe ON qpe.session_id = qs.session_id
			WHERE qs.session_id = ?
		),
		student_completion AS (
			SELECT 
				ase.student_id,
				COUNT(DISTINCT ase.question_id) as questions_answered
			FROM answer_submitted_events ase
			WHERE ase.session_id = ?
			GROUP BY ase.student_id
		)
		SELECT 
			? as session_id,
			ss.total_students,
			COUNT(sc.student_id) as completed_students,
			ROUND(
				COUNT(sc.student_id) * 100.0 / GREATEST(ss.total_students, 1), 2
			) as completion_rate,
			ss.total_questions,
			ROUND(
				AVG(sc.questions_answered * 100.0 / GREATEST(ss.total_questions, 1)), 2
			) as average_completion
		FROM session_stats ss
		LEFT JOIN student_completion sc ON sc.questions_answered = ss.total_questions
		GROUP BY ss.total_students, ss.total_questions
	`, sessionID, sessionID, sessionID).Scan(&data).Error

	return &data, err
}

func (r *eventRepository) GetDropoffPoints(sessionID uuid.UUID) ([]DropoffPoint, error) {
	var dropoffs []DropoffPoint

	err := r.db.Raw(`
		WITH question_order AS (
			SELECT 
				question_id,
				ROW_NUMBER() OVER (ORDER BY published_at) as question_order
			FROM question_published_events
			WHERE session_id = ?
		),
		student_progress AS (
			SELECT 
				qo.question_id,
				qo.question_order,
				COUNT(DISTINCT ase.student_id) as students_answered
			FROM question_order qo
			LEFT JOIN answer_submitted_events ase ON ase.question_id = qo.question_id 
				AND ase.session_id = ?
			GROUP BY qo.question_id, qo.question_order
		),
		total_students AS (
			SELECT COUNT(DISTINCT cs.student_id) as total_count
			FROM quiz_sessions qs
			JOIN classroom_students cs ON cs.classroom_id = qs.classroom_id
			WHERE qs.session_id = ?
		)
		SELECT 
			sp.question_id,
			sp.question_order,
			ts.total_count as students_at_start,
			sp.students_answered as students_at_end,
			ROUND(
				(ts.total_count - sp.students_answered) * 100.0 / 
				GREATEST(ts.total_count, 1), 2
			) as dropoff_rate
		FROM student_progress sp
		CROSS JOIN total_students ts
		ORDER BY sp.question_order
	`, sessionID, sessionID, sessionID).Scan(&dropoffs).Error

	return dropoffs, err
}

// QuizRepository implementations
func (r *quizRepository) CreateQuiz(quiz *models.Quiz) error {
	return r.db.Create(quiz).Error
}

func (r *quizRepository) GetQuizByID(quizID uuid.UUID) (*models.Quiz, error) {
	var quiz models.Quiz
	err := r.db.Where("quiz_id = ?", quizID).First(&quiz).Error
	return &quiz, err
}

func (r *quizRepository) CreateQuestion(question *models.Question) error {
	return r.db.Create(question).Error
}

func (r *quizRepository) GetQuestionByID(questionID uuid.UUID) (*models.Question, error) {
	var question models.Question
	err := r.db.Where("question_id = ?", questionID).First(&question).Error
	return &question, err
}

// SessionRepository implementations
func (r *sessionRepository) CreateSession(session *models.QuizSession) error {
	return r.db.Create(session).Error
}

func (r *sessionRepository) GetSessionByID(sessionID uuid.UUID) (*models.QuizSession, error) {
	var session models.QuizSession
	err := r.db.Where("session_id = ?", sessionID).First(&session).Error
	return &session, err
}

func (r *sessionRepository) UpdateSession(session *models.QuizSession) error {
	return r.db.Save(session).Error
}

// ClassroomRepository implementations
func (r *classroomRepository) CreateClassroom(classroom *models.Classroom) error {
	return r.db.Create(classroom).Error
}

func (r *classroomRepository) CreateStudent(student *models.Student) error {
	return r.db.Create(student).Error
}

func (r *classroomRepository) AddStudentToClassroom(classroomID, studentID uuid.UUID) error {
	classroomStudent := &models.ClassroomStudent{
		ClassroomID: classroomID,
		StudentID:   studentID,
	}
	return r.db.Create(classroomStudent).Error
}

func (r *classroomRepository) GetClassroomStudents(classroomID uuid.UUID) ([]models.Student, error) {
	var students []models.Student
	err := r.db.Raw(`
		SELECT s.* 
		FROM students s
		JOIN classroom_students cs ON s.student_id = cs.student_id
		WHERE cs.classroom_id = ?
		ORDER BY s.name
	`, classroomID).Scan(&students).Error
	return students, err
}

// Paginated version of GetClassroomStudents
func (r *classroomRepository) GetClassroomStudentsPaginated(classroomID uuid.UUID, pagination PaginationParams) (*PaginatedResponse[models.Student], error) {
	var students []models.Student
	var totalCount int64

	// Get total count
	err := r.db.Raw(`
		SELECT COUNT(*)
		FROM students s
		JOIN classroom_students cs ON s.student_id = cs.student_id
		WHERE cs.classroom_id = ?
	`, classroomID).Scan(&totalCount).Error

	if err != nil {
		return nil, err
	}

	// Get paginated results
	err = r.db.Raw(`
		SELECT s.* 
		FROM students s
		JOIN classroom_students cs ON s.student_id = cs.student_id
		WHERE cs.classroom_id = ?
		ORDER BY s.name
		LIMIT ? OFFSET ?
	`, classroomID, pagination.PageSize, pagination.Offset).Scan(&students).Error

	if err != nil {
		return nil, err
	}

	response := NewPaginatedResponse(students, pagination, int(totalCount))
	return &response, nil
}

// New paginated method for student performance across a classroom
func (r *eventRepository) GetStudentPerformanceList(classroomID uuid.UUID, pagination PaginationParams) (*PaginatedResponse[StudentPerformanceData], error) {
	var results []StudentPerformanceData
	var totalCount int64

	// Get total count of students in classroom
	err := r.db.Raw(`
		SELECT COUNT(DISTINCT s.student_id)
		FROM students s
		JOIN classroom_students cs ON s.student_id = cs.student_id
		WHERE cs.classroom_id = ?
	`, classroomID).Scan(&totalCount).Error

	if err != nil {
		return nil, err
	}

	// Get paginated student performance data
	err = r.db.Raw(`
		SELECT 
			s.student_id,
			COALESCE(COUNT(ase.event_id), 0) as questions_attempted,
			COALESCE(SUM(CASE WHEN ase.is_correct THEN 1 ELSE 0 END), 0) as correct_answers,
			COALESCE(ROUND(AVG(CASE WHEN ase.is_correct THEN 1.0 ELSE 0.0 END) * 100, 2), 0) as overall_accuracy,
			'N/A' as average_response_time
		FROM students s
		JOIN classroom_students cs ON s.student_id = cs.student_id
		LEFT JOIN answer_submitted_events ase ON ase.student_id = s.student_id
		LEFT JOIN quiz_sessions qs ON ase.session_id = qs.session_id AND qs.classroom_id = cs.classroom_id
		WHERE cs.classroom_id = ?
		GROUP BY s.student_id
		ORDER BY s.name
		LIMIT ? OFFSET ?
	`, classroomID, pagination.PageSize, pagination.Offset).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	response := NewPaginatedResponse(results, pagination, int(totalCount))
	return &response, nil
}

// New paginated method for classroom engagement history
func (r *eventRepository) GetClassroomEngagementHistory(classroomID uuid.UUID, dateRange time.Duration, pagination PaginationParams) (*PaginatedResponse[ClassroomEngagementData], error) {
	var results []ClassroomEngagementData
	var totalCount int64

	cutoffTime := time.Now().Add(-dateRange)

	// For this example, we'll return daily engagement data
	// In a real implementation, you might want different time granularities

	// Get total count of days in the range
	err := r.db.Raw(`
		SELECT COUNT(DISTINCT DATE(qs.started_at))
		FROM quiz_sessions qs
		WHERE qs.classroom_id = ? AND qs.started_at >= ?
	`, classroomID, cutoffTime).Scan(&totalCount).Error

	if err != nil {
		return nil, err
	}

	// Get paginated daily engagement data
	err = r.db.Raw(`
		WITH daily_sessions AS (
			SELECT 
				DATE(qs.started_at) as session_date,
				COUNT(DISTINCT qs.session_id) as daily_sessions,
				COUNT(DISTINCT ase.student_id) as active_students,
				COUNT(DISTINCT qpe.question_id) as total_questions,
				COUNT(ase.event_id) as total_answers,
				ROUND(AVG(CASE WHEN ase.is_correct THEN 1.0 ELSE 0.0 END) * 100, 2) as average_accuracy
			FROM quiz_sessions qs
			LEFT JOIN question_published_events qpe ON qpe.session_id = qs.session_id
			LEFT JOIN answer_submitted_events ase ON ase.session_id = qs.session_id
			WHERE qs.classroom_id = ? AND qs.started_at >= ?
			GROUP BY DATE(qs.started_at)
			ORDER BY session_date DESC
			LIMIT ? OFFSET ?
		),
		classroom_info AS (
			SELECT COUNT(*) as total_students
			FROM classroom_students 
			WHERE classroom_id = ?
		)
		SELECT 
			ci.total_students,
			ds.active_students,
			ROUND(ds.active_students * 100.0 / ci.total_students, 2) as engagement_rate,
			ds.average_accuracy,
			ds.total_questions,
			CASE 
				WHEN ds.total_questions > 0 AND ds.active_students > 0 
				THEN ROUND(ds.total_answers * 100.0 / (ds.active_students * ds.total_questions), 2)
				ELSE 0.0 
			END as response_rate
		FROM daily_sessions ds
		CROSS JOIN classroom_info ci
	`, classroomID, cutoffTime, pagination.PageSize, pagination.Offset, classroomID).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	response := NewPaginatedResponse(results, pagination, int(totalCount))
	return &response, nil
}

// NEW: Missing basic metrics implementations

// GetQuizSummary - aggregates quiz performance across all sessions and classrooms
func (r *eventRepository) GetQuizSummary(quizID uuid.UUID) (*QuizSummaryData, error) {
	var summary QuizSummaryData

	err := r.db.Raw(`
		WITH quiz_stats AS (
			SELECT 
				q.quiz_id,
				q.title,
				COUNT(DISTINCT qs.session_id) as total_sessions,
				COUNT(DISTINCT qs.classroom_id) as total_classrooms,
				COUNT(DISTINCT ase.student_id) as total_students,
				COUNT(DISTINCT qpe.question_id) as total_questions,
				MIN(qs.started_at) as first_used,
				MAX(qs.started_at) as last_used,
				ROUND(AVG(CASE WHEN ase.is_correct THEN 1.0 ELSE 0.0 END) * 100, 2) as average_accuracy
			FROM quizzes q
			LEFT JOIN quiz_sessions qs ON q.quiz_id = qs.quiz_id
			LEFT JOIN question_published_events qpe ON qs.session_id = qpe.session_id
			LEFT JOIN answer_submitted_events ase ON qs.session_id = ase.session_id
			WHERE q.quiz_id = ?
			GROUP BY q.quiz_id, q.title
		),
		completion_stats AS (
			SELECT 
				AVG(session_completion.completion_rate) as average_completion
			FROM (
				SELECT 
					qs.session_id,
					COUNT(DISTINCT ase.student_id) * 100.0 / 
					GREATEST(COUNT(DISTINCT cs.student_id), 1) as completion_rate
				FROM quiz_sessions qs
				JOIN classroom_students cs ON cs.classroom_id = qs.classroom_id
				LEFT JOIN answer_submitted_events ase ON ase.session_id = qs.session_id
				WHERE qs.quiz_id = ?
				GROUP BY qs.session_id
			) session_completion
		),
		engagement_stats AS (
			SELECT 
				COUNT(DISTINCT ase.student_id) * 100.0 / 
				GREATEST(COUNT(DISTINCT cs.student_id), 1) as overall_engagement
			FROM quiz_sessions qs
			JOIN classroom_students cs ON cs.classroom_id = qs.classroom_id
			LEFT JOIN answer_submitted_events ase ON ase.session_id = qs.session_id
			WHERE qs.quiz_id = ?
		)
		SELECT 
			qs.*,
			COALESCE(cs.average_completion, 0) as average_completion,
			COALESCE(es.overall_engagement, 0) as overall_engagement,
			ROUND((COALESCE(qs.average_accuracy, 0) + COALESCE(es.overall_engagement, 0)) / 2, 2) as effectiveness_score
		FROM quiz_stats qs
		CROSS JOIN completion_stats cs
		CROSS JOIN engagement_stats es
	`, quizID, quizID, quizID).Scan(&summary).Error

	return &summary, err
}

// GetQuestionAnalysis - performance of individual questions across all sessions
func (r *eventRepository) GetQuestionAnalysis(questionID uuid.UUID) (*QuestionAnalysisData, error) {
	var analysis QuestionAnalysisData

	// Get basic question stats
	err := r.db.Raw(`
		SELECT 
			q.question_id,
			q.quiz_id,
			COUNT(ase.event_id) as total_attempts,
			SUM(CASE WHEN ase.is_correct THEN 1 ELSE 0 END) as correct_attempts,
			ROUND(AVG(CASE WHEN ase.is_correct THEN 1.0 ELSE 0.0 END) * 100, 2) as accuracy_rate,
			ROUND(AVG(EXTRACT(EPOCH FROM (ase.submitted_at - qpe.published_at))), 2) as average_response_time,
			COUNT(DISTINCT qpe.session_id) as usage_count
		FROM questions q
		LEFT JOIN question_published_events qpe ON q.question_id = qpe.question_id
		LEFT JOIN answer_submitted_events ase ON q.question_id = ase.question_id
		WHERE q.question_id = ?
		GROUP BY q.question_id, q.quiz_id
	`, questionID).Scan(&analysis).Error

	if err != nil {
		return nil, err
	}

	// Set difficulty rating based on accuracy
	if analysis.AccuracyRate >= 80 {
		analysis.DifficultyRating = "Easy"
	} else if analysis.AccuracyRate >= 60 {
		analysis.DifficultyRating = "Medium"
	} else if analysis.AccuracyRate >= 40 {
		analysis.DifficultyRating = "Hard"
	} else {
		analysis.DifficultyRating = "Very Hard"
	}

	// Get answer distribution
	var distributions []AnswerDistribution
	err = r.db.Raw(`
		SELECT 
			answer,
			COUNT(*) as count,
			ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER(), 2) as percentage
		FROM answer_submitted_events
		WHERE question_id = ?
		GROUP BY answer
		ORDER BY count DESC
	`, questionID).Scan(&distributions).Error

	if err != nil {
		return &analysis, err
	}

	// Convert to map for easier access
	analysis.AnswerDistribution = make(map[string]int)
	for _, dist := range distributions {
		analysis.AnswerDistribution[dist.Answer] = dist.Count
	}

	return &analysis, nil
}

// GetQuizQuestionsList - paginated list of all questions in a quiz with their analytics
func (r *eventRepository) GetQuizQuestionsList(quizID uuid.UUID, pagination PaginationParams) (*PaginatedResponse[QuestionAnalysisData], error) {
	var results []QuestionAnalysisData
	var totalCount int64

	// Get total count
	err := r.db.Raw(`
		SELECT COUNT(*)
		FROM questions
		WHERE quiz_id = ?
	`, quizID).Scan(&totalCount).Error

	if err != nil {
		return nil, err
	}

	// Get paginated results
	err = r.db.Raw(`
		SELECT 
			q.question_id,
			q.quiz_id,
			COALESCE(COUNT(ase.event_id), 0) as total_attempts,
			COALESCE(SUM(CASE WHEN ase.is_correct THEN 1 ELSE 0 END), 0) as correct_attempts,
			COALESCE(ROUND(AVG(CASE WHEN ase.is_correct THEN 1.0 ELSE 0.0 END) * 100, 2), 0) as accuracy_rate,
			COALESCE(ROUND(AVG(EXTRACT(EPOCH FROM (ase.submitted_at - qpe.published_at))), 2), 0) as average_response_time,
			COALESCE(COUNT(DISTINCT qpe.session_id), 0) as usage_count
		FROM questions q
		LEFT JOIN question_published_events qpe ON q.question_id = qpe.question_id
		LEFT JOIN answer_submitted_events ase ON q.question_id = ase.question_id
		WHERE q.quiz_id = ?
		GROUP BY q.question_id, q.quiz_id
		ORDER BY q.question_id
		LIMIT ? OFFSET ?
	`, quizID, pagination.PageSize, pagination.Offset).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// Set difficulty ratings
	for i := range results {
		if results[i].AccuracyRate >= 80 {
			results[i].DifficultyRating = "Easy"
		} else if results[i].AccuracyRate >= 60 {
			results[i].DifficultyRating = "Medium"
		} else if results[i].AccuracyRate >= 40 {
			results[i].DifficultyRating = "Hard"
		} else {
			results[i].DifficultyRating = "Very Hard"
		}
	}

	response := NewPaginatedResponse(results, pagination, int(totalCount))
	return &response, nil
}

// GetClassroomSessions - paginated list of all sessions in a classroom
func (r *eventRepository) GetClassroomSessions(classroomID uuid.UUID, pagination PaginationParams) (*PaginatedResponse[SessionComparisonData], error) {
	var results []SessionComparisonData
	var totalCount int64

	// Get total count
	err := r.db.Raw(`
		SELECT COUNT(*)
		FROM quiz_sessions
		WHERE classroom_id = ?
	`, classroomID).Scan(&totalCount).Error

	if err != nil {
		return nil, err
	}

	// Get paginated results
	err = r.db.Raw(`
		SELECT 
			qs.session_id,
			q.title as quiz_title,
			c.name as classroom_name,
			qs.started_at,
			CASE 
				WHEN qs.ended_at IS NOT NULL 
				THEN EXTRACT(EPOCH FROM (qs.ended_at - qs.started_at)) || ' seconds'
				ELSE NULL 
			END as duration,
			COUNT(DISTINCT cs.student_id) as total_students,
			COUNT(DISTINCT ase.student_id) as participating_students,
			COUNT(DISTINCT qpe.question_id) as total_questions,
			ROUND(AVG(CASE WHEN ase.is_correct THEN 1.0 ELSE 0.0 END) * 100, 2) as average_accuracy,
			ROUND(
				COUNT(DISTINCT ase.student_id) * 100.0 / 
				GREATEST(COUNT(DISTINCT cs.student_id), 1), 2
			) as completion_rate,
			ROUND(
				(COALESCE(AVG(CASE WHEN ase.is_correct THEN 1.0 ELSE 0.0 END) * 100, 0) + 
				 COUNT(DISTINCT ase.student_id) * 100.0 / GREATEST(COUNT(DISTINCT cs.student_id), 1)) / 2, 2
			) as engagement_score
		FROM quiz_sessions qs
		JOIN quizzes q ON qs.quiz_id = q.quiz_id
		JOIN classrooms c ON qs.classroom_id = c.classroom_id
		JOIN classroom_students cs ON cs.classroom_id = qs.classroom_id
		LEFT JOIN question_published_events qpe ON qpe.session_id = qs.session_id
		LEFT JOIN answer_submitted_events ase ON ase.session_id = qs.session_id
		WHERE qs.classroom_id = ?
		GROUP BY qs.session_id, q.title, c.name, qs.started_at, qs.ended_at
		ORDER BY qs.started_at DESC
		LIMIT ? OFFSET ?
	`, classroomID, pagination.PageSize, pagination.Offset).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	response := NewPaginatedResponse(results, pagination, int(totalCount))
	return &response, nil
}

// GetQuizSessions - paginated list of all sessions for a specific quiz
func (r *eventRepository) GetQuizSessions(quizID uuid.UUID, pagination PaginationParams) (*PaginatedResponse[SessionComparisonData], error) {
	var results []SessionComparisonData
	var totalCount int64

	// Get total count
	err := r.db.Raw(`
		SELECT COUNT(*)
		FROM quiz_sessions
		WHERE quiz_id = ?
	`, quizID).Scan(&totalCount).Error

	if err != nil {
		return nil, err
	}

	// Get paginated results
	err = r.db.Raw(`
		SELECT 
			qs.session_id,
			q.title as quiz_title,
			c.name as classroom_name,
			qs.started_at,
			CASE 
				WHEN qs.ended_at IS NOT NULL 
				THEN EXTRACT(EPOCH FROM (qs.ended_at - qs.started_at)) || ' seconds'
				ELSE NULL 
			END as duration,
			COUNT(DISTINCT cs.student_id) as total_students,
			COUNT(DISTINCT ase.student_id) as participating_students,
			COUNT(DISTINCT qpe.question_id) as total_questions,
			ROUND(AVG(CASE WHEN ase.is_correct THEN 1.0 ELSE 0.0 END) * 100, 2) as average_accuracy,
			ROUND(
				COUNT(DISTINCT ase.student_id) * 100.0 / 
				GREATEST(COUNT(DISTINCT cs.student_id), 1), 2
			) as completion_rate,
			ROUND(
				(COALESCE(AVG(CASE WHEN ase.is_correct THEN 1.0 ELSE 0.0 END) * 100, 0) + 
				 COUNT(DISTINCT ase.student_id) * 100.0 / GREATEST(COUNT(DISTINCT cs.student_id), 1)) / 2, 2
			) as engagement_score
		FROM quiz_sessions qs
		JOIN quizzes q ON qs.quiz_id = q.quiz_id
		JOIN classrooms c ON qs.classroom_id = c.classroom_id
		JOIN classroom_students cs ON cs.classroom_id = qs.classroom_id
		LEFT JOIN question_published_events qpe ON qpe.session_id = qs.session_id
		LEFT JOIN answer_submitted_events ase ON ase.session_id = qs.session_id
		WHERE qs.quiz_id = ?
		GROUP BY qs.session_id, q.title, c.name, qs.started_at, qs.ended_at
		ORDER BY qs.started_at DESC
		LIMIT ? OFFSET ?
	`, quizID, pagination.PageSize, pagination.Offset).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	response := NewPaginatedResponse(results, pagination, int(totalCount))
	return &response, nil
}

// GetClassroomStudentRankings - ranked student performance within a classroom
func (r *eventRepository) GetClassroomStudentRankings(classroomID uuid.UUID, pagination PaginationParams) (*PaginatedResponse[StudentRankingData], error) {
	var results []StudentRankingData
	var totalCount int64

	// Get total count
	err := r.db.Raw(`
		SELECT COUNT(DISTINCT s.student_id)
		FROM students s
		JOIN classroom_students cs ON s.student_id = cs.student_id
		WHERE cs.classroom_id = ?
	`, classroomID).Scan(&totalCount).Error

	if err != nil {
		return nil, err
	}

	// Get paginated ranked results
	err = r.db.Raw(`
		WITH student_stats AS (
			SELECT 
				s.student_id,
				COALESCE(s.name, 'Unknown') as student_name,
				COUNT(ase.event_id) as questions_attempted,
				SUM(CASE WHEN ase.is_correct THEN 1 ELSE 0 END) as correct_answers,
				ROUND(AVG(CASE WHEN ase.is_correct THEN 1.0 ELSE 0.0 END) * 100, 2) as accuracy_rate,
				ROUND(AVG(EXTRACT(EPOCH FROM (ase.submitted_at - qpe.published_at))), 2) as average_response_time,
				COUNT(DISTINCT ase.session_id) as sessions_participated
			FROM students s
			JOIN classroom_students cs ON s.student_id = cs.student_id
			LEFT JOIN answer_submitted_events ase ON s.student_id = ase.student_id
			LEFT JOIN quiz_sessions qs ON ase.session_id = qs.session_id AND qs.classroom_id = cs.classroom_id
			LEFT JOIN question_published_events qpe ON ase.question_id = qpe.question_id AND ase.session_id = qpe.session_id
			WHERE cs.classroom_id = ?
			GROUP BY s.student_id, s.name
		),
		ranked_students AS (
			SELECT 
				*,
				ROW_NUMBER() OVER (ORDER BY accuracy_rate DESC, questions_attempted DESC) as rank,
				ROUND(
					(ROW_NUMBER() OVER (ORDER BY accuracy_rate DESC, questions_attempted DESC) - 1) * 100.0 / 
					GREATEST(COUNT(*) OVER() - 1, 1), 2
				) as percentile
			FROM student_stats
		)
		SELECT *
		FROM ranked_students
		ORDER BY rank
		LIMIT ? OFFSET ?
	`, classroomID, pagination.PageSize, pagination.Offset).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	response := NewPaginatedResponse(results, pagination, int(totalCount))
	return &response, nil
}

// GetSessionStudentRankings - ranked student performance within a specific session
func (r *eventRepository) GetSessionStudentRankings(sessionID uuid.UUID, pagination PaginationParams) (*PaginatedResponse[StudentRankingData], error) {
	var results []StudentRankingData
	var totalCount int64

	// Get total count
	err := r.db.Raw(`
		SELECT COUNT(DISTINCT ase.student_id)
		FROM answer_submitted_events ase
		WHERE ase.session_id = ?
	`, sessionID).Scan(&totalCount).Error

	if err != nil {
		return nil, err
	}

	// Get paginated ranked results
	err = r.db.Raw(`
		WITH student_stats AS (
			SELECT 
				s.student_id,
				COALESCE(s.name, 'Unknown') as student_name,
				COUNT(ase.event_id) as questions_attempted,
				SUM(CASE WHEN ase.is_correct THEN 1 ELSE 0 END) as correct_answers,
				ROUND(AVG(CASE WHEN ase.is_correct THEN 1.0 ELSE 0.0 END) * 100, 2) as accuracy_rate,
				ROUND(AVG(EXTRACT(EPOCH FROM (ase.submitted_at - qpe.published_at))), 2) as average_response_time,
				1 as sessions_participated
			FROM answer_submitted_events ase
			JOIN students s ON ase.student_id = s.student_id
			LEFT JOIN question_published_events qpe ON ase.question_id = qpe.question_id AND ase.session_id = qpe.session_id
			WHERE ase.session_id = ?
			GROUP BY s.student_id, s.name
		),
		ranked_students AS (
			SELECT 
				*,
				ROW_NUMBER() OVER (ORDER BY accuracy_rate DESC, questions_attempted DESC) as rank,
				ROUND(
					(ROW_NUMBER() OVER (ORDER BY accuracy_rate DESC, questions_attempted DESC) - 1) * 100.0 / 
					GREATEST(COUNT(*) OVER() - 1, 1), 2
				) as percentile
			FROM student_stats
		)
		SELECT *
		FROM ranked_students
		ORDER BY rank
		LIMIT ? OFFSET ?
	`, sessionID, pagination.PageSize, pagination.Offset).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	response := NewPaginatedResponse(results, pagination, int(totalCount))
	return &response, nil
}

// NEW: Basic Overview Implementations

// GetClassroomOverview - basic classroom statistics and recent activity
func (r *eventRepository) GetClassroomOverview(classroomID uuid.UUID) (*ClassroomOverviewData, error) {
	var overview ClassroomOverviewData

	err := r.db.Raw(`
		WITH classroom_info AS (
			SELECT 
				c.classroom_id,
				c.name as classroom_name,
				COUNT(cs.student_id) as total_students
			FROM classrooms c
			LEFT JOIN classroom_students cs ON c.classroom_id = cs.classroom_id
			WHERE c.classroom_id = ?
			GROUP BY c.classroom_id, c.name
		),
		activity_stats AS (
			SELECT 
				COUNT(DISTINCT qs.session_id) as total_sessions,
				COUNT(DISTINCT CASE 
					WHEN qs.started_at >= NOW() - INTERVAL '7 days' 
					THEN qs.session_id 
				END) as recent_sessions,
				MAX(qs.started_at) as last_activity,
				COUNT(DISTINCT CASE 
					WHEN ase.submitted_at >= NOW() - INTERVAL '30 days' 
					THEN ase.student_id 
				END) as active_students
			FROM quiz_sessions qs
			LEFT JOIN answer_submitted_events ase ON qs.session_id = ase.session_id
			WHERE qs.classroom_id = ?
		)
		SELECT 
			ci.classroom_id,
			ci.classroom_name,
			ci.total_students,
			COALESCE(as_stats.active_students, 0) as active_students,
			COALESCE(as_stats.total_sessions, 0) as total_sessions,
			COALESCE(as_stats.recent_sessions, 0) as recent_sessions,
			as_stats.last_activity,
			NOW() as created_at
		FROM classroom_info ci
		CROSS JOIN activity_stats as_stats
	`, classroomID, classroomID).Scan(&overview).Error

	return &overview, err
}

// GetClassPerformanceSummary - overall class performance metrics
func (r *eventRepository) GetClassPerformanceSummary(classroomID uuid.UUID) (*ClassPerformanceSummaryData, error) {
	var summary ClassPerformanceSummaryData

	err := r.db.Raw(`
		WITH classroom_info AS (
			SELECT 
				c.classroom_id,
				c.name as classroom_name,
				COUNT(cs.student_id) as total_students
			FROM classrooms c
			LEFT JOIN classroom_students cs ON c.classroom_id = cs.classroom_id
			WHERE c.classroom_id = ?
			GROUP BY c.classroom_id, c.name
		),
		performance_stats AS (
			SELECT 
				COUNT(DISTINCT ase.student_id) as participating_students,
				COUNT(DISTINCT qs.session_id) as session_count,
				COUNT(DISTINCT qs.quiz_id) as total_quizzes_taken,
				COUNT(ase.event_id) as total_questions_answered,
				ROUND(AVG(CASE WHEN ase.is_correct THEN 1.0 ELSE 0.0 END) * 100, 2) as overall_accuracy,
				ROUND(AVG(EXTRACT(EPOCH FROM (ase.submitted_at - qpe.published_at))), 2) as average_response_time
			FROM quiz_sessions qs
			LEFT JOIN answer_submitted_events ase ON qs.session_id = ase.session_id
			LEFT JOIN question_published_events qpe ON ase.question_id = qpe.question_id 
				AND ase.session_id = qpe.session_id
			WHERE qs.classroom_id = ?
		)
		SELECT 
			ci.classroom_id,
			ci.classroom_name,
			ci.total_students,
			COALESCE(ps_stats.participating_students, 0) as participating_students,
			COALESCE(ps_stats.overall_accuracy, 0) as overall_accuracy,
			ROUND(
				COALESCE(ps_stats.participating_students, 0) * 100.0 / 
				GREATEST(ci.total_students, 1), 2
			) as overall_participation_rate,
			COALESCE(ps_stats.total_quizzes_taken, 0) as total_quizzes_taken,
			COALESCE(ps_stats.total_questions_answered, 0) as total_questions_answered,
			COALESCE(ps_stats.average_response_time, 0) as average_response_time,
			COALESCE(ps_stats.session_count, 0) as session_count
		FROM classroom_info ci
		CROSS JOIN performance_stats ps_stats
	`, classroomID, classroomID).Scan(&summary).Error

	return &summary, err
}

// GetStudentActivitySummary - individual student participation and performance summary
func (r *eventRepository) GetStudentActivitySummary(studentID, classroomID uuid.UUID) (*StudentActivitySummaryData, error) {
	var summary StudentActivitySummaryData

	err := r.db.Raw(`
		WITH student_info AS (
			SELECT 
				s.student_id,
				COALESCE(s.name, 'Unknown') as student_name,
				c.classroom_id,
				c.name as classroom_name
			FROM students s
			JOIN classroom_students cs ON s.student_id = cs.student_id
			JOIN classrooms c ON cs.classroom_id = c.classroom_id
			WHERE s.student_id = ? AND c.classroom_id = ?
		),
		activity_stats AS (
			SELECT 
				COUNT(DISTINCT ase.session_id) as total_sessions_participated,
				COUNT(DISTINCT qs.quiz_id) as unique_quizzes_taken,
				COUNT(ase.event_id) as total_questions_answered,
				ROUND(AVG(CASE WHEN ase.is_correct THEN 1.0 ELSE 0.0 END) * 100, 2) as overall_accuracy,
				ROUND(AVG(EXTRACT(EPOCH FROM (ase.submitted_at - qpe.published_at))), 2) as average_response_time,
				MIN(ase.submitted_at) as first_activity,
				MAX(ase.submitted_at) as last_activity
			FROM answer_submitted_events ase
			JOIN quiz_sessions qs ON ase.session_id = qs.session_id
			LEFT JOIN question_published_events qpe ON ase.question_id = qpe.question_id 
				AND ase.session_id = qpe.session_id
			WHERE ase.student_id = ? AND qs.classroom_id = ?
		)
		SELECT 
			si.student_id,
			si.student_name,
			si.classroom_id,
			si.classroom_name,
			COALESCE(act_stats.total_sessions_participated, 0) as total_sessions_participated,
			COALESCE(act_stats.unique_quizzes_taken, 0) as unique_quizzes_taken,
			COALESCE(act_stats.total_questions_answered, 0) as total_questions_answered,
			COALESCE(act_stats.overall_accuracy, 0) as overall_accuracy,
			COALESCE(act_stats.average_response_time, 0) as average_response_time,
			act_stats.first_activity,
			act_stats.last_activity
		FROM student_info si
		CROSS JOIN activity_stats act_stats
	`, studentID, classroomID, studentID, classroomID).Scan(&summary).Error

	return &summary, err
}
