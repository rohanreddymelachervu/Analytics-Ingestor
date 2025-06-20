package analytics

import (
	"fmt"
	"strings"
	"time"
)

// Measure represents a quantitative metric that can be aggregated
type Measure struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Type        string `json:"type"` // count, sum, avg, min, max
	SQL         string `json:"sql"`
	Format      string `json:"format,omitempty"`
}

// Dimension represents a categorical attribute for grouping/filtering
type Dimension struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Type        string `json:"type"` // string, number, time
	SQL         string `json:"sql"`
	Format      string `json:"format,omitempty"`
}

// QueryRequest represents a generic analytics query
type QueryRequest struct {
	Measures   []string          `json:"measures"`
	Dimensions []string          `json:"dimensions"`
	Filters    map[string]string `json:"filters,omitempty"`
	TimeRange  *TimeRange        `json:"time_range,omitempty"`
	Limit      int               `json:"limit,omitempty"`
	OrderBy    []OrderBy         `json:"order_by,omitempty"`
}

type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type OrderBy struct {
	Field string `json:"field"`
	Order string `json:"order"` // ASC, DESC
}

// Analytics cube definitions
var QuizAnalyticsCube = map[string]interface{}{
	"measures": map[string]Measure{
		"total_answers": {
			Name:        "total_answers",
			DisplayName: "Total Answers",
			Type:        "count",
			SQL:         "COUNT(ase.event_id)",
		},
		"correct_answers": {
			Name:        "correct_answers",
			DisplayName: "Correct Answers",
			Type:        "count",
			SQL:         "COUNT(CASE WHEN ase.is_correct = true THEN 1 END)",
		},
		"accuracy_rate": {
			Name:        "accuracy_rate",
			DisplayName: "Accuracy Rate",
			Type:        "avg",
			SQL:         "ROUND(AVG(CASE WHEN ase.is_correct THEN 100.0 ELSE 0.0 END), 2)",
			Format:      "percentage",
		},
		"active_students": {
			Name:        "active_students",
			DisplayName: "Active Students",
			Type:        "count",
			SQL:         "COUNT(DISTINCT ase.student_id)",
		},
		"questions_published": {
			Name:        "questions_published",
			DisplayName: "Questions Published",
			Type:        "count",
			SQL:         "COUNT(DISTINCT qpe.question_id)",
		},
		// STUDENT PERFORMANCE ANALYSIS MEASURES
		"wrong_answers": {
			Name:        "wrong_answers",
			DisplayName: "Wrong Answers",
			Type:        "count",
			SQL:         "COUNT(CASE WHEN ase.is_correct = false THEN 1 END)",
		},
		"performance_variance": {
			Name:        "performance_variance",
			DisplayName: "Performance Variance",
			Type:        "variance",
			SQL:         "VARIANCE(CASE WHEN ase.is_correct THEN 100.0 ELSE 0.0 END)",
		},
		"student_attempts_per_question": {
			Name:        "student_attempts_per_question",
			DisplayName: "Avg Attempts Per Question",
			Type:        "avg",
			SQL:         "ROUND(COUNT(ase.event_id) * 1.0 / GREATEST(COUNT(DISTINCT qpe.question_id), 1), 2)",
		},
		// CLASSROOM ENGAGEMENT METRICS
		"participation_rate": {
			Name:        "participation_rate",
			DisplayName: "Participation Rate",
			Type:        "percentage",
			SQL:         "ROUND(COUNT(DISTINCT ase.student_id) * 100.0 / GREATEST(COUNT(DISTINCT s.student_id), 1), 2)",
			Format:      "percentage",
		},
		"engagement_score": {
			Name:        "engagement_score",
			DisplayName: "Engagement Score",
			Type:        "calculated",
			SQL:         "ROUND((COUNT(DISTINCT ase.student_id) * 100.0 / GREATEST(COUNT(DISTINCT s.student_id), 1) + AVG(CASE WHEN ase.is_correct THEN 100.0 ELSE 0.0 END)) / 2, 2)",
			Format:      "percentage",
		},
		"session_completion_rate": {
			Name:        "session_completion_rate",
			DisplayName: "Session Completion Rate",
			Type:        "percentage",
			SQL:         "ROUND(COUNT(ase.event_id) * 100.0 / GREATEST(COUNT(DISTINCT ase.student_id) * COUNT(DISTINCT qpe.question_id), 1), 2)",
			Format:      "percentage",
		},
		"unique_sessions": {
			Name:        "unique_sessions",
			DisplayName: "Unique Sessions",
			Type:        "count",
			SQL:         "COUNT(DISTINCT qs.session_id)",
		},
		"average_session_duration": {
			Name:        "average_session_duration",
			DisplayName: "Average Session Duration",
			Type:        "avg",
			SQL:         "ROUND(AVG(EXTRACT(EPOCH FROM (qs.ended_at - qs.started_at))), 2)",
			Format:      "seconds",
		},
		"questions_per_minute": {
			Name:        "questions_per_minute",
			DisplayName: "Questions Per Minute",
			Type:        "calculated",
			SQL:         "ROUND(COUNT(DISTINCT qpe.question_id) / GREATEST(EXTRACT(EPOCH FROM (MAX(qpe.published_at) - MIN(qpe.published_at))) / 60, 1), 2)",
		},
		// CONTENT EFFECTIVENESS EVALUATION
		"question_difficulty_score": {
			Name:        "question_difficulty_score",
			DisplayName: "Question Difficulty Score",
			Type:        "calculated",
			SQL:         "ROUND(100 - AVG(CASE WHEN ase.is_correct THEN 100.0 ELSE 0.0 END), 2)",
			Format:      "percentage",
		},
		"content_effectiveness_score": {
			Name:        "content_effectiveness_score",
			DisplayName: "Content Effectiveness Score",
			Type:        "calculated",
			SQL:         "ROUND((AVG(CASE WHEN ase.is_correct THEN 100.0 ELSE 0.0 END) + (COUNT(DISTINCT ase.student_id) * 100.0 / GREATEST(COUNT(DISTINCT s.student_id), 1))) / 2, 2)",
			Format:      "percentage",
		},
		"time_to_first_answer": {
			Name:        "time_to_first_answer",
			DisplayName: "Time to First Answer",
			Type:        "avg",
			SQL:         "ROUND(AVG(EXTRACT(EPOCH FROM (ase.submitted_at - qpe.published_at))), 2)",
			Format:      "seconds",
		},
		"question_engagement_rate": {
			Name:        "question_engagement_rate",
			DisplayName: "Question Engagement Rate",
			Type:        "percentage",
			SQL:         "ROUND(COUNT(ase.event_id) * 100.0 / GREATEST(COUNT(DISTINCT s.student_id), 1), 2)",
			Format:      "percentage",
		},
		"quiz_completion_rate": {
			Name:        "quiz_completion_rate",
			DisplayName: "Quiz Completion Rate",
			Type:        "percentage",
			SQL:         "ROUND(COUNT(DISTINCT CASE WHEN ase.is_correct IS NOT NULL THEN ase.student_id END) * 100.0 / GREATEST(COUNT(DISTINCT s.student_id), 1), 2)",
			Format:      "percentage",
		},
		// TIME-BASED MEASURES
		"response_speed_score": {
			Name:        "response_speed_score",
			DisplayName: "Response Speed Score",
			Type:        "calculated",
			SQL:         "ROUND(AVG(EXTRACT(EPOCH FROM (ase.submitted_at - qpe.published_at))), 2)",
			Format:      "seconds",
		},
	},
	"dimensions": map[string]Dimension{
		"session_id": {
			Name:        "session_id",
			DisplayName: "Quiz Session",
			Type:        "string",
			SQL:         "qs.session_id",
		},
		"classroom_name": {
			Name:        "classroom_name",
			DisplayName: "Classroom",
			Type:        "string",
			SQL:         "c.name",
		},
		"student_name": {
			Name:        "student_name",
			DisplayName: "Student",
			Type:        "string",
			SQL:         "s.name",
		},
		"question_id": {
			Name:        "question_id",
			DisplayName: "Question",
			Type:        "string",
			SQL:         "ase.question_id",
		},
		"answer_option": {
			Name:        "answer_option",
			DisplayName: "Answer Choice",
			Type:        "string",
			SQL:         "ase.answer",
		},
		"event_date": {
			Name:        "event_date",
			DisplayName: "Date",
			Type:        "time",
			SQL:         "DATE(ase.submitted_at)",
			Format:      "YYYY-MM-DD",
		},
		"event_hour": {
			Name:        "event_hour",
			DisplayName: "Hour",
			Type:        "time",
			SQL:         "EXTRACT(hour FROM ase.submitted_at)",
			Format:      "HH",
		},
		// STUDENT PERFORMANCE DIMENSIONS
		"performance_level": {
			Name:        "performance_level",
			DisplayName: "Performance Level",
			Type:        "string",
			SQL:         "CASE WHEN ase.is_correct = true THEN 'Correct' ELSE 'Incorrect' END",
		},
		"response_speed_category": {
			Name:        "response_speed_category",
			DisplayName: "Response Speed Category",
			Type:        "string",
			SQL:         "CASE WHEN EXTRACT(EPOCH FROM (ase.submitted_at - qpe.published_at)) < 30 THEN 'Fast' WHEN EXTRACT(EPOCH FROM (ase.submitted_at - qpe.published_at)) < 60 THEN 'Medium' ELSE 'Slow' END",
		},
		"correctness_flag": {
			Name:        "correctness_flag",
			DisplayName: "Answer Correctness",
			Type:        "boolean",
			SQL:         "ase.is_correct",
		},
		// ENGAGEMENT DIMENSIONS
		"engagement_level": {
			Name:        "engagement_level",
			DisplayName: "Engagement Level",
			Type:        "string",
			SQL:         "CASE WHEN ase.student_id IS NOT NULL THEN 'Active' ELSE 'Inactive' END",
		},
		"session_duration_category": {
			Name:        "session_duration_category",
			DisplayName: "Session Length",
			Type:        "string",
			SQL:         "CASE WHEN EXTRACT(EPOCH FROM (qs.ended_at - qs.started_at)) < 1800 THEN 'Short' WHEN EXTRACT(EPOCH FROM (qs.ended_at - qs.started_at)) < 3600 THEN 'Medium' ELSE 'Long' END",
		},
		// CONTENT EFFECTIVENESS DIMENSIONS
		"quiz_title": {
			Name:        "quiz_title",
			DisplayName: "Quiz Name",
			Type:        "string",
			SQL:         "qs.session_id",
		},
		"difficulty_level": {
			Name:        "difficulty_level",
			DisplayName: "Difficulty Level",
			Type:        "string",
			SQL:         "CASE WHEN qpe.timer_duration_sec < 30 THEN 'Hard' WHEN qpe.timer_duration_sec < 60 THEN 'Medium' ELSE 'Easy' END",
		},
		"timer_duration_category": {
			Name:        "timer_duration_category",
			DisplayName: "Question Timer",
			Type:        "string",
			SQL:         "CASE WHEN qpe.timer_duration_sec < 30 THEN 'Fast' WHEN qpe.timer_duration_sec < 60 THEN 'Medium' ELSE 'Slow' END",
		},
		"teacher_id": {
			Name:        "teacher_id",
			DisplayName: "Teacher",
			Type:        "string",
			SQL:         "qpe.teacher_id",
		},
		// TEMPORAL DIMENSIONS
		"event_week": {
			Name:        "event_week",
			DisplayName: "Week",
			Type:        "time",
			SQL:         "EXTRACT(week FROM ase.submitted_at)",
			Format:      "WW",
		},
		"event_month": {
			Name:        "event_month",
			DisplayName: "Month",
			Type:        "time",
			SQL:         "EXTRACT(month FROM ase.submitted_at)",
			Format:      "MM",
		},
		"event_day_of_week": {
			Name:        "event_day_of_week",
			DisplayName: "Day of Week",
			Type:        "time",
			SQL:         "TO_CHAR(ase.submitted_at, 'Day')",
			Format:      "string",
		},
		"time_bucket": {
			Name:        "time_bucket",
			DisplayName: "Time Bucket",
			Type:        "time",
			SQL:         "CASE WHEN EXTRACT(hour FROM ase.submitted_at) < 12 THEN 'Morning' WHEN EXTRACT(hour FROM ase.submitted_at) < 18 THEN 'Afternoon' ELSE 'Evening' END",
		},
	},
}

// QueryBuilder generates SQL from generic query request
func (qr *QueryRequest) BuildSQL() (string, error) {
	measures := QuizAnalyticsCube["measures"].(map[string]Measure)
	dimensions := QuizAnalyticsCube["dimensions"].(map[string]Dimension)

	var selectFields []string
	var groupByFields []string

	// Add requested measures
	for _, measureName := range qr.Measures {
		if measure, exists := measures[measureName]; exists {
			selectFields = append(selectFields, fmt.Sprintf("%s as %s", measure.SQL, measure.Name))
		}
	}

	// Add requested dimensions
	for _, dimName := range qr.Dimensions {
		if dim, exists := dimensions[dimName]; exists {
			selectFields = append(selectFields, fmt.Sprintf("%s as %s", dim.SQL, dim.Name))
			groupByFields = append(groupByFields, dim.SQL)
		}
	}

	if len(selectFields) == 0 {
		return "", fmt.Errorf("no valid measures or dimensions specified")
	}

	// Build base query
	query := fmt.Sprintf(`
		SELECT %s
		FROM answer_submitted_events ase
		LEFT JOIN quiz_sessions qs ON ase.session_id = qs.session_id
		LEFT JOIN classrooms c ON qs.classroom_id = c.classroom_id
		LEFT JOIN students s ON ase.student_id = s.student_id
		LEFT JOIN question_published_events qpe ON ase.question_id = qpe.question_id AND ase.session_id = qpe.session_id
	`, strings.Join(selectFields, ", "))

	// Add filters
	var whereConditions []string
	if qr.TimeRange != nil {
		whereConditions = append(whereConditions,
			fmt.Sprintf("ase.submitted_at BETWEEN '%s' AND '%s'",
				qr.TimeRange.Start.Format("2006-01-02 15:04:05"),
				qr.TimeRange.End.Format("2006-01-02 15:04:05")))
	}

	for field, value := range qr.Filters {
		if dim, exists := dimensions[field]; exists {
			whereConditions = append(whereConditions, fmt.Sprintf("%s = '%s'", dim.SQL, value))
		}
	}

	if len(whereConditions) > 0 {
		query += " WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Add GROUP BY
	if len(groupByFields) > 0 {
		query += " GROUP BY " + strings.Join(groupByFields, ", ")
	}

	// Add ORDER BY
	if len(qr.OrderBy) > 0 {
		var orderFields []string
		for _, order := range qr.OrderBy {
			orderFields = append(orderFields, fmt.Sprintf("%s %s", order.Field, order.Order))
		}
		query += " ORDER BY " + strings.Join(orderFields, ", ")
	}

	// Add LIMIT
	if qr.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", qr.Limit)
	}

	return query, nil
}
