# Educational Analytics Framework - Technical Design Document

## Executive Summary

This document outlines the design and implementation of a comprehensive reporting framework for our educational applications ecosystem, consisting of a Whiteboard App (teacher-facing) and Notebook App (student-facing) that synchronize during quiz sessions.

## 1. Data Collection Strategy

### 1.1 Metrics Tracked Across Applications

#### Whiteboard App (Teacher Events)
- **Quiz Session Management**
  - Session start/end timestamps
  - Quiz metadata (title, description, subject)
  - Classroom and participant counts

- **Question Publishing Events**
  - Question ID and content
  - Publish timestamp
  - Timer duration
  - Question difficulty level
  - Expected answer format

- **Real-time Control Events**
  - Question transitions
  - Timer modifications
  - Session pause/resume actions

#### Notebook App (Student Events)
- **Answer Submission Events**
  - Student ID and answer content
  - Submission timestamp
  - Response time (time from question publish to submission)
  - Answer correctness validation

- **Engagement Events**
  - App focus/unfocus events
  - Question view duration
  - Navigation patterns within the app

- **Performance Indicators**
  - Consecutive correct/incorrect streaks
  - Time-to-first-response patterns
  - Confidence levels (if collected)

### 1.2 Event-Based Tracking Methodology

```json
{
  "event_type": "answer_submitted",
  "timestamp": "2024-01-15T10:30:45.123Z",
  "session_id": "550e8400-e29b-41d4-a716-446655440000",
  "data": {
    "student_id": "550e8400-e29b-41d4-a716-446655440001",
    "question_id": "550e8400-e29b-41d4-a716-446655440002", 
    "answer": "B",
    "response_time_ms": 12300,
    "is_correct": true,
    "confidence_level": "high"
  }
}
```

**Event Processing Pipeline:**
1. **Collection**: Apps send events to `/api/events` endpoint
2. **Validation**: Event schema validation and authentication
3. **Enrichment**: Add derived metrics (correctness, timing calculations)
4. **Storage**: Persist to PostgreSQL with proper indexing
5. **Aggregation**: Real-time and batch processing for reports

## 2. Database Schema Design

### 2.1 Core Entities

```sql
-- Core Educational Entities
CREATE TABLE quizzes (
  quiz_id      UUID      PRIMARY KEY,
  title        VARCHAR   NOT NULL,
  description  TEXT,
  subject      VARCHAR(50),
  difficulty   VARCHAR(20) DEFAULT 'medium',
  created_at   TIMESTAMP DEFAULT NOW()
);

CREATE TABLE classrooms (
  classroom_id  UUID    PRIMARY KEY,
  name          VARCHAR NOT NULL,
  school_id     UUID,
  grade_level   VARCHAR(10),
  created_at    TIMESTAMP DEFAULT NOW()
);

CREATE TABLE students (
  student_id  UUID    PRIMARY KEY,
  name        VARCHAR NOT NULL,
  grade_level VARCHAR(10),
  created_at  TIMESTAMP DEFAULT NOW()
);

CREATE TABLE questions (
  question_id     UUID    PRIMARY KEY,
  quiz_id         UUID    NOT NULL,
  question_text   TEXT    NOT NULL,
  question_type   VARCHAR(20) DEFAULT 'multiple_choice',
  options         JSON,
  correct_answer  VARCHAR NOT NULL,
  difficulty      VARCHAR(20) DEFAULT 'medium',
  points          INTEGER DEFAULT 1,
  FOREIGN KEY (quiz_id) REFERENCES quizzes(quiz_id) ON DELETE CASCADE
);
```

### 2.2 Session and Event Tables

```sql
-- Quiz Session Management
CREATE TABLE quiz_sessions (
  session_id     UUID      PRIMARY KEY,
  quiz_id        UUID      NOT NULL,
  classroom_id   UUID      NOT NULL,
  teacher_id     UUID,
  started_at     TIMESTAMP NOT NULL,
  ended_at       TIMESTAMP,
  status         VARCHAR(20) DEFAULT 'active',
  total_questions INTEGER DEFAULT 0,
  FOREIGN KEY (quiz_id) REFERENCES quizzes(quiz_id) ON DELETE RESTRICT,
  FOREIGN KEY (classroom_id) REFERENCES classrooms(classroom_id) ON DELETE RESTRICT
);

-- Teacher Action Events
CREATE TABLE question_published_events (
  event_id           UUID      PRIMARY KEY,
  session_id         UUID      NOT NULL,
  question_id        UUID      NOT NULL,
  teacher_id         UUID,
  published_at       TIMESTAMP NOT NULL,
  timer_duration_sec INT       NOT NULL,
  sequence_number    INTEGER,
  FOREIGN KEY (session_id) REFERENCES quiz_sessions(session_id) ON DELETE CASCADE,
  FOREIGN KEY (question_id) REFERENCES questions(question_id) ON DELETE CASCADE
);

-- Student Response Events  
CREATE TABLE answer_submitted_events (
  event_id       UUID      PRIMARY KEY,
  session_id     UUID      NOT NULL,
  question_id    UUID      NOT NULL,
  student_id     UUID      NOT NULL,
  answer         VARCHAR   NOT NULL,
  is_correct     BOOLEAN   NOT NULL,
  submitted_at   TIMESTAMP NOT NULL,
  response_time_ms INTEGER,
  confidence_level VARCHAR(10),
  FOREIGN KEY (session_id) REFERENCES quiz_sessions(session_id) ON DELETE CASCADE,
  FOREIGN KEY (question_id) REFERENCES questions(question_id) ON DELETE CASCADE,
  FOREIGN KEY (student_id) REFERENCES students(student_id) ON DELETE CASCADE
);
```

### 2.3 Performance Optimization Indexes

```sql
-- Performance indexes for common query patterns
CREATE INDEX idx_ase_session_question ON answer_submitted_events (session_id, question_id);
CREATE INDEX idx_ase_session_student ON answer_submitted_events (session_id, student_id);
CREATE INDEX idx_ase_submitted_at ON answer_submitted_events (submitted_at);
CREATE INDEX idx_ase_correctness ON answer_submitted_events (is_correct, submitted_at);

CREATE INDEX idx_qpe_session_published_at ON question_published_events (session_id, published_at);
CREATE INDEX idx_sessions_classroom_time ON quiz_sessions (classroom_id, started_at);
```

## 3. API Design

### 3.1 Authentication & Authorization

**JWT-Based Authentication:**
- Role-based access control: `writer` (apps) vs `reader` (analytics dashboard)
- Scope-based permissions: `WRITE` (event ingestion) vs `READ` (report access)

```bash
# Authentication Flow
POST /api/auth/signup
POST /api/auth/login
```

### 3.2 Data Ingestion Endpoints

#### Single Event Ingestion
```bash
POST /api/events
Authorization: Bearer <writer_token>
Content-Type: application/json

{
  "event_type": "answer_submitted",
  "session_id": "550e8400-e29b-41d4-a716-446655440000",
  "data": {
    "student_id": "550e8400-e29b-41d4-a716-446655440001",
    "question_id": "550e8400-e29b-41d4-a716-446655440002",
    "answer": "B",
    "response_time_ms": 12300
  }
}
```

#### Batch Event Ingestion
```bash
POST /api/events/batch
Authorization: Bearer <writer_token>
Content-Type: application/json

[
  { "event_type": "question_published", "session_id": "...", "data": {...} },
  { "event_type": "answer_submitted", "session_id": "...", "data": {...} },
  { "event_type": "answer_submitted", "session_id": "...", "data": {...} }
]
```

### 3.3 Query Endpoints for Reports

#### Student Performance Analysis
```bash
GET /api/reports/student-performance?student_id=<uuid>&classroom_id=<uuid>
Authorization: Bearer <reader_token>

Response:
{
  "report_type": "student_performance",
  "student_id": "550e8400-e29b-41d4-a716-446655440001",
  "overall_accuracy": 0.78,
  "questions_attempted": 45,
  "correct_answers": 35,
  "average_response_time": "12.3s",
  "improvement_trend": "+5.2%",
  "subject_breakdown": [
    {"subject": "Mathematics", "accuracy": 0.85, "questions": 20},
    {"subject": "Science", "accuracy": 0.72, "questions": 15}
  ]
}
```

#### Classroom Engagement Metrics
```bash
GET /api/reports/classroom-engagement?classroom_id=<uuid>&date_range=7d
Authorization: Bearer <reader_token>

Response:
{
  "report_type": "classroom_engagement", 
  "classroom_id": "550e8400-e29b-41d4-a716-446655440003",
  "total_students": 28,
  "active_students": 25,
  "engagement_rate": 0.89,
  "average_session_duration": "42m",
  "participation_trend": [
    {"date": "2024-01-08", "participation": 0.82},
    {"date": "2024-01-09", "participation": 0.85}
  ]
}
```

#### Content Effectiveness Evaluation
```bash
GET /api/reports/content-effectiveness?quiz_id=<uuid>
Authorization: Bearer <reader_token>

Response:
{
  "report_type": "content_effectiveness",
  "quiz_id": "550e8400-e29b-41d4-a716-446655440004",
  "total_sessions": 12,
  "overall_accuracy": 0.74,
  "question_analysis": [
    {
      "question_id": "q1",
      "accuracy": 0.91,
      "avg_response_time": "8.2s", 
      "effectiveness_score": 0.95
    }
  ],
  "recommendations": [
    "Question q2 shows low accuracy - consider adding hints"
  ]
}
```

## 4. System Architecture & Scale Considerations

### 4.1 Scale Requirements
- **1,000 schools** × **30 classrooms** × **30 students** = **900,000 total users**
- **Peak concurrent users**: ~180,000 (20% activity rate)
- **Event volume**: ~50,000 events/minute during peak quiz times
- **Storage growth**: ~1TB/year for event data

### 4.2 Architecture Components

```
┌─────────────────┐    ┌─────────────────┐
│  Whiteboard App │    │   Notebook App  │
│   (Teachers)    │    │   (Students)    │
└─────────┬───────┘    └─────────┬───────┘
          │                      │
          └──────────┬───────────┘
                     │ HTTPS/JSON
                     ▼
          ┌─────────────────────────┐
          │    Load Balancer        │
          │   (nginx/cloudflare)    │
          └─────────┬───────────────┘
                    │
          ┌─────────▼───────────────┐
          │     Gin API Server      │
          │  (Go + Authentication)  │
          └─────────┬───────────────┘
                    │
          ┌─────────▼───────────────┐
          │   PostgreSQL Cluster   │
          │  (Primary + Replicas)   │
          └─────────────────────────┘
```

### 4.3 Performance Optimizations
1. **Database Connection Pooling**: Configure GORM with appropriate pool sizes
2. **Read Replicas**: Separate read queries (reports) from write queries (events)
3. **Caching Layer**: Redis for frequently accessed aggregated data
4. **Batch Processing**: Process events in batches for better throughput
5. **Asynchronous Processing**: Use message queues for non-critical processing

## 5. Security Considerations

### 5.1 Data Privacy
- **Student Data Protection**: Comply with FERPA/COPPA regulations
- **Anonymization**: Option to anonymize student identifiers in reports
- **Data Retention**: Configurable retention policies for event data

### 5.2 Access Control
- **Role-Based Permissions**: Teachers can only access their classroom data
- **API Rate Limiting**: Prevent abuse and ensure fair usage
- **Audit Logging**: Track all data access for compliance

## 6. Implementation Roadmap

### Phase 1: Core Infrastructure (Completed)
- ✅ Database schema and migrations
- ✅ JWT authentication system
- ✅ Basic API endpoints for event ingestion
- ✅ Mock reporting endpoints

### Phase 2: Production Readiness
- [ ] Database connection optimization
- [ ] Proper error handling and logging
- [ ] Input validation and sanitization  
- [ ] API documentation (Swagger)
- [ ] Unit and integration tests

### Phase 3: Advanced Analytics
- [ ] Real-time aggregation pipelines
- [ ] Advanced reporting queries
- [ ] Data visualization endpoints
- [ ] Export functionality (CSV, PDF)

### Phase 4: Scale & Performance
- [ ] Load testing and optimization
- [ ] Caching implementation
- [ ] Database read replicas
- [ ] Monitoring and alerting

## 7. Alternative Approaches Considered

### 7.1 Time-Series Database (Rejected)
**Considered**: InfluxDB or TimescaleDB for event storage
**Reason for Rejection**: PostgreSQL with proper indexing sufficient for current scale; adds complexity

### 7.2 Event Streaming (Future Enhancement)
**Considered**: Apache Kafka for real-time event processing
**Status**: Valuable for future real-time analytics but overkill for initial MVP

### 7.3 NoSQL Document Store (Rejected)
**Considered**: MongoDB for flexible event schema
**Reason for Rejection**: Relational data model better suited for educational domain with clear entity relationships

## 8. Assumptions Made

1. **Network Reliability**: Stable internet connectivity in classroom environments
2. **Device Capabilities**: Modern Android devices with adequate storage/processing
3. **User Behavior**: Teachers and students will adapt to digital quiz workflow
4. **Data Volume**: Quiz sessions average 15-20 questions over 30-45 minute periods
5. **Compliance**: Standard educational data privacy requirements (FERPA compliance)

---

**Document Version**: 1.0  
**Last Updated**: January 2024  
**Author**: Educational Analytics Team 