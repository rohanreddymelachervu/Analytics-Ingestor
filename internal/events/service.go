package events

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/rohanreddymelachervu/ingestor/internal/models"
	"github.com/rohanreddymelachervu/ingestor/internal/repository"
)

type Service struct {
	EventRepo     repository.EventRepository
	QuizRepo      repository.QuizRepository
	SessionRepo   repository.SessionRepository
	ClassroomRepo repository.ClassroomRepository
}

func NewService(eventRepo repository.EventRepository, quizRepo repository.QuizRepository,
	sessionRepo repository.SessionRepository, classroomRepo repository.ClassroomRepository) *Service {
	return &Service{
		EventRepo:     eventRepo,
		QuizRepo:      quizRepo,
		SessionRepo:   sessionRepo,
		ClassroomRepo: classroomRepo,
	}
}

func (s *Service) ProcessEvent(event EventPayload, userID interface{}) error {
	switch event.EventType {
	case "QUESTION_PUBLISHED":
		return s.processQuestionPublishedEvent(event, userID)
	case "ANSWER_SUBMITTED":
		return s.processAnswerSubmittedEvent(event, userID)
	case "SESSION_STARTED":
		return s.processSessionStartedEvent(event, userID)
	default:
		return fmt.Errorf("unknown event type: %s", event.EventType)
	}
}

func (s *Service) processQuestionPublishedEvent(event EventPayload, userID interface{}) error {
	// Parse UUIDs
	eventID, err := uuid.Parse(event.EventID)
	if err != nil {
		return fmt.Errorf("invalid event_id: %v", err)
	}
	sessionID, err := uuid.Parse(event.SessionID)
	if err != nil {
		return fmt.Errorf("invalid session_id: %v", err)
	}
	questionID, err := uuid.Parse(event.QuestionID)
	if err != nil {
		return fmt.Errorf("invalid question_id: %v", err)
	}

	var teacherID *uuid.UUID
	if event.TeacherID != nil {
		parsed, err := uuid.Parse(*event.TeacherID)
		if err != nil {
			return fmt.Errorf("invalid teacher_id: %v", err)
		}
		teacherID = &parsed
	}

	// Create question published event
	questionEvent := &models.QuestionPublishedEvent{
		EventID:     eventID,
		SessionID:   sessionID,
		QuestionID:  questionID,
		TeacherID:   teacherID,
		PublishedAt: event.Timestamp,
	}

	if event.TimerSec != nil {
		questionEvent.TimerDurationSec = *event.TimerSec
	}

	return s.EventRepo.SaveQuestionPublishedEvent(questionEvent)
}

func (s *Service) processAnswerSubmittedEvent(event EventPayload, userID interface{}) error {
	// Parse UUIDs
	eventID, err := uuid.Parse(event.EventID)
	if err != nil {
		return fmt.Errorf("invalid event_id: %v", err)
	}
	sessionID, err := uuid.Parse(event.SessionID)
	if err != nil {
		return fmt.Errorf("invalid session_id: %v", err)
	}
	questionID, err := uuid.Parse(event.QuestionID)
	if err != nil {
		return fmt.Errorf("invalid question_id: %v", err)
	}

	if event.StudentID == nil || event.Answer == nil {
		return fmt.Errorf("student_id and answer are required for ANSWER_SUBMITTED events")
	}

	studentID, err := uuid.Parse(*event.StudentID)
	if err != nil {
		return fmt.Errorf("invalid student_id: %v", err)
	}

	// Timer validation: check if answer is submitted within allowed time
	err = s.EventRepo.ValidateAnswerTiming(sessionID, questionID, event.Timestamp)
	if err != nil {
		return fmt.Errorf("answer submitted after deadline: %v", err)
	}

	// Determine if answer is correct (A and C = correct, B and D = incorrect)
	isCorrect := (*event.Answer == "A" || *event.Answer == "C")

	// Create answer submitted event
	answerEvent := &models.AnswerSubmittedEvent{
		EventID:     eventID,
		SessionID:   sessionID,
		QuestionID:  questionID,
		StudentID:   studentID,
		Answer:      *event.Answer,
		IsCorrect:   isCorrect,
		SubmittedAt: event.Timestamp,
	}

	return s.EventRepo.SaveAnswerSubmittedEvent(answerEvent)
}

func (s *Service) processSessionStartedEvent(event EventPayload, userID interface{}) error {
	// Parse UUIDs
	sessionID, err := uuid.Parse(event.SessionID)
	if err != nil {
		return fmt.Errorf("invalid session_id: %v", err)
	}
	quizID, err := uuid.Parse(event.QuizID)
	if err != nil {
		return fmt.Errorf("invalid quiz_id: %v", err)
	}
	classroomID, err := uuid.Parse(event.ClassroomID)
	if err != nil {
		return fmt.Errorf("invalid classroom_id: %v", err)
	}

	// Create or update session
	session := &models.QuizSession{
		SessionID:   sessionID,
		QuizID:      quizID,
		ClassroomID: classroomID,
		StartedAt:   event.Timestamp,
	}

	return s.SessionRepo.CreateSession(session)
}
