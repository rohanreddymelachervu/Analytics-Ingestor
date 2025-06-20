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
		"response_time_avg": {
			Name:        "response_time_avg",
			DisplayName: "Average Response Time",
			Type:        "avg",
			SQL:         "AVG(ase.response_time_ms)",
			Format:      "milliseconds",
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
			SQL:         "q.question_id",
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
		LEFT JOIN questions q ON ase.question_id = q.question_id
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
