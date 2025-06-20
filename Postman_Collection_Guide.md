# Educational Analytics API - Postman Collection Guide

## Overview
This Postman collection provides complete testing coverage for the Educational Analytics Framework API, designed to support Whiteboard and Notebook applications in educational environments.

## 📦 Collection Contents

### Folder Structure
```
🔐 Authentication (4 requests)
├── Sign Up Writer
├── Login Writer  
├── Sign Up Reader
└── Login Reader

📊 Event Ingestion (4 requests)
├── Session Started Event
├── Question Published Event
├── Answer Submitted Event
└── Batch Events

👤 Student Analytics (5 requests)
├── Student Performance
├── Student Performance List
├── Student Activity Summary
├── Session Student Rankings
└── Classroom Student Rankings

🏫 Classroom Analytics (5 requests)
├── Classroom Overview
├── Classroom Engagement
├── Classroom Engagement History
├── Class Performance Summary
└── Classroom Sessions

📝 Quiz & Content Analytics (5 requests)
├── Quiz Summary
├── Content Effectiveness
├── Question Analysis
├── Quiz Questions List
└── Quiz Sessions

🔴 Real-time Analytics (4 requests)
├── Active Participants
├── Questions Per Minute
├── Response Rate
└── Completion Rate

⚡ Performance Analytics (3 requests)
├── Latency Analysis
├── Timeout Analysis
└── Dropoff Analysis
```

## 🚀 Getting Started

### 1. Import Collection
1. Download `Educational_Analytics_API.postman_collection.json`
2. Open Postman
3. Click "Import" → "Upload Files"
4. Select the collection file
5. Click "Import"

### 2. Environment Variables
The collection includes pre-configured variables:

| Variable | Purpose | Example Value |
|----------|---------|---------------|
| `base_url` | API endpoint | `http://localhost:8080/api` |
| `writer_token` | JWT for event ingestion | Auto-populated after login |
| `reader_token` | JWT for analytics access | Auto-populated after login |

### 3. Server Setup
Before testing, ensure your server is running:
```bash
# Start PostgreSQL database
# Start server
export DATABASE_URL='postgres://postgres:root@localhost:5432/ingestor'
export JWT_SECRET='your-jwt-secret'
export KAFKA_ENABLED=false  # or true for Kafka mode
./server
```

## 🔐 Authentication Flow

### User Roles & Scopes
- **Writer Role**: `WRITE` scope - Can ingest events (Whiteboard/Notebook apps)
- **Reader Role**: `READ` scope - Can access analytics reports (Dashboard)

### Authentication Steps
1. **Create Writer User**: Run "Sign Up Writer"
2. **Login Writer**: Run "Login Writer" → Updates `writer_token` variable
3. **Create Reader User**: Run "Sign Up Reader"  
4. **Login Reader**: Run "Login Reader" → Updates `reader_token` variable

### Token Auto-Management
The collection automatically:
- Saves JWT tokens to collection variables
- Uses appropriate tokens for each request type
- Handles token refresh when needed

## 📊 Event Ingestion Testing

### Event Types Supported
- `SESSION_STARTED` - Quiz session initialization
- `QUESTION_PUBLISHED` - Teacher publishes question to students
- `ANSWER_SUBMITTED` - Student submits answer response
- `SESSION_ENDED` - Quiz session completion

### Sample Event Payloads

#### Session Started Event
```json
{
  "event_id": "unique-uuid",
  "event_type": "SESSION_STARTED",
  "timestamp": "2025-06-20T10:01:00Z",
  "session_id": "session-uuid",
  "quiz_id": "quiz-uuid",
  "classroom_id": "classroom-uuid",
  "question_id": "question-uuid"
}
```

#### Answer Submitted Event
```json
{
  "event_id": "unique-uuid",
  "event_type": "ANSWER_SUBMITTED",
  "timestamp": "2025-06-20T10:01:45Z",
  "session_id": "session-uuid",
  "quiz_id": "quiz-uuid",
  "classroom_id": "classroom-uuid",
  "question_id": "question-uuid",
  "student_id": "student-uuid",
  "answer": "A"
}
```

### Batch Processing
Use "Batch Events" for high-throughput scenarios. Send as a **direct array** of events:
```json
[
  {
    "event_id": "unique-uuid-1",
    "event_type": "ANSWER_SUBMITTED",
    "timestamp": "2025-06-20T10:01:50Z",
    "session_id": "session-uuid",
    "quiz_id": "quiz-uuid",
    "classroom_id": "classroom-uuid",
    "question_id": "question-uuid",
    "student_id": "student-uuid-1",
    "answer": "A"
  },
  {
    "event_id": "unique-uuid-2", 
    "event_type": "ANSWER_SUBMITTED",
    "timestamp": "2025-06-20T10:01:55Z",
    "session_id": "session-uuid",
    "quiz_id": "quiz-uuid",
    "classroom_id": "classroom-uuid",
    "question_id": "question-uuid",
    "student_id": "student-uuid-2",
    "answer": "B"
  }
]
```

**Important**: The batch endpoint expects a **direct JSON array**, not an object with an "events" property.

## 📈 Analytics Testing

### Core Metrics Available

#### Student-Level Analytics
- **Individual Performance**: Accuracy, response time, trends
- **Activity Summary**: Cross-session participation data
- **Ranking Systems**: Classroom and session-based rankings

#### Classroom-Level Analytics  
- **Engagement Metrics**: Participation rates, activity levels
- **Performance Summary**: Aggregate classroom statistics
- **Session Management**: Historical session data

#### Quiz & Content Analytics
- **Content Effectiveness**: Quiz performance optimization
- **Question Analysis**: Individual question difficulty metrics
- **Cross-Session Analytics**: Multi-session quiz data

#### Real-time Analytics
- **Active Participants**: Live session monitoring
- **Response Rates**: Real-time participation tracking
- **Throughput Metrics**: Questions per minute analysis

### Query Parameters

#### Common Parameters
| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `classroom_id` | UUID | Classroom identifier | `1a2b3c4d-5e6f-7890-1234-567890abcdef` |
| `session_id` | UUID | Quiz session identifier | `11111111-1111-1111-1111-111111111111` |
| `student_id` | UUID | Student identifier | `9f8e7d6c-5b4a-3928-1746-5a6b7c8d9e0f` |
| `quiz_id` | UUID | Quiz identifier | `3e4d5e6f-7890-1234-5678-90abcdef1234` |

#### Pagination Parameters
| Parameter | Type | Description | Default |
|-----------|------|-------------|---------|
| `page` | Integer | Page number | 1 |
| `limit` | Integer | Items per page | 10 |

#### Time Range Parameters
| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `from_ts` | ISO DateTime | Start time | `2025-01-01T00:00:00Z` |
| `to_ts` | ISO DateTime | End time | `2025-12-31T23:59:59Z` |
| `time_range` | Duration | Relative time range | `1h`, `24h`, `7d` |

## 🧪 Testing Scenarios

### 1. Complete Event Flow Test
1. Run "Session Started Event"
2. Run "Question Published Event"  
3. Run "Answer Submitted Event" (multiple times with different students)
4. Test analytics endpoints to verify data processing

### 2. Real-time Analytics Test
1. Start quiz session with events
2. Call "Active Participants" for live monitoring
3. Check "Response Rate" for participation tracking
4. Monitor "Questions Per Minute" for throughput

### 3. Historical Analytics Test
1. Create multiple quiz sessions over time
2. Test "Student Performance List" with pagination
3. Check "Classroom Engagement History"
4. Verify "Quiz Sessions" data accuracy

### 4. Performance Testing
1. Use "Batch Events" for high-volume ingestion
2. Test "Latency Analysis" for response time metrics
3. Check "Timeout Analysis" for timer behavior
4. Monitor "Dropoff Analysis" for engagement patterns

## ✅ Response Validation

### Success Response Patterns

#### Event Ingestion Success
```json
{
  "event_id": "test-event-123",
  "message": "Event queued successfully",
  "mode": "kafka",  // or "database"
  "timestamp": "2025-06-20T15:00:00Z"
}
```

#### Analytics Response Pattern
```json
{
  "data": { /* analytics results */ },
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total_count": 25,
    "total_pages": 3,
    "has_more": true,
    "has_previous": false
  },
  "metadata": {
    "query_time_ms": 15,
    "cache_hit": false
  }
}
```

### Error Response Patterns

#### Authentication Errors
```json
{
  "error": "invalid or expired token",
  "code": 401
}
```

#### Authorization Errors  
```json
{
  "error": "insufficient permissions for this resource",
  "code": 403
}
```

#### Validation Errors
```json
{
  "error": "student_id and classroom_id are required",
  "code": 400
}
```

## 🔧 Troubleshooting

### Common Issues

1. **Server Not Running**
   - **Error**: `Failed to connect to localhost port 8080`
   - **Solution**: Start the analytics server first

2. **Database Connection Issues**
   - **Error**: `DATABASE_URL is required`
   - **Solution**: Set environment variables properly

3. **Invalid JWT Token**
   - **Error**: `invalid or expired token`
   - **Solution**: Re-run login endpoints to refresh tokens

4. **Insufficient Permissions**
   - **Error**: `insufficient permissions`
   - **Solution**: Use correct token type (writer vs reader)

5. **Missing Required Parameters**
   - **Error**: `student_id and classroom_id are required`
   - **Solution**: Check request documentation for required parameters

### Debug Steps
1. Check server logs for detailed error messages
2. Verify environment variables are set correctly
3. Ensure database is running and accessible
4. Validate request payloads match expected schema
5. Confirm JWT tokens haven't expired

## 📋 Test Data

### Sample UUIDs for Testing
```
Classroom: 1a2b3c4d-5e6f-7890-1234-567890abcdef
Quiz: 3e4d5e6f-7890-1234-5678-90abcdef1234
Question: f1e2d3c4-b5a6-9788-1234-567890abcdef
Session: 11111111-1111-1111-1111-111111111111
Student 1: 9f8e7d6c-5b4a-3928-1746-5a6b7c8d9e0f
Student 2: 8e7d6c5b-4a39-2817-4659-4a5b6c7d8e9f
```

### Creating Test Data
1. Use event ingestion endpoints to create test events
2. Vary timestamps to simulate realistic timing
3. Use multiple students for comprehensive analytics
4. Create multiple sessions for historical data testing

## 🚀 Advanced Usage

### Collection Variables
You can customize the collection by modifying variables:
- Change `base_url` for different environments
- Update timeout values for slow networks
- Modify pagination defaults for testing

### Environment Setup
Create multiple environments for:
- **Development**: `localhost:8080`
- **Staging**: `staging-api.example.com`
- **Production**: `api.example.com`

### Automated Testing
Use Postman's test scripts for:
- Response validation
- Data integrity checks  
- Performance monitoring
- Regression testing

## 📝 Notes

- All timestamps should be in ISO 8601 format
- UUIDs must be valid UUID v4 format
- Event IDs must be unique to prevent duplicates
- Pagination is 1-indexed (page starts at 1)
- Time ranges support various formats (1h, 24h, 7d)
- JWT tokens have expiration times - refresh as needed

## 🔗 Related Documentation

- [API Technical Design Document](EDUCATIONAL_ANALYTICS_TECHNICAL_DESIGN.md)
- [Assignment Deliverables Summary](ASSIGNMENT_DELIVERABLES_SUMMARY.md)
- [Kafka Setup Guide](KAFKA_SETUP.md)
- [Database Schema Documentation](migrations/postgres/) 