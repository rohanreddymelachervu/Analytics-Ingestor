package reports

import (
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/rohanreddymelachervu/ingestor/internal/repository"
)

type Service struct {
	EventRepo     repository.EventRepository
	ClassroomRepo repository.ClassroomRepository
}

func NewService(eventRepo repository.EventRepository, classroomRepo repository.ClassroomRepository) *Service {
	return &Service{
		EventRepo:     eventRepo,
		ClassroomRepo: classroomRepo,
	}
}

func (s *Service) GetActiveParticipants(sessionID uuid.UUID, timeRange time.Duration, pagination repository.PaginationParams) (interface{}, error) {
	paginatedData, err := s.EventRepo.GetActiveParticipants(sessionID, timeRange, pagination)
	if err != nil {
		return nil, err
	}

	// Calculate average accuracy from the current page data
	avgAccuracy := calculateAverageAccuracy(paginatedData.Data)

	response := map[string]interface{}{
		"session_id": sessionID,
		"time_range": timeRange.String(),
		"pagination": map[string]interface{}{
			"page":         paginatedData.Page,
			"page_size":    paginatedData.PageSize,
			"total_count":  paginatedData.TotalCount,
			"total_pages":  paginatedData.TotalPages,
			"has_more":     paginatedData.HasMore,
			"has_previous": paginatedData.HasPrevious,
		},
		"active_participants":      paginatedData.Data,
		"total_participants":       paginatedData.TotalCount,
		"page_participants":        len(paginatedData.Data),
		"average_accuracy_percent": avgAccuracy,
	}

	return response, nil
}

func (s *Service) GetQuestionsPerMinute(sessionID uuid.UUID) (interface{}, error) {
	stats, err := s.EventRepo.GetQuestionsPerMinuteStats(sessionID)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"session_id":      sessionID,
		"total_questions": stats.TotalQuestions,
		"average_qpm":     stats.AverageQPM,
		"peak_qpm":        stats.PeakQPM,
		"analysis_window": "session_duration",
	}

	return response, nil
}

func (s *Service) GetStudentPerformance(studentID, classroomID uuid.UUID) (interface{}, error) {
	performance, err := s.EventRepo.GetStudentPerformance(studentID, classroomID)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"student_id":               studentID,
		"classroom_id":             classroomID,
		"questions_attempted":      performance.QuestionsAttempted,
		"correct_answers":          performance.CorrectAnswers,
		"overall_accuracy_percent": performance.OverallAccuracy,
		"average_response_time":    performance.AverageResponseTime,
		"performance_trend":        "stable", // Could be enhanced with historical data
	}

	return response, nil
}

func (s *Service) GetClassroomEngagement(classroomID uuid.UUID, dateRange time.Duration) (interface{}, error) {
	engagement, err := s.EventRepo.GetClassroomEngagement(classroomID, dateRange)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"classroom_id":             classroomID,
		"date_range":               dateRange.String(),
		"total_students":           engagement.TotalStudents,
		"active_students":          engagement.ActiveStudents,
		"engagement_rate_percent":  engagement.EngagementRate,
		"average_accuracy_percent": engagement.AverageAccuracy,
		"total_questions":          engagement.TotalQuestions,
		"response_rate_percent":    engagement.ResponseRate,
		"trends": map[string]string{
			"engagement": "increasing",
			"accuracy":   "stable",
		},
	}

	return response, nil
}

func (s *Service) GetContentEffectiveness(quizID uuid.UUID) (interface{}, error) {
	effectiveness, err := s.EventRepo.GetContentEffectiveness(quizID)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"quiz_id":             effectiveness.QuizID,
		"total_questions":     effectiveness.TotalQuestions,
		"average_accuracy":    effectiveness.AverageAccuracy,
		"overall_engagement":  effectiveness.OverallEngagement,
		"effectiveness_score": effectiveness.EffectivenessScore,
		"recommendations":     effectiveness.Recommendations,
		"content_analysis": map[string]interface{}{
			"difficulty_level": "moderate",
			"engagement_level": "high",
			"optimization_suggestions": []string{
				"Consider adding more interactive elements",
				"Review questions with low accuracy rates",
			},
		},
	}

	return response, nil
}

func (s *Service) GetResponseRate(sessionID, questionID uuid.UUID) (interface{}, error) {
	data, err := s.EventRepo.GetResponseRate(sessionID, questionID)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"question_id":       data.QuestionID,
		"session_id":        data.SessionID,
		"students_received": data.StudentsReceived,
		"students_answered": data.StudentsAnswered,
		"response_rate":     data.ResponseRate,
		"analysis": map[string]interface{}{
			"engagement_level": getEngagementLevel(data.ResponseRate),
			"benchmark":        "Industry average: 85%",
		},
	}

	return response, nil
}

func (s *Service) GetLatencyAnalysis(sessionID, questionID uuid.UUID) (interface{}, error) {
	data, err := s.EventRepo.GetLatencyToFirstAnswer(sessionID, questionID)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"question_id":          data.QuestionID,
		"session_id":           data.SessionID,
		"first_answer_latency": data.FirstAnswerLatency.String(),
		"average_latency":      data.AverageLatency.String(),
		"median_latency":       data.MedianLatency.String(),
		"analysis": map[string]interface{}{
			"speed_rating": getSpeedRating(data.AverageLatency),
			"benchmark":    "Target: <30s for optimal engagement",
		},
	}

	return response, nil
}

func (s *Service) GetTimeoutAnalysis(sessionID, questionID uuid.UUID) (interface{}, error) {
	data, err := s.EventRepo.GetTimeoutAndSkippedRate(sessionID, questionID)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"question_id":    data.QuestionID,
		"session_id":     data.SessionID,
		"total_students": data.TotalStudents,
		"timeout_count":  data.TimeoutCount,
		"skipped_count":  data.SkippedCount,
		"timeout_rate":   data.TimeoutRate,
		"skipped_rate":   data.SkippedRate,
		"analysis": map[string]interface{}{
			"difficulty_indicator": getDifficultyIndicator(data.TimeoutRate, data.SkippedRate),
			"recommendation":       getTimeoutRecommendation(data.TimeoutRate),
		},
	}

	return response, nil
}

func (s *Service) GetCompletionRate(sessionID uuid.UUID) (interface{}, error) {
	data, err := s.EventRepo.GetCompletionRate(sessionID)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"session_id":         data.SessionID,
		"total_students":     data.TotalStudents,
		"completed_students": data.CompletedStudents,
		"completion_rate":    data.CompletionRate,
		"total_questions":    data.TotalQuestions,
		"average_completion": data.AverageCompletion,
		"analysis": map[string]interface{}{
			"retention_level": getRetentionLevel(data.CompletionRate),
			"benchmark":       "Target completion rate: >80%",
		},
	}

	return response, nil
}

func (s *Service) GetDropoffAnalysis(sessionID uuid.UUID) (interface{}, error) {
	dropoffPoints, err := s.EventRepo.GetDropoffPoints(sessionID)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"session_id":     sessionID,
		"dropoff_points": dropoffPoints,
		"analysis": map[string]interface{}{
			"critical_points": getCriticalDropoffPoints(dropoffPoints),
			"recommendations": getDropoffRecommendations(dropoffPoints),
		},
	}

	return response, nil
}

// Helper functions for analysis
func calculateAverageAccuracy(participants []repository.ParticipantMetrics) float64 {
	if len(participants) == 0 {
		return 0.0
	}

	total := 0.0
	for _, p := range participants {
		total += p.Accuracy
	}
	return total / float64(len(participants))
}

func getEngagementLevel(responseRate float64) string {
	if responseRate >= 90 {
		return "excellent"
	} else if responseRate >= 75 {
		return "good"
	} else if responseRate >= 60 {
		return "moderate"
	}
	return "low"
}

func getSpeedRating(avgLatency time.Duration) string {
	if avgLatency <= 15*time.Second {
		return "fast"
	} else if avgLatency <= 30*time.Second {
		return "moderate"
	}
	return "slow"
}

func getDifficultyIndicator(timeoutRate, skippedRate float64) string {
	combined := timeoutRate + skippedRate
	if combined >= 30 {
		return "high_difficulty"
	} else if combined >= 15 {
		return "moderate_difficulty"
	}
	return "appropriate_difficulty"
}

func getTimeoutRecommendation(timeoutRate float64) string {
	if timeoutRate >= 20 {
		return "Consider increasing time limit or simplifying question"
	} else if timeoutRate >= 10 {
		return "Monitor question difficulty and time allocation"
	}
	return "Time allocation appears appropriate"
}

func getRetentionLevel(completionRate float64) string {
	if completionRate >= 85 {
		return "excellent"
	} else if completionRate >= 70 {
		return "good"
	} else if completionRate >= 50 {
		return "moderate"
	}
	return "concerning"
}

func getCriticalDropoffPoints(dropoffs []repository.DropoffPoint) []repository.DropoffPoint {
	var critical []repository.DropoffPoint
	for _, point := range dropoffs {
		if point.DropoffRate >= 25.0 {
			critical = append(critical, point)
		}
	}
	return critical
}

func getDropoffRecommendations(dropoffs []repository.DropoffPoint) []string {
	recommendations := []string{}

	for _, point := range dropoffs {
		if point.DropoffRate >= 30 {
			recommendations = append(recommendations,
				"Review question difficulty and clarity for high drop-off points")
			break
		}
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Retention rates are healthy")
	}

	return recommendations
}

// New paginated service for student performance list
func (s *Service) GetStudentPerformanceList(classroomID uuid.UUID, pagination repository.PaginationParams) (interface{}, error) {
	paginatedData, err := s.EventRepo.GetStudentPerformanceList(classroomID, pagination)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"classroom_id": classroomID,
		"pagination": map[string]interface{}{
			"page":         paginatedData.Page,
			"page_size":    paginatedData.PageSize,
			"total_count":  paginatedData.TotalCount,
			"total_pages":  paginatedData.TotalPages,
			"has_more":     paginatedData.HasMore,
			"has_previous": paginatedData.HasPrevious,
		},
		"students": paginatedData.Data,
		"summary": map[string]interface{}{
			"total_students": paginatedData.TotalCount,
			"page_students":  len(paginatedData.Data),
		},
	}

	return response, nil
}

// New paginated service for classroom engagement history
func (s *Service) GetClassroomEngagementHistory(classroomID uuid.UUID, dateRange time.Duration, pagination repository.PaginationParams) (interface{}, error) {
	paginatedData, err := s.EventRepo.GetClassroomEngagementHistory(classroomID, dateRange, pagination)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"classroom_id": classroomID,
		"date_range":   dateRange.String(),
		"pagination": map[string]interface{}{
			"page":         paginatedData.Page,
			"page_size":    paginatedData.PageSize,
			"total_count":  paginatedData.TotalCount,
			"total_pages":  paginatedData.TotalPages,
			"has_more":     paginatedData.HasMore,
			"has_previous": paginatedData.HasPrevious,
		},
		"engagement_history": paginatedData.Data,
		"summary": map[string]interface{}{
			"total_periods": paginatedData.TotalCount,
			"page_periods":  len(paginatedData.Data),
		},
	}

	return response, nil
}

// NEW: Missing basic metrics service methods

func (s *Service) GetQuizSummary(quizID uuid.UUID) (interface{}, error) {
	summary, err := s.EventRepo.GetQuizSummary(quizID)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"quiz_id": summary.QuizID,
		"title":   summary.Title,
		"usage_statistics": map[string]interface{}{
			"total_sessions":   summary.TotalSessions,
			"total_classrooms": summary.TotalClassrooms,
			"total_students":   summary.TotalStudents,
			"total_questions":  summary.TotalQuestions,
			"first_used":       summary.FirstUsed,
			"last_used":        summary.LastUsed,
		},
		"performance_metrics": map[string]interface{}{
			"average_accuracy":    summary.AverageAccuracy,
			"average_completion":  summary.AverageCompletion,
			"overall_engagement":  summary.OverallEngagement,
			"effectiveness_score": summary.EffectivenessScore,
		},
		"insights": map[string]interface{}{
			"performance_rating": getPerformanceRating(summary.EffectivenessScore),
			"usage_frequency":    getUsageFrequency(summary.TotalSessions),
			"reach":              getRatingReach(summary.TotalClassrooms),
		},
	}

	return response, nil
}

func (s *Service) GetQuestionAnalysis(questionID uuid.UUID) (interface{}, error) {
	analysis, err := s.EventRepo.GetQuestionAnalysis(questionID)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"question_id": analysis.QuestionID,
		"quiz_id":     analysis.QuizID,
		"usage_stats": map[string]interface{}{
			"total_attempts":   analysis.TotalAttempts,
			"correct_attempts": analysis.CorrectAttempts,
			"usage_count":      analysis.UsageCount,
		},
		"performance_metrics": map[string]interface{}{
			"accuracy_rate":         analysis.AccuracyRate,
			"average_response_time": analysis.AverageResponseTime,
			"difficulty_rating":     analysis.DifficultyRating,
		},
		"answer_distribution": analysis.AnswerDistribution,
		"insights": map[string]interface{}{
			"difficulty_level": getDifficultyInsight(analysis.AccuracyRate),
			"response_quality": getResponseQuality(analysis.AverageResponseTime),
			"effectiveness":    getQuestionEffectiveness(analysis.AccuracyRate, analysis.TotalAttempts),
		},
	}

	return response, nil
}

func (s *Service) GetQuizQuestionsList(quizID uuid.UUID, pagination repository.PaginationParams) (interface{}, error) {
	paginatedData, err := s.EventRepo.GetQuizQuestionsList(quizID, pagination)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"quiz_id": quizID,
		"pagination": map[string]interface{}{
			"page":         paginatedData.Page,
			"page_size":    paginatedData.PageSize,
			"total_count":  paginatedData.TotalCount,
			"total_pages":  paginatedData.TotalPages,
			"has_more":     paginatedData.HasMore,
			"has_previous": paginatedData.HasPrevious,
		},
		"questions": paginatedData.Data,
		"summary": map[string]interface{}{
			"total_questions": paginatedData.TotalCount,
			"page_questions":  len(paginatedData.Data),
		},
	}

	return response, nil
}

func (s *Service) GetClassroomSessions(classroomID uuid.UUID, pagination repository.PaginationParams) (interface{}, error) {
	paginatedData, err := s.EventRepo.GetClassroomSessions(classroomID, pagination)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"classroom_id": classroomID,
		"pagination": map[string]interface{}{
			"page":         paginatedData.Page,
			"page_size":    paginatedData.PageSize,
			"total_count":  paginatedData.TotalCount,
			"total_pages":  paginatedData.TotalPages,
			"has_more":     paginatedData.HasMore,
			"has_previous": paginatedData.HasPrevious,
		},
		"sessions": paginatedData.Data,
		"summary": map[string]interface{}{
			"total_sessions": paginatedData.TotalCount,
			"page_sessions":  len(paginatedData.Data),
		},
	}

	return response, nil
}

func (s *Service) GetQuizSessions(quizID uuid.UUID, pagination repository.PaginationParams) (interface{}, error) {
	paginatedData, err := s.EventRepo.GetQuizSessions(quizID, pagination)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"quiz_id": quizID,
		"pagination": map[string]interface{}{
			"page":         paginatedData.Page,
			"page_size":    paginatedData.PageSize,
			"total_count":  paginatedData.TotalCount,
			"total_pages":  paginatedData.TotalPages,
			"has_more":     paginatedData.HasMore,
			"has_previous": paginatedData.HasPrevious,
		},
		"sessions": paginatedData.Data,
		"summary": map[string]interface{}{
			"total_sessions": paginatedData.TotalCount,
			"page_sessions":  len(paginatedData.Data),
		},
	}

	return response, nil
}

func (s *Service) GetClassroomStudentRankings(classroomID uuid.UUID, pagination repository.PaginationParams) (interface{}, error) {
	paginatedData, err := s.EventRepo.GetClassroomStudentRankings(classroomID, pagination)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"classroom_id": classroomID,
		"pagination": map[string]interface{}{
			"page":         paginatedData.Page,
			"page_size":    paginatedData.PageSize,
			"total_count":  paginatedData.TotalCount,
			"total_pages":  paginatedData.TotalPages,
			"has_more":     paginatedData.HasMore,
			"has_previous": paginatedData.HasPrevious,
		},
		"rankings": paginatedData.Data,
		"summary": map[string]interface{}{
			"total_students": paginatedData.TotalCount,
			"page_students":  len(paginatedData.Data),
		},
	}

	return response, nil
}

func (s *Service) GetSessionStudentRankings(sessionID uuid.UUID, pagination repository.PaginationParams) (interface{}, error) {
	paginatedData, err := s.EventRepo.GetSessionStudentRankings(sessionID, pagination)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"session_id": sessionID,
		"pagination": map[string]interface{}{
			"page":         paginatedData.Page,
			"page_size":    paginatedData.PageSize,
			"total_count":  paginatedData.TotalCount,
			"total_pages":  paginatedData.TotalPages,
			"has_more":     paginatedData.HasMore,
			"has_previous": paginatedData.HasPrevious,
		},
		"rankings": paginatedData.Data,
		"summary": map[string]interface{}{
			"total_students": paginatedData.TotalCount,
			"page_students":  len(paginatedData.Data),
		},
	}

	return response, nil
}

// Helper functions for insights and analysis
func getPerformanceRating(score float64) string {
	if score >= 90 {
		return "excellent"
	} else if score >= 75 {
		return "good"
	} else if score >= 60 {
		return "average"
	} else if score >= 40 {
		return "below_average"
	}
	return "needs_improvement"
}

func getUsageFrequency(sessions int) string {
	if sessions >= 20 {
		return "high"
	} else if sessions >= 10 {
		return "moderate"
	} else if sessions >= 5 {
		return "low"
	}
	return "minimal"
}

func getRatingReach(classrooms int) string {
	if classrooms >= 10 {
		return "wide"
	} else if classrooms >= 5 {
		return "moderate"
	} else if classrooms >= 2 {
		return "limited"
	}
	return "single_classroom"
}

func getDifficultyInsight(accuracy float64) string {
	if accuracy >= 90 {
		return "too_easy"
	} else if accuracy >= 70 {
		return "appropriate"
	} else if accuracy >= 50 {
		return "challenging"
	}
	return "too_difficult"
}

func getResponseQuality(responseTime float64) string {
	if responseTime <= 10 {
		return "quick"
	} else if responseTime <= 30 {
		return "moderate"
	} else if responseTime <= 60 {
		return "slow"
	}
	return "very_slow"
}

func getQuestionEffectiveness(accuracy float64, attempts int) string {
	if accuracy >= 70 && attempts >= 10 {
		return "highly_effective"
	} else if accuracy >= 50 && attempts >= 5 {
		return "effective"
	} else if attempts < 5 {
		return "needs_more_data"
	}
	return "needs_improvement"
}

// NEW: Basic Overview Service Methods

func (s *Service) GetClassroomOverview(classroomID uuid.UUID) (interface{}, error) {
	overview, err := s.EventRepo.GetClassroomOverview(classroomID)
	if err != nil {
		return nil, err
	}

	// Add insights based on the data
	insights := map[string]interface{}{
		"activity_level":    getActivityLevel(overview.RecentSessions, overview.TotalSessions),
		"engagement_status": getEngagementStatus(overview.ActiveStudents, overview.TotalStudents),
		"growth_trend":      getGrowthTrend(overview.RecentSessions),
	}

	response := map[string]interface{}{
		"overview": overview,
		"insights": insights,
		"summary": map[string]interface{}{
			"participation_rate": calculateParticipationRate(overview.ActiveStudents, overview.TotalStudents),
			"activity_score":     calculateActivityScore(overview.RecentSessions, overview.TotalSessions),
		},
	}

	return response, nil
}

func (s *Service) GetClassPerformanceSummary(classroomID uuid.UUID) (interface{}, error) {
	summary, err := s.EventRepo.GetClassPerformanceSummary(classroomID)
	if err != nil {
		return nil, err
	}

	// Add performance insights
	insights := map[string]interface{}{
		"performance_level":   getPerformanceLevel(summary.OverallAccuracy),
		"participation_level": getParticipationLevel(summary.OverallParticipation),
		"engagement_quality":  getEngagementQuality(summary.TotalQuestionsAnswered, summary.ParticipatingStudents),
		"response_speed":      getResponseSpeed(summary.AverageResponseTime),
	}

	response := map[string]interface{}{
		"performance_summary": summary,
		"insights":            insights,
		"benchmarks": map[string]interface{}{
			"target_accuracy":       75.0,
			"target_participation":  80.0,
			"optimal_response_time": 30.0,
		},
	}

	return response, nil
}

func (s *Service) GetStudentActivitySummary(studentID, classroomID uuid.UUID) (interface{}, error) {
	summary, err := s.EventRepo.GetStudentActivitySummary(studentID, classroomID)
	if err != nil {
		return nil, err
	}

	// Add activity insights
	insights := map[string]interface{}{
		"activity_level":         getStudentActivityLevel(summary.TotalSessionsParticipated),
		"performance_trend":      getStudentPerformanceTrend(summary.OverallAccuracy),
		"engagement_consistency": getEngagementConsistency(summary.TotalSessionsParticipated, summary.UniqueQuizzesTaken),
		"response_efficiency":    getResponseEfficiency(summary.AverageResponseTime),
	}

	response := map[string]interface{}{
		"activity_summary": summary,
		"insights":         insights,
		"metrics": map[string]interface{}{
			"questions_per_session": calculateQuestionsPerSession(summary.TotalQuestionsAnswered, summary.TotalSessionsParticipated),
			"quiz_variety_score":    calculateQuizVarietyScore(summary.UniqueQuizzesTaken, summary.TotalSessionsParticipated),
		},
	}

	return response, nil
}

// NEW: Additional helper functions for overview insights

func getActivityLevel(recentSessions, totalSessions int) string {
	if totalSessions == 0 {
		return "no_activity"
	}
	activityRate := float64(recentSessions) / float64(totalSessions) * 100
	if activityRate >= 50 {
		return "high"
	} else if activityRate >= 25 {
		return "moderate"
	} else if activityRate > 0 {
		return "low"
	}
	return "inactive"
}

func getEngagementStatus(activeStudents, totalStudents int) string {
	if totalStudents == 0 {
		return "no_students"
	}
	engagementRate := float64(activeStudents) / float64(totalStudents) * 100
	if engagementRate >= 80 {
		return "highly_engaged"
	} else if engagementRate >= 60 {
		return "well_engaged"
	} else if engagementRate >= 40 {
		return "moderately_engaged"
	} else if engagementRate > 0 {
		return "low_engagement"
	}
	return "no_engagement"
}

func getGrowthTrend(recentSessions int) string {
	if recentSessions >= 5 {
		return "accelerating"
	} else if recentSessions >= 3 {
		return "growing"
	} else if recentSessions >= 1 {
		return "steady"
	}
	return "stagnant"
}

func calculateParticipationRate(activeStudents, totalStudents int) float64 {
	if totalStudents == 0 {
		return 0.0
	}
	return math.Round(float64(activeStudents)/float64(totalStudents)*10000) / 100
}

func calculateActivityScore(recentSessions, totalSessions int) float64 {
	if totalSessions == 0 {
		return 0.0
	}
	score := float64(recentSessions) / float64(totalSessions) * 100
	return math.Round(score*100) / 100
}

func getPerformanceLevel(accuracy float64) string {
	if accuracy >= 90 {
		return "excellent"
	} else if accuracy >= 75 {
		return "good"
	} else if accuracy >= 60 {
		return "satisfactory"
	} else if accuracy >= 40 {
		return "needs_improvement"
	}
	return "poor"
}

func getParticipationLevel(participation float64) string {
	if participation >= 90 {
		return "outstanding"
	} else if participation >= 75 {
		return "strong"
	} else if participation >= 50 {
		return "moderate"
	} else if participation >= 25 {
		return "weak"
	}
	return "very_low"
}

func getEngagementQuality(totalQuestions, participatingStudents int) string {
	if participatingStudents == 0 {
		return "no_engagement"
	}
	questionsPerStudent := float64(totalQuestions) / float64(participatingStudents)
	if questionsPerStudent >= 20 {
		return "deep_engagement"
	} else if questionsPerStudent >= 10 {
		return "good_engagement"
	} else if questionsPerStudent >= 5 {
		return "moderate_engagement"
	} else if questionsPerStudent >= 1 {
		return "light_engagement"
	}
	return "minimal_engagement"
}

func getResponseSpeed(avgResponseTime float64) string {
	if avgResponseTime <= 10 {
		return "very_fast"
	} else if avgResponseTime <= 30 {
		return "fast"
	} else if avgResponseTime <= 60 {
		return "moderate"
	} else if avgResponseTime <= 120 {
		return "slow"
	}
	return "very_slow"
}

func getStudentActivityLevel(sessionsParticipated int) string {
	if sessionsParticipated >= 10 {
		return "highly_active"
	} else if sessionsParticipated >= 5 {
		return "active"
	} else if sessionsParticipated >= 2 {
		return "moderately_active"
	} else if sessionsParticipated >= 1 {
		return "minimally_active"
	}
	return "inactive"
}

func getStudentPerformanceTrend(accuracy float64) string {
	return getPerformanceLevel(accuracy) // Reuse existing logic
}

func getEngagementConsistency(sessions, quizzes int) string {
	if sessions == 0 {
		return "no_data"
	}
	ratio := float64(quizzes) / float64(sessions)
	if ratio >= 0.8 {
		return "very_consistent"
	} else if ratio >= 0.6 {
		return "consistent"
	} else if ratio >= 0.4 {
		return "somewhat_consistent"
	} else if ratio >= 0.2 {
		return "inconsistent"
	}
	return "very_inconsistent"
}

func getResponseEfficiency(avgResponseTime float64) string {
	return getResponseSpeed(avgResponseTime) // Reuse existing logic
}

func calculateQuestionsPerSession(totalQuestions, totalSessions int) float64 {
	if totalSessions == 0 {
		return 0.0
	}
	return math.Round(float64(totalQuestions)/float64(totalSessions)*100) / 100
}

func calculateQuizVarietyScore(uniqueQuizzes, totalSessions int) float64 {
	if totalSessions == 0 {
		return 0
	}
	return math.Min(float64(uniqueQuizzes)/float64(totalSessions)*100, 100)
}

// ExecuteGenericQuery executes a generic SQL query for cube.dev-style analytics
func (s *Service) ExecuteGenericQuery(sql string) ([]map[string]interface{}, error) {
	return s.EventRepo.ExecuteGenericQuery(sql)
}
