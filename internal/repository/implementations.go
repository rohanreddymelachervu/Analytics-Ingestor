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

func (r *eventRepository) GetActiveParticipants(sessionID uuid.UUID, timeRange time.Duration) ([]ParticipantMetrics, error) {
	var results []ParticipantMetrics

	cutoffTime := time.Now().Add(-timeRange)

	err := r.db.Raw(`
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
	`, sessionID, cutoffTime).Scan(&results).Error

	return results, err
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
	`, classroomID).Scan(&students).Error
	return students, err
}
