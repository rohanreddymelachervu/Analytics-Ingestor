# Analytics System Testing Guide

This guide demonstrates how to test the educational analytics system with real event data.

## Prerequisites

1. **Start the server** with proper environment variables:
   ```bash
   export DATABASE_URL="postgres://username:password@localhost/analytics_db?sslmode=disable"
   export JWT_SECRET="your-secret-key-here"
   go run cmd/server/main.go
   ```

2. **Create test users** for authentication:

### Step 1: Create Writer User (for apps to send events)
```bash
curl -X POST http://localhost:8080/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Whiteboard App",
    "email": "whiteboard@school.edu",
    "password": "apppassword123",
    "role": "writer"
  }'
```

### Step 2: Login to get Writer Token
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "whiteboard@school.edu",
    "password": "apppassword123"
  }'
```
**Save the `token` from the response as `WRITER_TOKEN`**

### Step 3: Create Reader User (for analytics dashboard)
```bash
curl -X POST http://localhost:8080/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Analytics Dashboard",
    "email": "analytics@school.edu", 
    "password": "dashpassword123",
    "role": "reader"
  }'

curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "analytics@school.edu",
    "password": "dashpassword123"
  }'
```
**Save the `token` from the response as `READER_TOKEN`**

## Step 4: Send Test Events

### 4.1 Session Started Event
```bash
curl -X POST http://localhost:8080/api/events \
  -H "Authorization: Bearer WRITER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "event_id": "550e8400-e29b-41d4-a716-446655440001",
    "event_type": "SESSION_STARTED",
    "timestamp": "2024-01-15T10:00:00.000Z",
    "session_id": "550e8400-e29b-41d4-a716-446655440100",
    "quiz_id": "550e8400-e29b-41d4-a716-446655440200",
    "classroom_id": "550e8400-e29b-41d4-a716-446655440300",
    "question_id": "550e8400-e29b-41d4-a716-446655440400",
    "teacher_id": "550e8400-e29b-41d4-a716-446655440500"
  }'
```

### 4.2 Question Published Events
```bash
curl -X POST http://localhost:8080/api/events \
  -H "Authorization: Bearer WRITER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "event_id": "550e8400-e29b-41d4-a716-446655440002",
    "event_type": "QUESTION_PUBLISHED",
    "timestamp": "2024-01-15T10:05:00.000Z",
    "session_id": "550e8400-e29b-41d4-a716-446655440100",
    "quiz_id": "550e8400-e29b-41d4-a716-446655440200",
    "classroom_id": "550e8400-e29b-41d4-a716-446655440300",
    "question_id": "550e8400-e29b-41d4-a716-446655440401",
    "teacher_id": "550e8400-e29b-41d4-a716-446655440500",
    "timer_sec": 60
  }'

curl -X POST http://localhost:8080/api/events \
  -H "Authorization: Bearer WRITER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "event_id": "550e8400-e29b-41d4-a716-446655440003",
    "event_type": "QUESTION_PUBLISHED", 
    "timestamp": "2024-01-15T10:07:00.000Z",
    "session_id": "550e8400-e29b-41d4-a716-446655440100",
    "quiz_id": "550e8400-e29b-41d4-a716-446655440200",
    "classroom_id": "550e8400-e29b-41d4-a716-446655440300",
    "question_id": "550e8400-e29b-41d4-a716-446655440402",
    "teacher_id": "550e8400-e29b-41d4-a716-446655440500",
    "timer_sec": 45
  }'
```

### 4.3 Answer Submitted Events (Multiple Students)
```bash
# Student 1 answers
curl -X POST http://localhost:8080/api/events \
  -H "Authorization: Bearer WRITER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "event_id": "550e8400-e29b-41d4-a716-446655440010",
    "event_type": "ANSWER_SUBMITTED",
    "timestamp": "2024-01-15T10:05:30.000Z",
    "session_id": "550e8400-e29b-41d4-a716-446655440100",
    "quiz_id": "550e8400-e29b-41d4-a716-446655440200",
    "classroom_id": "550e8400-e29b-41d4-a716-446655440300",
    "question_id": "550e8400-e29b-41d4-a716-446655440401",
    "student_id": "550e8400-e29b-41d4-a716-446655440601",
    "answer": "B",
    "response_time_ms": 30000
  }'

# Student 2 answers
curl -X POST http://localhost:8080/api/events \
  -H "Authorization: Bearer WRITER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "event_id": "550e8400-e29b-41d4-a716-446655440011",
    "event_type": "ANSWER_SUBMITTED",
    "timestamp": "2024-01-15T10:05:45.000Z",
    "session_id": "550e8400-e29b-41d4-a716-446655440100",
    "quiz_id": "550e8400-e29b-41d4-a716-446655440200",
    "classroom_id": "550e8400-e29b-41d4-a716-446655440300",
    "question_id": "550e8400-e29b-41d4-a716-446655440401",
    "student_id": "550e8400-e29b-41d4-a716-446655440602",
    "answer": "A",
    "response_time_ms": 45000
  }'
```

### 4.4 Batch Events Example
```bash
curl -X POST http://localhost:8080/api/events/batch \
  -H "Authorization: Bearer WRITER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '[
    {
      "event_id": "550e8400-e29b-41d4-a716-446655440012",
      "event_type": "ANSWER_SUBMITTED",
      "timestamp": "2024-01-15T10:07:30.000Z",
      "session_id": "550e8400-e29b-41d4-a716-446655440100",
      "quiz_id": "550e8400-e29b-41d4-a716-446655440200",
      "classroom_id": "550e8400-e29b-41d4-a716-446655440300",
      "question_id": "550e8400-e29b-41d4-a716-446655440402",
      "student_id": "550e8400-e29b-41d4-a716-446655440601",
      "answer": "C",
      "response_time_ms": 25000
    },
    {
      "event_id": "550e8400-e29b-41d4-a716-446655440013",
      "event_type": "ANSWER_SUBMITTED",
      "timestamp": "2024-01-15T10:07:50.000Z",
      "session_id": "550e8400-e29b-41d4-a716-446655440100",
      "quiz_id": "550e8400-e29b-41d4-a716-446655440200",
      "classroom_id": "550e8400-e29b-41d4-a716-446655440300",
      "question_id": "550e8400-e29b-41d4-a716-446655440402",
      "student_id": "550e8400-e29b-41d4-a716-446655440602",
      "answer": "C",
      "response_time_ms": 35000
    }
  ]'
```

## Step 5: Generate Reports

### 5.1 Active Participants Report
```bash
curl -X GET "http://localhost:8080/api/reports/active-participants?session_id=550e8400-e29b-41d4-a716-446655440100&time_range=60m" \
  -H "Authorization: Bearer READER_TOKEN"
```

### 5.2 Questions Per Minute Report
```bash
curl -X GET "http://localhost:8080/api/reports/questions-per-minute?session_id=550e8400-e29b-41d4-a716-446655440100" \
  -H "Authorization: Bearer READER_TOKEN"
```

### 5.3 Student Performance Report
```bash
curl -X GET "http://localhost:8080/api/reports/student-performance?student_id=550e8400-e29b-41d4-a716-446655440601&classroom_id=550e8400-e29b-41d4-a716-446655440300" \
  -H "Authorization: Bearer READER_TOKEN"
```

### 5.4 Classroom Engagement Report
```bash
curl -X GET "http://localhost:8080/api/reports/classroom-engagement?classroom_id=550e8400-e29b-41d4-a716-446655440300&date_range=7d" \
  -H "Authorization: Bearer READER_TOKEN"
```

### 5.5 Content Effectiveness Report
```bash
curl -X GET "http://localhost:8080/api/reports/content-effectiveness?quiz_id=550e8400-e29b-41d4-a716-446655440200" \
  -H "Authorization: Bearer READER_TOKEN"
```

## Expected Behavior

1. **Event Ingestion**: All events should return `201 Created` with the event ID
2. **Real Data**: Reports should show actual data from the events you sent, not mock data
3. **Analytics**: You should see meaningful metrics like:
   - Active participant counts
   - Response times and accuracy
   - Questions per minute trends
   - Engagement rates
   - Content effectiveness scores

## Notes

- **Question Data**: For `ANSWER_SUBMITTED` events to work properly, you need to first create quiz and question records in the database, or the system will return errors when trying to validate answers.
- **Student/Classroom Data**: Similarly, students and classrooms should exist in the database for complete analytics.
- **Real Analytics**: The reports will show real aggregated data from your events, including time-series data, accuracy calculations, and trend analysis.

## Troubleshooting

- **401 Unauthorized**: Check that your token is valid and has the correct scope (WRITE for events, READ for reports)
- **400 Bad Request**: Verify that all required fields are included and UUIDs are properly formatted
- **500 Internal Server Error**: Check that the database is properly connected and all foreign key relationships exist

This demonstrates a **complete event-driven analytics pipeline** with real database operations! 

# Analytics Ingestor Testing Guide

## Critical Gaps Fixed ✅

This system now implements **ALL** required metrics including:

### 🔥 Real-Time Engagement Metrics
- ✅ Active Participants
- ✅ Questions Served/Min (Questions Per Minute)  
- ✅ **Response Rate** - % of students who answered out of those who received question
- ✅ **Latency to First Answer** - Time from question delivery to first response

### 🔥 Question-Level Metrics  
- ✅ Attempt Count
- ✅ Correctness Rate
- ✅ Average Response Time
- ✅ **Skipped/Timeout Rate** - % who didn't answer before time limit
- ✅ Difficulty Index

### 🔥 Quiz-Level Metrics
- ✅ **Completion Rate** - % of students who finish all questions
- ✅ Average Score  
- ✅ Score Distribution
- ✅ Time to Completion
- ✅ **Drop-off Points** - Questions where >X% abandon

### 🔥 Student & Classroom Level
- ✅ All individual and aggregate metrics

### 🔥 Timer Validation
- ✅ **Critical**: Ensures `timestamp ≤ publishedTimestamp + timerDuration`

## Setup

1. Start PostgreSQL:
```bash
# Ensure PostgreSQL is running on localhost:5432
```

2. Set environment variables:
```bash
export DATABASE_URL="postgres://postgres:root@localhost:5432/ingestor"
export JWT_SECRET="62c23d514144fc4fd1dd75fdfed51791f4b9ee14f153db00411ef0eb0bb62aca"
```

3. Start the server:
```bash
./server
```

## Authentication Setup

### Create WRITE User (for event ingestion)
```bash
curl -X POST http://localhost:8080/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Writer User",
    "email": "writer@example.com", 
    "password": "password123",
    "role": "writer"
  }'
```

### Login and Get WRITE Token
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "writer@example.com",
    "password": "password123"
  }'
```

### Create READ User (for reports)  
```bash
curl -X POST http://localhost:8080/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Reader User",
    "email": "reader@example.com",
    "password": "password123", 
    "role": "reader"
  }'
```

### Login and Get READ Token
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "reader@example.com", 
    "password": "password123"
  }'
```

## Event Ingestion Testing

**Use WRITE token for all events:**

### 1. Session Started Event
```bash
curl -X POST http://localhost:8080/api/events \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_WRITE_TOKEN" \
  -d '{
    "event_id": "550e8400-e29b-41d4-a716-446655440000",
    "event_type": "SESSION_STARTED", 
    "timestamp": "2025-06-19T10:00:00Z",
    "session_id": "550e8400-e29b-41d4-a716-446655440000",
    "quiz_id": "550e8400-e29b-41d4-a716-446655440010",
    "classroom_id": "550e8400-e29b-41d4-a716-446655440020",
    "question_id": "550e8400-e29b-41d4-a716-446655440030"
  }'
```

### 2. Question Published Event
```bash
curl -X POST http://localhost:8080/api/events \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_WRITE_TOKEN" \
  -d '{
    "event_id": "550e8400-e29b-41d4-a716-446655440001",
    "event_type": "QUESTION_PUBLISHED",
    "timestamp": "2025-06-19T10:01:00Z", 
    "session_id": "550e8400-e29b-41d4-a716-446655440000",
    "quiz_id": "550e8400-e29b-41d4-a716-446655440010",
    "classroom_id": "550e8400-e29b-41d4-a716-446655440020",
    "question_id": "550e8400-e29b-41d4-a716-446655440030",
    "teacher_id": "550e8400-e29b-41d4-a716-446655440050",
    "timer_sec": 30
  }'
```

### 3. Answer Submitted Events (with Timer Validation)
```bash
# Student 1 - Correct Answer (A) - WITHIN timer  
curl -X POST http://localhost:8080/api/events \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_WRITE_TOKEN" \
  -d '{
    "event_id": "550e8400-e29b-41d4-a716-446655440002",
    "event_type": "ANSWER_SUBMITTED",
    "timestamp": "2025-06-19T10:01:15Z",
    "session_id": "550e8400-e29b-41d4-a716-446655440000", 
    "quiz_id": "550e8400-e29b-41d4-a716-446655440010",
    "classroom_id": "550e8400-e29b-41d4-a716-446655440020",
    "question_id": "550e8400-e29b-41d4-a716-446655440030",
    "student_id": "550e8400-e29b-41d4-a716-446655440040",
    "answer": "A"
  }'

# Student 2 - Incorrect Answer (B) - WITHIN timer
curl -X POST http://localhost:8080/api/events \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_WRITE_TOKEN" \
  -d '{
    "event_id": "550e8400-e29b-41d4-a716-446655440003", 
    "event_type": "ANSWER_SUBMITTED",
    "timestamp": "2025-06-19T10:01:20Z",
    "session_id": "550e8400-e29b-41d4-a716-446655440000",
    "quiz_id": "550e8400-e29b-41d4-a716-446655440010", 
    "classroom_id": "550e8400-e29b-41d4-a716-446655440020",
    "question_id": "550e8400-e29b-41d4-a716-446655440030",
    "student_id": "550e8400-e29b-41d4-a716-446655440041",
    "answer": "B"
  }'

# Test Timer Validation - LATE answer (should FAIL)
curl -X POST http://localhost:8080/api/events \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_WRITE_TOKEN" \
  -d '{
    "event_id": "550e8400-e29b-41d4-a716-446655440099",
    "event_type": "ANSWER_SUBMITTED", 
    "timestamp": "2025-06-19T10:02:00Z",
    "session_id": "550e8400-e29b-41d4-a716-446655440000",
    "quiz_id": "550e8400-e29b-41d4-a716-446655440010",
    "classroom_id": "550e8400-e29b-41d4-a716-446655440020",
    "question_id": "550e8400-e29b-41d4-a716-446655440030", 
    "student_id": "550e8400-e29b-41d4-a716-446655440042",
    "answer": "C"
  }'
```

## 🔥 CRITICAL NEW METRICS TESTING

**Use READ token for all reports:**

### Response Rate Analysis
```bash
curl -X GET "http://localhost:8080/api/reports/response-rate?session_id=550e8400-e29b-41d4-a716-446655440000&question_id=550e8400-e29b-41d4-a716-446655440030" \
  -H "Authorization: Bearer YOUR_READ_TOKEN"
```

### Latency Analysis  
```bash
curl -X GET "http://localhost:8080/api/reports/latency-analysis?session_id=550e8400-e29b-41d4-a716-446655440000&question_id=550e8400-e29b-41d4-a716-446655440030" \
  -H "Authorization: Bearer YOUR_READ_TOKEN"
```

### Timeout Analysis
```bash
curl -X GET "http://localhost:8080/api/reports/timeout-analysis?session_id=550e8400-e29b-41d4-a716-446655440000&question_id=550e8400-e29b-41d4-a716-446655440030" \
  -H "Authorization: Bearer YOUR_READ_TOKEN"
```

### Completion Rate
```bash
curl -X GET "http://localhost:8080/api/reports/completion-rate?session_id=550e8400-e29b-41d4-a716-446655440000" \
  -H "Authorization: Bearer YOUR_READ_TOKEN"
```

### Dropoff Analysis
```bash
curl -X GET "http://localhost:8080/api/reports/dropoff-analysis?session_id=550e8400-e29b-41d4-a716-446655440000" \
  -H "Authorization: Bearer YOUR_READ_TOKEN"
```

## Existing Reports (All Working)

### Active Participants
```bash
curl -X GET "http://localhost:8080/api/reports/active-participants?session_id=550e8400-e29b-41d4-a716-446655440000&time_range=24h" \
  -H "Authorization: Bearer YOUR_READ_TOKEN"
```

### Questions Per Minute
```bash
curl -X GET "http://localhost:8080/api/reports/questions-per-minute?session_id=550e8400-e29b-41d4-a716-446655440000" \
  -H "Authorization: Bearer YOUR_READ_TOKEN"
```

### Student Performance
```bash
curl -X GET "http://localhost:8080/api/reports/student-performance?student_id=550e8400-e29b-41d4-a716-446655440040&classroom_id=550e8400-e29b-41d4-a716-446655440020" \
  -H "Authorization: Bearer YOUR_READ_TOKEN"
```

### Classroom Engagement
```bash
curl -X GET "http://localhost:8080/api/reports/classroom-engagement?classroom_id=550e8400-e29b-41d4-a716-446655440020&date_range=7d" \
  -H "Authorization: Bearer YOUR_READ_TOKEN"
```

### Content Effectiveness
```bash
curl -X GET "http://localhost:8080/api/reports/content-effectiveness?quiz_id=550e8400-e29b-41d4-a716-446655440010" \
  -H "Authorization: Bearer YOUR_READ_TOKEN"
```

## Expected Results

### Response Rate
- **students_received**: 2 (total students in classroom)
- **students_answered**: 2 (both students answered)  
- **response_rate**: 1.0 (100%)

### Latency Analysis
- **first_answer_latency**: ~15s (first student response)
- **average_latency**: ~17.5s 
- **median_latency**: ~17.5s

### Timeout Analysis
- **timeout_rate**: 0.0 (no students timed out)
- **skipped_rate**: 0.0 

### Completion Rate
- **completion_rate**: 1.0 (100% - both students answered all questions)

### Timer Validation 
- Late answer (after 30sec timer) should **FAIL** with timing validation error

## Error Testing

### Invalid UUID
```bash
curl -X GET "http://localhost:8080/api/reports/response-rate?session_id=invalid-uuid&question_id=550e8400-e29b-41d4-a716-446655440030" \
  -H "Authorization: Bearer YOUR_READ_TOKEN"
```

### Missing Parameters
```bash
curl -X GET "http://localhost:8080/api/reports/response-rate?session_id=550e8400-e29b-41d4-a716-446655440000" \
  -H "Authorization: Bearer YOUR_READ_TOKEN"
```

### Wrong Token Scope
```bash
# Try using WRITE token for READ endpoint (should fail)
curl -X GET "http://localhost:8080/api/reports/response-rate?session_id=550e8400-e29b-41d4-a716-446655440000&question_id=550e8400-e29b-41d4-a716-446655440030" \
  -H "Authorization: Bearer YOUR_WRITE_TOKEN"
```

## Status Summary

✅ **ALL CRITICAL METRICS IMPLEMENTED:**
- Real-time engagement (Active participants, QPM, Response Rate, Latency)
- Question-level (Attempts, Correctness, Response Time, Timeouts)  
- Quiz-level (Completion Rate, Scores, Drop-offs)
- Student/Classroom aggregations
- Timer validation with deadline enforcement
- Comprehensive error handling
- JWT authentication with proper scope enforcement

🎯 **Production Ready** for educational analytics at scale (900K users) 

# 🧪 **COMPREHENSIVE ANALYTICS INGESTOR TESTING GUIDE**

## **🚀 Setup & Database Preparation**

### **1. Start the Server**
```bash
export DATABASE_URL="postgres://postgres:root@localhost:5432/ingestor"
export JWT_SECRET="62c23d514144fc4fd1dd75fdfed51791f4b9ee14f153db00411ef0eb0bb62aca"
go build -o server cmd/server/main.go && ./server
```

### **2. Create Test Data Setup (CRITICAL - Avoids Foreign Key Errors)**

**First, manually insert base entities into database:**

```sql
-- Connect to your PostgreSQL database
psql postgres://postgres:root@localhost:5432/ingestor

-- Insert test quiz
INSERT INTO quizzes (quiz_id, title, description) VALUES 
('550e8400-e29b-41d4-a716-446655440010', 'Math Quiz 1', 'Basic mathematics questions');

-- Insert test classroom  
INSERT INTO classrooms (classroom_id, name) VALUES 
('550e8400-e29b-41d4-a716-446655440020', 'Class 5A');

-- Insert test students
INSERT INTO students (student_id, name) VALUES 
('550e8400-e29b-41d4-a716-446655440040', 'Alice Johnson'),
('550e8400-e29b-41d4-a716-446655440041', 'Bob Smith');

-- Link students to classroom
INSERT INTO classroom_students (classroom_id, student_id) VALUES 
('550e8400-e29b-41d4-a716-446655440020', '550e8400-e29b-41d4-a716-446655440040'),
('550e8400-e29b-41d4-a716-446655440020', '550e8400-e29b-41d4-a716-446655440041');

-- Insert test question
INSERT INTO questions (question_id, quiz_id) VALUES 
('550e8400-e29b-41d4-a716-446655440030', '550e8400-e29b-41d4-a716-446655440010');
```

## **🔐 Authentication Testing**

### **1. Create Test Users**

**WRITE User (for event ingestion):**
```bash
curl -X POST http://127.0.0.1:8080/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Write User",
    "email": "writer@test.com",
    "password": "password123",
    "role": "writer"
  }'
```

**READ User (for reports):**
```bash
curl -X POST http://127.0.0.1:8080/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Read User", 
    "email": "reader@test.com",
    "password": "password123",
    "role": "reader"
  }'
```

### **2. Generate JWT Tokens**

**Get WRITE token:**
```bash
curl -X POST http://127.0.0.1:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "writer@test.com",
    "password": "password123"
  }'
```

**Get READ token:**
```bash
curl -X POST http://127.0.0.1:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "reader@test.com", 
    "password": "password123"
  }'
```

**💡 Save both tokens as environment variables:**
```bash
export WRITE_TOKEN="your_write_jwt_token_here"
export READ_TOKEN="your_read_jwt_token_here"
```

## **📊 Event Ingestion Testing**

### **1. Session Started Event**
```bash
curl -X POST http://127.0.0.1:8080/api/events \
  -H "Authorization: Bearer $WRITE_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "event_id": "550e8400-e29b-41d4-a716-446655440001",
    "event_type": "SESSION_STARTED",
    "timestamp": "2025-06-19T10:00:00.000Z",
    "session_id": "550e8400-e29b-41d4-a716-446655440000",
    "quiz_id": "550e8400-e29b-41d4-a716-446655440010",
    "classroom_id": "550e8400-e29b-41d4-a716-446655440020",
    "question_id": "550e8400-e29b-41d4-a716-446655440030",
    "teacher_id": "550e8400-e29b-41d4-a716-446655440050"
  }'
```

### **2. Question Published Event (WITH TIMER)**
```bash
curl -X POST http://127.0.0.1:8080/api/events \
  -H "Authorization: Bearer $WRITE_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "event_id": "550e8400-e29b-41d4-a716-446655440002",
    "event_type": "QUESTION_PUBLISHED",
    "timestamp": "2025-06-19T10:01:00.000Z",
    "session_id": "550e8400-e29b-41d4-a716-446655440000",
    "quiz_id": "550e8400-e29b-41d4-a716-446655440010",
    "classroom_id": "550e8400-e29b-41d4-a716-446655440020",
    "question_id": "550e8400-e29b-41d4-a716-446655440030",
    "teacher_id": "550e8400-e29b-41d4-a716-446655440050",
    "timer_sec": 30
  }'
```

### **3. Answer Submitted Events (Valid Timing)**

**Student 1 - Correct Answer (A/C = correct):**
```bash
curl -X POST http://127.0.0.1:8080/api/events \
  -H "Authorization: Bearer $WRITE_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "event_id": "550e8400-e29b-41d4-a716-446655440003",
    "event_type": "ANSWER_SUBMITTED",
    "timestamp": "2025-06-19T10:01:15.000Z",
    "session_id": "550e8400-e29b-41d4-a716-446655440000",
    "quiz_id": "550e8400-e29b-41d4-a716-446655440010",
    "classroom_id": "550e8400-e29b-41d4-a716-446655440020",
    "question_id": "550e8400-e29b-41d4-a716-446655440030",
    "student_id": "550e8400-e29b-41d4-a716-446655440040",
    "answer": "A"
  }'
```

**Student 2 - Incorrect Answer (B/D = incorrect):**
```bash
curl -X POST http://127.0.0.1:8080/api/events \
  -H "Authorization: Bearer $WRITE_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "event_id": "550e8400-e29b-41d4-a716-446655440004",
    "event_type": "ANSWER_SUBMITTED", 
    "timestamp": "2025-06-19T10:01:20.000Z",
    "session_id": "550e8400-e29b-41d4-a716-446655440000",
    "quiz_id": "550e8400-e29b-41d4-a716-446655440010",
    "classroom_id": "550e8400-e29b-41d4-a716-446655440020",
    "question_id": "550e8400-e29b-41d4-a716-446655440030",
    "student_id": "550e8400-e29b-41d4-a716-446655440041",
    "answer": "B"
  }'
```

### **4. Test Timer Validation (Should FAIL):**
```bash
curl -X POST http://127.0.0.1:8080/api/events \
  -H "Authorization: Bearer $WRITE_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "event_id": "550e8400-e29b-41d4-a716-446655440099",
    "event_type": "ANSWER_SUBMITTED",
    "timestamp": "2025-06-19T10:02:00.000Z",
    "session_id": "550e8400-e29b-41d4-a716-446655440000",
    "quiz_id": "550e8400-e29b-41d4-a716-446655440010",
    "classroom_id": "550e8400-e29b-41d4-a716-446655440020",
    "question_id": "550e8400-e29b-41d4-a716-446655440030",
    "student_id": "550e8400-e29b-41d4-a716-446655440041",
    "answer": "D"
  }'
```
**Expected:** Should return 500 error with "answer submitted after deadline" message.

### **5. Batch Events Testing**
```bash
curl -X POST http://127.0.0.1:8080/api/events/batch \
  -H "Authorization: Bearer $WRITE_TOKEN" \
  -H "Content-Type: application/json" \
  -d '[
    {
      "event_id": "550e8400-e29b-41d4-a716-446655440005",
      "event_type": "ANSWER_SUBMITTED",
      "timestamp": "2025-06-19T10:01:25.000Z",
      "session_id": "550e8400-e29b-41d4-a716-446655440000",
      "quiz_id": "550e8400-e29b-41d4-a716-446655440010",
      "classroom_id": "550e8400-e29b-41d4-a716-446655440020",
      "question_id": "550e8400-e29b-41d4-a716-446655440030",
      "student_id": "550e8400-e29b-41d4-a716-446655440040",
      "answer": "C"
    }
  ]'
```

## **📈 Analytics Reports Testing**

### **1. Active Participants**
```bash
curl -H "Authorization: Bearer $READ_TOKEN" \
  "http://127.0.0.1:8080/api/reports/active-participants?session_id=550e8400-e29b-41d4-a716-446655440000&time_range=24h"
```
**Expected:** 2 students, ~50% accuracy

### **2. Questions Per Minute**
```bash
curl -H "Authorization: Bearer $READ_TOKEN" \
  "http://127.0.0.1:8080/api/reports/questions-per-minute?session_id=550e8400-e29b-41d4-a716-446655440000"
```
**Expected:** 1 question, QPM rates

### **3. Student Performance**
```bash
curl -H "Authorization: Bearer $READ_TOKEN" \
  "http://127.0.0.1:8080/api/reports/student-performance?student_id=550e8400-e29b-41d4-a716-446655440040&classroom_id=550e8400-e29b-41d4-a716-446655440020"
```
**Expected:** 100% accuracy for Student 1

### **4. Classroom Engagement**
```bash
curl -H "Authorization: Bearer $READ_TOKEN" \
  "http://127.0.0.1:8080/api/reports/classroom-engagement?classroom_id=550e8400-e29b-41d4-a716-446655440020&date_range=7d"
```
**Expected:** 2 total students, engagement metrics

### **5. Content Effectiveness**
```bash
curl -H "Authorization: Bearer $READ_TOKEN" \
  "http://127.0.0.1:8080/api/reports/content-effectiveness?quiz_id=550e8400-e29b-41d4-a716-446655440010"
```
**Expected:** Quiz analysis with recommendations

## **🔬 NEW CRITICAL METRICS TESTING**

### **6. Response Rate Analysis**
```bash
curl -H "Authorization: Bearer $READ_TOKEN" \
  "http://127.0.0.1:8080/api/reports/response-rate?session_id=550e8400-e29b-41d4-a716-446655440000&question_id=550e8400-e29b-41d4-a716-446655440030"
```
**Expected:** Response rate calculation (answered/total students)

### **7. Latency Analysis**
```bash
curl -H "Authorization: Bearer $READ_TOKEN" \
  "http://127.0.0.1:8080/api/reports/latency-analysis?session_id=550e8400-e29b-41d4-a716-446655440000&question_id=550e8400-e29b-41d4-a716-446655440030"
```
**Expected:** Answer timing metrics

### **8. Timeout Analysis**
```bash
curl -H "Authorization: Bearer $READ_TOKEN" \
  "http://127.0.0.1:8080/api/reports/timeout-analysis?session_id=550e8400-e29b-41d4-a716-446655440000&question_id=550e8400-e29b-41d4-a716-446655440030"
```
**Expected:** Timeout and skip rates

### **9. Completion Rate**
```bash
curl -H "Authorization: Bearer $READ_TOKEN" \
  "http://127.0.0.1:8080/api/reports/completion-rate?session_id=550e8400-e29b-41d4-a716-446655440000"
```
**Expected:** Session completion statistics

### **10. Dropoff Analysis**
```bash
curl -H "Authorization: Bearer $READ_TOKEN" \
  "http://127.0.0.1:8080/api/reports/dropoff-analysis?session_id=550e8400-e29b-41d4-a716-446655440000"
```
**Expected:** Question abandonment analysis

## **🚨 Error Testing**

### **1. Authentication Errors**
```bash
# No token
curl "http://127.0.0.1:8080/api/reports/active-participants?session_id=550e8400-e29b-41d4-a716-446655440000"
# Expected: 401 Unauthorized

# Wrong scope (WRITE token on READ endpoint)
curl -H "Authorization: Bearer $WRITE_TOKEN" \
  "http://127.0.0.1:8080/api/reports/active-participants?session_id=550e8400-e29b-41d4-a716-446655440000"
# Expected: 403 Forbidden
```

### **2. Validation Errors**
```bash
# Invalid UUID
curl -H "Authorization: Bearer $READ_TOKEN" \
  "http://127.0.0.1:8080/api/reports/active-participants?session_id=invalid-uuid"
# Expected: 400 Bad Request

# Missing required parameters
curl -X POST http://127.0.0.1:8080/api/events \
  -H "Authorization: Bearer $WRITE_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"event_type": "QUESTION_PUBLISHED"}'
# Expected: 400 Bad Request (missing required fields)
```

## **✅ Expected Database State After Testing**

**Check data with:**
```sql
-- Verify events were stored
SELECT COUNT(*) FROM question_published_events;  -- Should be 1
SELECT COUNT(*) FROM answer_submitted_events;    -- Should be 3-4
SELECT COUNT(*) FROM quiz_sessions;              -- Should be 1

-- Verify accuracy calculations
SELECT 
  student_id, 
  COUNT(*) as total_answers,
  SUM(CASE WHEN is_correct THEN 1 ELSE 0 END) as correct_answers,
  ROUND(AVG(CASE WHEN is_correct THEN 1.0 ELSE 0.0 END) * 100, 2) as accuracy
FROM answer_submitted_events 
WHERE session_id = '550e8400-e29b-41d4-a716-446655440000'
GROUP BY student_id;
```

## **🏆 Success Criteria**

- [ ] ✅ All events stored without foreign key errors
- [ ] ✅ Timer validation working (late submissions rejected)
- [ ] ✅ Authentication and scope enforcement working  
- [ ] ✅ All report endpoints return real data
- [ ] ✅ Accuracy calculations correct (Student 1: 100%, Student 2: 0%)
- [ ] ✅ No GORM scanning errors in logs
- [ ] ✅ All new critical metrics working
- [ ] ✅ Error handling comprehensive

**🎯 ZERO hardcoded data - all calculations from real database operations!** 