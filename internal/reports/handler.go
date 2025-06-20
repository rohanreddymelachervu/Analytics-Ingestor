package reports

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rohanreddymelachervu/ingestor/internal/analytics"
	"github.com/rohanreddymelachervu/ingestor/internal/repository"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

// Helper function to parse pagination parameters
func parsePaginationParams(c *gin.Context) repository.PaginationParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))
	return repository.NewPaginationParams(page, pageSize)
}

func (h *Handler) GetActiveParticipants(c *gin.Context) {
	sessionIDStr := c.Query("session_id")
	if sessionIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id is required"})
		return
	}

	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session_id format"})
		return
	}

	timeRangeStr := c.DefaultQuery("time_range", "60m")
	timeRange, err := time.ParseDuration(timeRangeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid time_range format"})
		return
	}

	// Parse pagination parameters
	pagination := parsePaginationParams(c)

	data, err := h.service.GetActiveParticipants(sessionID, timeRange, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *Handler) GetQuestionsPerMinute(c *gin.Context) {
	sessionIDStr := c.Query("session_id")
	if sessionIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id is required"})
		return
	}

	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session_id format"})
		return
	}

	data, err := h.service.GetQuestionsPerMinute(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *Handler) GetStudentPerformance(c *gin.Context) {
	studentIDStr := c.Query("student_id")
	classroomIDStr := c.Query("classroom_id")

	if studentIDStr == "" || classroomIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "student_id and classroom_id are required"})
		return
	}

	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student_id format"})
		return
	}

	classroomID, err := uuid.Parse(classroomIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid classroom_id format"})
		return
	}

	data, err := h.service.GetStudentPerformance(studentID, classroomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *Handler) GetClassroomEngagement(c *gin.Context) {
	classroomIDStr := c.Query("classroom_id")
	if classroomIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "classroom_id is required"})
		return
	}

	classroomID, err := uuid.Parse(classroomIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid classroom_id format"})
		return
	}

	dateRangeStr := c.DefaultQuery("date_range", "7d")
	dateRange, err := parseDateRange(dateRangeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date_range format"})
		return
	}

	data, err := h.service.GetClassroomEngagement(classroomID, dateRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *Handler) GetContentEffectiveness(c *gin.Context) {
	quizIDStr := c.Query("quiz_id")
	if quizIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quiz_id is required"})
		return
	}

	quizID, err := uuid.Parse(quizIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quiz_id format"})
		return
	}

	data, err := h.service.GetContentEffectiveness(quizID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *Handler) GetResponseRate(c *gin.Context) {
	sessionIDStr := c.Query("session_id")
	questionIDStr := c.Query("question_id")

	if sessionIDStr == "" || questionIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id and question_id are required"})
		return
	}

	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session_id format"})
		return
	}

	questionID, err := uuid.Parse(questionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid question_id format"})
		return
	}

	data, err := h.service.GetResponseRate(sessionID, questionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *Handler) GetLatencyAnalysis(c *gin.Context) {
	sessionIDStr := c.Query("session_id")
	questionIDStr := c.Query("question_id")

	if sessionIDStr == "" || questionIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id and question_id are required"})
		return
	}

	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session_id format"})
		return
	}

	questionID, err := uuid.Parse(questionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid question_id format"})
		return
	}

	data, err := h.service.GetLatencyAnalysis(sessionID, questionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *Handler) GetTimeoutAnalysis(c *gin.Context) {
	sessionIDStr := c.Query("session_id")
	questionIDStr := c.Query("question_id")

	if sessionIDStr == "" || questionIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id and question_id are required"})
		return
	}

	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session_id format"})
		return
	}

	questionID, err := uuid.Parse(questionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid question_id format"})
		return
	}

	data, err := h.service.GetTimeoutAnalysis(sessionID, questionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *Handler) GetCompletionRate(c *gin.Context) {
	sessionIDStr := c.Query("session_id")
	if sessionIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id is required"})
		return
	}

	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session_id format"})
		return
	}

	data, err := h.service.GetCompletionRate(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *Handler) GetDropoffAnalysis(c *gin.Context) {
	sessionIDStr := c.Query("session_id")
	if sessionIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id is required"})
		return
	}

	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session_id format"})
		return
	}

	data, err := h.service.GetDropoffAnalysis(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

// New handler for paginated student performance list
func (h *Handler) GetStudentPerformanceList(c *gin.Context) {
	classroomIDStr := c.Query("classroom_id")
	if classroomIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "classroom_id is required"})
		return
	}

	classroomID, err := uuid.Parse(classroomIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid classroom_id format"})
		return
	}

	// Parse pagination parameters
	pagination := parsePaginationParams(c)

	data, err := h.service.GetStudentPerformanceList(classroomID, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

// New handler for paginated classroom engagement history
func (h *Handler) GetClassroomEngagementHistory(c *gin.Context) {
	classroomIDStr := c.Query("classroom_id")
	if classroomIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "classroom_id is required"})
		return
	}

	classroomID, err := uuid.Parse(classroomIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid classroom_id format"})
		return
	}

	dateRangeStr := c.DefaultQuery("date_range", "7d")
	dateRange, err := parseDateRange(dateRangeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date_range format"})
		return
	}

	// Parse pagination parameters
	pagination := parsePaginationParams(c)

	data, err := h.service.GetClassroomEngagementHistory(classroomID, dateRange, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

// Helper function to parse date ranges like "7d", "30d", "1h"
func parseDateRange(rangeStr string) (time.Duration, error) {
	if len(rangeStr) < 2 {
		return 0, nil
	}

	unit := rangeStr[len(rangeStr)-1]
	valueStr := rangeStr[:len(rangeStr)-1]
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, err
	}

	switch unit {
	case 'd':
		return time.Duration(value) * 24 * time.Hour, nil
	case 'h':
		return time.Duration(value) * time.Hour, nil
	case 'm':
		return time.Duration(value) * time.Minute, nil
	default:
		return time.ParseDuration(rangeStr)
	}
}

// NEW: Missing basic metrics handlers

func (h *Handler) GetQuizSummary(c *gin.Context) {
	quizIDStr := c.Query("quiz_id")
	if quizIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quiz_id is required"})
		return
	}

	quizID, err := uuid.Parse(quizIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quiz_id format"})
		return
	}

	data, err := h.service.GetQuizSummary(quizID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *Handler) GetQuestionAnalysis(c *gin.Context) {
	questionIDStr := c.Query("question_id")
	if questionIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "question_id is required"})
		return
	}

	questionID, err := uuid.Parse(questionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid question_id format"})
		return
	}

	data, err := h.service.GetQuestionAnalysis(questionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *Handler) GetQuizQuestionsList(c *gin.Context) {
	quizIDStr := c.Query("quiz_id")
	if quizIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quiz_id is required"})
		return
	}

	quizID, err := uuid.Parse(quizIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quiz_id format"})
		return
	}

	// Parse pagination parameters
	pagination := parsePaginationParams(c)

	data, err := h.service.GetQuizQuestionsList(quizID, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *Handler) GetClassroomSessions(c *gin.Context) {
	classroomIDStr := c.Query("classroom_id")
	if classroomIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "classroom_id is required"})
		return
	}

	classroomID, err := uuid.Parse(classroomIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid classroom_id format"})
		return
	}

	// Parse pagination parameters
	pagination := parsePaginationParams(c)

	data, err := h.service.GetClassroomSessions(classroomID, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *Handler) GetQuizSessions(c *gin.Context) {
	quizIDStr := c.Query("quiz_id")
	if quizIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quiz_id is required"})
		return
	}

	quizID, err := uuid.Parse(quizIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quiz_id format"})
		return
	}

	// Parse pagination parameters
	pagination := parsePaginationParams(c)

	data, err := h.service.GetQuizSessions(quizID, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *Handler) GetClassroomStudentRankings(c *gin.Context) {
	classroomIDStr := c.Query("classroom_id")
	if classroomIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "classroom_id is required"})
		return
	}

	classroomID, err := uuid.Parse(classroomIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid classroom_id format"})
		return
	}

	// Parse pagination parameters
	pagination := parsePaginationParams(c)

	data, err := h.service.GetClassroomStudentRankings(classroomID, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *Handler) GetSessionStudentRankings(c *gin.Context) {
	sessionIDStr := c.Query("session_id")
	if sessionIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id is required"})
		return
	}

	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session_id format"})
		return
	}

	// Parse pagination parameters
	pagination := parsePaginationParams(c)

	data, err := h.service.GetSessionStudentRankings(sessionID, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

// NEW: Basic Overview HTTP Handlers

// GetClassroomOverview handles GET /api/reports/classroom-overview
func (h *Handler) GetClassroomOverview(c *gin.Context) {
	classroomID, err := uuid.Parse(c.Query("classroom_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid classroom_id format"})
		return
	}

	response, err := h.service.GetClassroomOverview(classroomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get classroom overview"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetClassPerformanceSummary handles GET /api/reports/class-performance-summary
func (h *Handler) GetClassPerformanceSummary(c *gin.Context) {
	classroomID, err := uuid.Parse(c.Query("classroom_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid classroom_id format"})
		return
	}

	response, err := h.service.GetClassPerformanceSummary(classroomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get class performance summary"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetStudentActivitySummary handles GET /api/reports/student-activity-summary
func (h *Handler) GetStudentActivitySummary(c *gin.Context) {
	studentID, err := uuid.Parse(c.Query("student_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student_id format"})
		return
	}

	classroomID, err := uuid.Parse(c.Query("classroom_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid classroom_id format"})
		return
	}

	response, err := h.service.GetStudentActivitySummary(studentID, classroomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get student activity summary"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GenericQuery handles cube.dev-style analytics queries with measures and dimensions
func (h *Handler) GenericQuery(c *gin.Context) {
	var request analytics.QueryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Validate that at least one measure or dimension is requested
	if len(request.Measures) == 0 && len(request.Dimensions) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one measure or dimension must be specified"})
		return
	}

	// Generate SQL from the query request
	sql, err := request.BuildSQL()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to build query", "details": err.Error()})
		return
	}

	// Execute the generated SQL
	results, err := h.service.ExecuteGenericQuery(sql)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query", "details": err.Error()})
		return
	}

	// Return successful response with query metadata
	response := gin.H{
		"query": gin.H{
			"measures":   request.Measures,
			"dimensions": request.Dimensions,
			"filters":    request.Filters,
			"time_range": request.TimeRange,
			"limit":      request.Limit,
			"order_by":   request.OrderBy,
		},
		"data":          results,
		"generated_sql": sql,
		"count":         len(results),
	}

	c.JSON(http.StatusOK, response)
}
