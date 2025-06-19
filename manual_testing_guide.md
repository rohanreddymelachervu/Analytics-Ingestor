# 🧪 Manual Testing Guide for Analytics Ingestor API

This guide provides step-by-step instructions for manually testing every endpoint after resetting the database.

## 📋 Prerequisites

1. **Reset Database**: Run `./reset_database.sh` first
2. **Start Server**: Ensure the Analytics server is running (`./server`)
3. **Install jq**: For JSON formatting (`brew install jq` on macOS)

## 🔐 Step 1: Authentication

First, get authentication tokens:

```bash
# Get Reader Token (READ scope)
READER_TOKEN=$(curl -s -X POST "http://localhost:8080/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"reader.test@example.com","password":"test123"}' | jq -r '.token')

echo "Reader Token: $READER_TOKEN"

# Get Writer Token (WRITE scope)
WRITER_TOKEN=$(curl -s -X POST "http://localhost:8080/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"writer.test@example.com","password":"test123"}' | jq -r '.token')

echo "Writer Token: $WRITER_TOKEN"
```

## 📊 Step 2: Test Core Analytics Endpoints (All 10)

### 1. Active Participants
```bash
curl -s -X GET "http://localhost:8080/api/reports/active-participants?session_id=900e8400-e29b-41d4-a716-446655440000&time_range=24h" \
  -H "Authorization: Bearer $READER_TOKEN" | jq .
```

**Expected**: 3 active participants with pagination

### 2. Questions Per Minute
```bash
curl -s -X GET "http://localhost:8080/api/reports/questions-per-minute?session_id=900e8400-e29b-41d4-a716-446655440000" \
  -H "Authorization: Bearer $READER_TOKEN" | jq .
```

**Expected**: Questions count and average QPM

### 3. Student Performance
```bash
curl -s -X GET "http://localhost:8080/api/reports/student-performance?student_id=900e8400-e29b-41d4-a716-446655440040&classroom_id=900e8400-e29b-41d4-a716-446655440020" \
  -H "Authorization: Bearer $READER_TOKEN" | jq .
```

**Expected**: Student's performance metrics

### 4. Classroom Engagement
```bash
curl -s -X GET "http://localhost:8080/api/reports/classroom-engagement?classroom_id=900e8400-e29b-41d4-a716-446655440020&date_range=7d" \
  -H "Authorization: Bearer $READER_TOKEN" | jq .
```

**Expected**: Classroom engagement data

### 5. Content Effectiveness
```bash
curl -s -X GET "http://localhost:8080/api/reports/content-effectiveness?quiz_id=900e8400-e29b-41d4-a716-446655440010" \
  -H "Authorization: Bearer $READER_TOKEN" | jq .
```

**Expected**: Content effectiveness score

### 6. Response Rate
```bash
curl -s -X GET "http://localhost:8080/api/reports/response-rate?session_id=900e8400-e29b-41d4-a716-446655440000&question_id=900e8400-e29b-41d4-a716-446655440030" \
  -H "Authorization: Bearer $READER_TOKEN" | jq .
```

**Expected**: Response rate metrics

### 7. Latency Analysis
```bash
curl -s -X GET "http://localhost:8080/api/reports/latency-analysis?session_id=900e8400-e29b-41d4-a716-446655440000&question_id=900e8400-e29b-41d4-a716-446655440030" \
  -H "Authorization: Bearer $READER_TOKEN" | jq .
```

**Expected**: Response time analytics

### 8. Timeout Analysis
```bash
curl -s -X GET "http://localhost:8080/api/reports/timeout-analysis?session_id=900e8400-e29b-41d4-a716-446655440000&question_id=900e8400-e29b-41d4-a716-446655440030" \
  -H "Authorization: Bearer $READER_TOKEN" | jq .
```

**Expected**: Timeout metrics

### 9. Completion Rate
```bash
curl -s -X GET "http://localhost:8080/api/reports/completion-rate?session_id=900e8400-e29b-41d4-a716-446655440000" \
  -H "Authorization: Bearer $READER_TOKEN" | jq .
```

**Expected**: Session completion rate

### 10. Drop-off Analysis
```bash
curl -s -X GET "http://localhost:8080/api/reports/dropoff-analysis?session_id=900e8400-e29b-41d4-a716-446655440000" \
  -H "Authorization: Bearer $READER_TOKEN" | jq .
```

**Expected**: Student drop-off patterns

## 🆕 Step 3: Test New Paginated Endpoints

### 11. Student Performance List
```bash
curl -s -X GET "http://localhost:8080/api/reports/student-performance-list?classroom_id=900e8400-e29b-41d4-a716-446655440020&page=1&page_size=5" \
  -H "Authorization: Bearer $READER_TOKEN" | jq .
```

**Expected**: Paginated student performance data

### 12. Classroom Engagement History
```bash
curl -s -X GET "http://localhost:8080/api/reports/classroom-engagement-history?classroom_id=900e8400-e29b-41d4-a716-446655440020&date_range=30d&page=1&page_size=10" \
  -H "Authorization: Bearer $READER_TOKEN" | jq .
```

**Expected**: Paginated engagement history

## 🚀 Step 4: Test Pagination Features

### Small Page Size
```bash
curl -s -X GET "http://localhost:8080/api/reports/active-participants?session_id=900e8400-e29b-41d4-a716-446655440000&time_range=24h&page=1&page_size=1" \
  -H "Authorization: Bearer $READER_TOKEN" | jq .
```

**Expected**: Only 1 participant, `has_more: true`

### Page 2
```bash
curl -s -X GET "http://localhost:8080/api/reports/active-participants?session_id=900e8400-e29b-41d4-a716-446655440000&time_range=24h&page=2&page_size=1" \
  -H "Authorization: Bearer $READER_TOKEN" | jq .
```

**Expected**: Second participant, `has_previous: true`

### Large Page Size (should be capped)
```bash
curl -s -X GET "http://localhost:8080/api/reports/active-participants?session_id=900e8400-e29b-41d4-a716-446655440000&time_range=24h&page=1&page_size=5000" \
  -H "Authorization: Bearer $READER_TOKEN" | jq .
```

**Expected**: page_size capped at 1000

## 🔒 Step 5: Test Security Features

### Unauthorized Access (should fail with 401)
```bash
curl -s -w "%{http_code}" -X GET "http://localhost:8080/api/reports/active-participants?session_id=900e8400-e29b-41d4-a716-446655440000&time_range=24h"
```

**Expected**: HTTP 401

### Wrong Scope (should fail with 403)
```bash
curl -s -w "%{http_code}" -X GET "http://localhost:8080/api/reports/active-participants?session_id=900e8400-e29b-41d4-a716-446655440000&time_range=24h" \
  -H "Authorization: Bearer $WRITER_TOKEN"
```

**Expected**: HTTP 403

### Invalid UUID (should fail with 400)
```bash
curl -s -w "%{http_code}" -X GET "http://localhost:8080/api/reports/active-participants?session_id=invalid-uuid&time_range=24h" \
  -H "Authorization: Bearer $READER_TOKEN"
```

**Expected**: HTTP 400

## ⚡ Step 6: Test Edge Cases

### Zero Page (should be corrected to 1)
```bash
curl -s -X GET "http://localhost:8080/api/reports/active-participants?session_id=900e8400-e29b-41d4-a716-446655440000&time_range=24h&page=0&page_size=10" \
  -H "Authorization: Bearer $READER_TOKEN" | jq '.pagination.page'
```

**Expected**: page = 1

### Negative Page Size (should be corrected)
```bash
curl -s -X GET "http://localhost:8080/api/reports/active-participants?session_id=900e8400-e29b-41d4-a716-446655440000&time_range=24h&page=1&page_size=-5" \
  -H "Authorization: Bearer $READER_TOKEN" | jq '.pagination.page_size'
```

**Expected**: positive page_size

## ✍️ Step 7: Test Event Creation (WRITE scope)

### Create Event (should work with WRITER_TOKEN)
```bash
curl -s -X POST "http://localhost:8080/api/events" \
  -H "Authorization: Bearer $WRITER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "event_id": "test-'$(date +%s)'-001",
    "event_type": "QUESTION_PUBLISHED",
    "timestamp": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'",
    "session_id": "900e8400-e29b-41d4-a716-446655440000",
    "quiz_id": "900e8400-e29b-41d4-a716-446655440010",
    "classroom_id": "900e8400-e29b-41d4-a716-446655440020",
    "question_id": "900e8400-e29b-41d4-a716-446655440030",
    "teacher_id": "900e8400-e29b-41d4-a716-446655440050",
    "timer_sec": 30
  }' | jq .
```

**Expected**: HTTP 201 or 500 (foreign key constraint)

## 📊 Step 8: Verify Data Consistency

### Check Pagination Math
```bash
# Get total count
TOTAL=$(curl -s -X GET "http://localhost:8080/api/reports/active-participants?session_id=900e8400-e29b-41d4-a716-446655440000&time_range=24h" \
  -H "Authorization: Bearer $READER_TOKEN" | jq '.pagination.total_count')

echo "Total participants: $TOTAL"

# Test page math with page_size=2
curl -s -X GET "http://localhost:8080/api/reports/active-participants?session_id=900e8400-e29b-41d4-a716-446655440000&time_range=24h&page=1&page_size=2" \
  -H "Authorization: Bearer $READER_TOKEN" | jq '.pagination'
```

**Expected**: Correct total_pages calculation

## 🏆 Validation Checklist

- [ ] All 10 core endpoints return valid JSON
- [ ] All 2 new paginated endpoints work correctly
- [ ] Pagination structure includes all required fields
- [ ] Pagination math is correct (total_pages calculation)
- [ ] Page size limits are enforced (max 1000)
- [ ] Page numbers are corrected (min 1)
- [ ] Authentication works (401 for no token)
- [ ] Authorization works (403 for wrong scope)
- [ ] Parameter validation works (400 for invalid UUIDs)
- [ ] Edge cases are handled gracefully
- [ ] Response times are reasonable (< 100ms typical)

## 🎯 Test Data Reference

- **Active Session ID**: `900e8400-e29b-41d4-a716-446655440000`
- **Classroom ID**: `900e8400-e29b-41d4-a716-446655440020`
- **Quiz ID**: `900e8400-e29b-41d4-a716-446655440010`
- **Question ID**: `900e8400-e29b-41d4-a716-446655440030`
- **Student ID**: `900e8400-e29b-41d4-a716-446655440040`
- **Teacher ID**: `900e8400-e29b-41d4-a716-446655440050`

## 🔧 Troubleshooting

- **Database connection issues**: Check PostgreSQL is running and credentials in `reset_database.sh`
- **Server not responding**: Ensure `./server` is running on port 8080
- **Authentication failures**: Verify test users exist after database reset
- **Empty responses**: Check if test data was inserted correctly

---

**🎉 Success Criteria**: All endpoints should return valid JSON with proper pagination structure and expected data counts matching the fresh test data! 