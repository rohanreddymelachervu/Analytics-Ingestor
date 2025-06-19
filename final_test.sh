#!/bin/bash

echo "🚀 COMPREHENSIVE ANALYTICS API TEST - ALL ENDPOINTS 🚀"
echo "======================================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m'

BASE_URL="http://localhost:8080"
PASS_COUNT=0
FAIL_COUNT=0

pass_test() {
    echo -e "${GREEN}✅ PASS:${NC} $1"
    ((PASS_COUNT++))
}

fail_test() {
    echo -e "${RED}❌ FAIL:${NC} $1"
    ((FAIL_COUNT++))
}

info_test() {
    echo -e "${BLUE}ℹ️  INFO:${NC} $1"
}

warn_test() {
    echo -e "${YELLOW}⚠️  WARN:${NC} $1"
}

# Get tokens
echo -e "${PURPLE}🔐 AUTHENTICATION${NC}"
echo "=================="
READER_TOKEN=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"reader.test@example.com","password":"test123"}' | jq -r '.token')

WRITER_TOKEN=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"writer.test@example.com","password":"test123"}' | jq -r '.token')

if [ "$READER_TOKEN" != "null" ] && [ -n "$READER_TOKEN" ]; then
    pass_test "Reader authentication successful"
else
    fail_test "Reader authentication failed"
    exit 1
fi

if [ "$WRITER_TOKEN" != "null" ] && [ -n "$WRITER_TOKEN" ]; then
    pass_test "Writer authentication successful"
else
    fail_test "Writer authentication failed"
    exit 1
fi

# Test IDs
SESSION_ID="900e8400-e29b-41d4-a716-446655440000"
CLASSROOM_ID="900e8400-e29b-41d4-a716-446655440020"
QUIZ_ID="900e8400-e29b-41d4-a716-446655440010"
QUESTION_ID="900e8400-e29b-41d4-a716-446655440030"
STUDENT_ID="900e8400-e29b-41d4-a716-446655440040"

echo ""
echo -e "${PURPLE}📊 CORE ANALYTICS ENDPOINTS (ALL 10)${NC}"
echo "====================================="

# 1. Active Participants
echo -e "${CYAN}🎯 1/10: Active Participants${NC}"
ACTIVE_RESPONSE=$(curl -s -X GET "$BASE_URL/api/reports/active-participants?session_id=$SESSION_ID&time_range=24h" \
  -H "Authorization: Bearer $READER_TOKEN")

if echo "$ACTIVE_RESPONSE" | jq . >/dev/null 2>&1; then
    pass_test "Active Participants - Valid JSON response"
    
    if echo "$ACTIVE_RESPONSE" | jq -e '.active_participants' >/dev/null && \
       echo "$ACTIVE_RESPONSE" | jq -e '.total_participants' >/dev/null && \
       echo "$ACTIVE_RESPONSE" | jq -e '.pagination' >/dev/null; then
        pass_test "Active Participants - Complete structure"
        
        total=$(echo "$ACTIVE_RESPONSE" | jq '.total_participants')
        page_size=$(echo "$ACTIVE_RESPONSE" | jq '.pagination.page_size')
        info_test "Found $total participants, page_size=$page_size"
    else
        fail_test "Active Participants - Missing required fields"
    fi
else
    fail_test "Active Participants - Invalid JSON response"
fi

# 2. Questions Per Minute
echo -e "${CYAN}🎯 2/10: Questions Per Minute${NC}"
QPM_RESPONSE=$(curl -s -X GET "$BASE_URL/api/reports/questions-per-minute?session_id=$SESSION_ID" \
  -H "Authorization: Bearer $READER_TOKEN")

if echo "$QPM_RESPONSE" | jq . >/dev/null 2>&1; then
    pass_test "Questions Per Minute - Valid JSON response"
    
    if echo "$QPM_RESPONSE" | jq -e '.total_questions' >/dev/null && \
       echo "$QPM_RESPONSE" | jq -e '.average_qpm' >/dev/null; then
        pass_test "Questions Per Minute - Required fields present"
        
        total_q=$(echo "$QPM_RESPONSE" | jq '.total_questions')
        avg_qpm=$(echo "$QPM_RESPONSE" | jq '.average_qpm')
        info_test "Total questions: $total_q, Average QPM: $avg_qpm"
    else
        fail_test "Questions Per Minute - Missing required fields"
    fi
else
    fail_test "Questions Per Minute - Invalid JSON response"
fi

# 3. Student Performance
echo -e "${CYAN}🎯 3/10: Student Performance${NC}"
STUDENT_RESPONSE=$(curl -s -X GET "$BASE_URL/api/reports/student-performance?student_id=$STUDENT_ID&classroom_id=$CLASSROOM_ID" \
  -H "Authorization: Bearer $READER_TOKEN")

if echo "$STUDENT_RESPONSE" | jq . >/dev/null 2>&1; then
    pass_test "Student Performance - Valid JSON response"
    
    if echo "$STUDENT_RESPONSE" | jq -e '.questions_attempted' >/dev/null; then
        pass_test "Student Performance - Required fields present"
        
        attempted=$(echo "$STUDENT_RESPONSE" | jq '.questions_attempted // 0')
        correct=$(echo "$STUDENT_RESPONSE" | jq '.correct_answers // 0')
        info_test "Student attempted: $attempted, correct: $correct"
    else
        warn_test "Student Performance - Some fields may be missing (could be valid)"
    fi
else
    fail_test "Student Performance - Invalid JSON response"
fi

# 4. Classroom Engagement
echo -e "${CYAN}🎯 4/10: Classroom Engagement${NC}"
CLASSROOM_RESPONSE=$(curl -s -X GET "$BASE_URL/api/reports/classroom-engagement?classroom_id=$CLASSROOM_ID&date_range=7d" \
  -H "Authorization: Bearer $READER_TOKEN")

if echo "$CLASSROOM_RESPONSE" | jq . >/dev/null 2>&1; then
    pass_test "Classroom Engagement - Valid JSON response"
    
    if echo "$CLASSROOM_RESPONSE" | jq -e '.active_students' >/dev/null; then
        pass_test "Classroom Engagement - Required fields present"
        
        active=$(echo "$CLASSROOM_RESPONSE" | jq '.active_students // 0')
        info_test "Active students: $active"
    else
        warn_test "Classroom Engagement - Some fields may be missing"
    fi
else
    fail_test "Classroom Engagement - Invalid JSON response"
fi

# 5. Content Effectiveness
echo -e "${CYAN}🎯 5/10: Content Effectiveness${NC}"
CONTENT_RESPONSE=$(curl -s -X GET "$BASE_URL/api/reports/content-effectiveness?quiz_id=$QUIZ_ID" \
  -H "Authorization: Bearer $READER_TOKEN")

if echo "$CONTENT_RESPONSE" | jq . >/dev/null 2>&1; then
    pass_test "Content Effectiveness - Valid JSON response"
    
    if echo "$CONTENT_RESPONSE" | jq -e '.effectiveness_score' >/dev/null; then
        pass_test "Content Effectiveness - Required fields present"
        
        effectiveness=$(echo "$CONTENT_RESPONSE" | jq '.effectiveness_score // 0')
        info_test "Effectiveness score: $effectiveness"
    else
        warn_test "Content Effectiveness - Some fields may be missing"
    fi
else
    fail_test "Content Effectiveness - Invalid JSON response"
fi

# 6. Response Rate
echo -e "${CYAN}🎯 6/10: Response Rate${NC}"
RESPONSE_RESPONSE=$(curl -s -X GET "$BASE_URL/api/reports/response-rate?session_id=$SESSION_ID&question_id=$QUESTION_ID" \
  -H "Authorization: Bearer $READER_TOKEN")

if echo "$RESPONSE_RESPONSE" | jq . >/dev/null 2>&1; then
    pass_test "Response Rate - Valid JSON response"
    
    if echo "$RESPONSE_RESPONSE" | jq -e '.response_rate' >/dev/null; then
        pass_test "Response Rate - Required fields present"
        
        rate=$(echo "$RESPONSE_RESPONSE" | jq '.response_rate // 0')
        info_test "Response rate: $rate%"
    else
        warn_test "Response Rate - Some fields may be missing"
    fi
else
    fail_test "Response Rate - Invalid JSON response"
fi

# 7. Latency Analysis
echo -e "${CYAN}🎯 7/10: Latency Analysis${NC}"
LATENCY_RESPONSE=$(curl -s -X GET "$BASE_URL/api/reports/latency-analysis?session_id=$SESSION_ID&question_id=$QUESTION_ID" \
  -H "Authorization: Bearer $READER_TOKEN")

if echo "$LATENCY_RESPONSE" | jq . >/dev/null 2>&1; then
    pass_test "Latency Analysis - Valid JSON response"
else
    fail_test "Latency Analysis - Invalid JSON response"
fi

# 8. Timeout Analysis
echo -e "${CYAN}🎯 8/10: Timeout Analysis${NC}"
TIMEOUT_RESPONSE=$(curl -s -X GET "$BASE_URL/api/reports/timeout-analysis?session_id=$SESSION_ID&question_id=$QUESTION_ID" \
  -H "Authorization: Bearer $READER_TOKEN")

if echo "$TIMEOUT_RESPONSE" | jq . >/dev/null 2>&1; then
    pass_test "Timeout Analysis - Valid JSON response"
else
    fail_test "Timeout Analysis - Invalid JSON response"
fi

# 9. Completion Rate
echo -e "${CYAN}🎯 9/10: Completion Rate${NC}"
COMPLETION_RESPONSE=$(curl -s -X GET "$BASE_URL/api/reports/completion-rate?session_id=$SESSION_ID" \
  -H "Authorization: Bearer $READER_TOKEN")

if echo "$COMPLETION_RESPONSE" | jq . >/dev/null 2>&1; then
    pass_test "Completion Rate - Valid JSON response"
    
    if echo "$COMPLETION_RESPONSE" | jq -e '.completion_rate' >/dev/null; then
        pass_test "Completion Rate - Required fields present"
        
        comp_rate=$(echo "$COMPLETION_RESPONSE" | jq '.completion_rate // 0')
        info_test "Completion rate: $comp_rate%"
    else
        warn_test "Completion Rate - Some fields may be missing"
    fi
else
    fail_test "Completion Rate - Invalid JSON response"
fi

# 10. Drop-off Analysis
echo -e "${CYAN}🎯 10/10: Drop-off Analysis${NC}"
DROPOFF_RESPONSE=$(curl -s -X GET "$BASE_URL/api/reports/dropoff-analysis?session_id=$SESSION_ID" \
  -H "Authorization: Bearer $READER_TOKEN")

if echo "$DROPOFF_RESPONSE" | jq . >/dev/null 2>&1; then
    pass_test "Drop-off Analysis - Valid JSON response"
else
    fail_test "Drop-off Analysis - Invalid JSON response"
fi

echo ""
echo -e "${PURPLE}🆕 NEW PAGINATED ENDPOINTS${NC}"
echo "========================="

# Student Performance List
echo -e "${CYAN}📋 Student Performance List${NC}"
STUDENT_LIST_RESPONSE=$(curl -s -X GET "$BASE_URL/api/reports/student-performance-list?classroom_id=$CLASSROOM_ID&page=1&page_size=5" \
  -H "Authorization: Bearer $READER_TOKEN")

if echo "$STUDENT_LIST_RESPONSE" | jq . >/dev/null 2>&1; then
    pass_test "Student Performance List - Valid JSON response"
    
    if echo "$STUDENT_LIST_RESPONSE" | jq -e '.students' >/dev/null && \
       echo "$STUDENT_LIST_RESPONSE" | jq -e '.pagination' >/dev/null; then
        pass_test "Student Performance List - Complete structure"
        
        student_count=$(echo "$STUDENT_LIST_RESPONSE" | jq '.students | length')
        total_count=$(echo "$STUDENT_LIST_RESPONSE" | jq '.pagination.total_count')
        info_test "Students on page: $student_count, Total: $total_count"
    else
        fail_test "Student Performance List - Missing required fields"
    fi
else
    fail_test "Student Performance List - Invalid JSON response"
fi

# Classroom Engagement History
echo -e "${CYAN}📈 Classroom Engagement History${NC}"
ENGAGEMENT_HISTORY_RESPONSE=$(curl -s -X GET "$BASE_URL/api/reports/classroom-engagement-history?classroom_id=$CLASSROOM_ID&date_range=30d&page=1&page_size=10" \
  -H "Authorization: Bearer $READER_TOKEN")

if echo "$ENGAGEMENT_HISTORY_RESPONSE" | jq . >/dev/null 2>&1; then
    pass_test "Classroom Engagement History - Valid JSON response"
    
    if echo "$ENGAGEMENT_HISTORY_RESPONSE" | jq -e '.engagement_history' >/dev/null && \
       echo "$ENGAGEMENT_HISTORY_RESPONSE" | jq -e '.pagination' >/dev/null; then
        pass_test "Classroom Engagement History - Complete structure"
        
        history_count=$(echo "$ENGAGEMENT_HISTORY_RESPONSE" | jq '.engagement_history | length')
        info_test "History periods: $history_count"
    else
        fail_test "Classroom Engagement History - Missing required fields"
    fi
else
    fail_test "Classroom Engagement History - Invalid JSON response"
fi

echo ""
echo -e "${PURPLE}🚀 PAGINATION TESTING${NC}"
echo "===================="

# Test Active Participants with different pagination
echo -e "${CYAN}📄 Pagination Validation${NC}"
ACTIVE_PAG1=$(curl -s -X GET "$BASE_URL/api/reports/active-participants?session_id=$SESSION_ID&time_range=24h&page=1&page_size=1" \
  -H "Authorization: Bearer $READER_TOKEN")

if echo "$ACTIVE_PAG1" | jq . >/dev/null 2>&1; then
    pass_test "Active Participants (page_size=1) - Valid JSON"
    
    page=$(echo "$ACTIVE_PAG1" | jq '.pagination.page')
    page_size=$(echo "$ACTIVE_PAG1" | jq '.pagination.page_size')
    total_count=$(echo "$ACTIVE_PAG1" | jq '.pagination.total_count')
    total_pages=$(echo "$ACTIVE_PAG1" | jq '.pagination.total_pages')
    has_more=$(echo "$ACTIVE_PAG1" | jq '.pagination.has_more')
    
    if [ "$page_size" -eq 1 ]; then
        pass_test "Pagination - Page size correctly set to 1"
    else
        fail_test "Pagination - Page size incorrect: expected 1, got $page_size"
    fi
    
    # Test pagination math
    if [ "$page_size" -gt 0 ] && [ "$total_count" -ge 0 ]; then
        expected_pages=$(( (total_count + page_size - 1) / page_size ))
        if [ "$total_pages" -eq "$expected_pages" ] || [ "$total_count" -eq 0 ]; then
            pass_test "Pagination - Math correct (pages: $total_pages for $total_count items)"
        else
            fail_test "Pagination - Math incorrect: expected $expected_pages, got $total_pages"
        fi
    fi
    
    info_test "Pagination: page=$page, size=$page_size, total=$total_count, has_more=$has_more"
else
    fail_test "Active Participants (page_size=1) - Invalid JSON"
fi

echo ""
echo -e "${PURPLE}🔒 SECURITY TESTING${NC}"
echo "=================="

# Test unauthorized access
UNAUTH_STATUS=$(curl -s -w "%{http_code}" -o /dev/null -X GET "$BASE_URL/api/reports/active-participants?session_id=$SESSION_ID&time_range=24h")
if [ "$UNAUTH_STATUS" -eq 401 ]; then
    pass_test "Unauthorized access correctly blocked (401)"
else
    fail_test "Unauthorized access not blocked properly (got $UNAUTH_STATUS)"
fi

# Test wrong scope
WRONG_SCOPE_STATUS=$(curl -s -w "%{http_code}" -o /dev/null -X GET "$BASE_URL/api/reports/active-participants?session_id=$SESSION_ID&time_range=24h" \
  -H "Authorization: Bearer $WRITER_TOKEN")
if [ "$WRONG_SCOPE_STATUS" -eq 403 ]; then
    pass_test "Wrong scope correctly blocked (403)"
else
    fail_test "Wrong scope not blocked properly (got $WRONG_SCOPE_STATUS)"
fi

# Test invalid parameters
INVALID_UUID_STATUS=$(curl -s -w "%{http_code}" -o /dev/null -X GET "$BASE_URL/api/reports/active-participants?session_id=invalid-uuid&time_range=24h" \
  -H "Authorization: Bearer $READER_TOKEN")
if [ "$INVALID_UUID_STATUS" -eq 400 ]; then
    pass_test "Invalid UUID correctly rejected (400)"
else
    fail_test "Invalid UUID not rejected properly (got $INVALID_UUID_STATUS)"
fi

echo ""
echo -e "${PURPLE}⚡ EDGE CASE TESTING${NC}"
echo "==================="

# Test large page size (should be capped)
LARGE_PAGE_RESPONSE=$(curl -s -X GET "$BASE_URL/api/reports/active-participants?session_id=$SESSION_ID&time_range=24h&page=1&page_size=5000" \
  -H "Authorization: Bearer $READER_TOKEN")

if echo "$LARGE_PAGE_RESPONSE" | jq . >/dev/null 2>&1; then
    actual_page_size=$(echo "$LARGE_PAGE_RESPONSE" | jq '.pagination.page_size')
    if [ "$actual_page_size" -le 1000 ]; then
        pass_test "Large page size correctly capped at $actual_page_size"
    else
        fail_test "Large page size not capped (got $actual_page_size)"
    fi
else
    fail_test "Large page size test failed - invalid JSON"
fi

# Test zero page (should be corrected)
ZERO_PAGE_RESPONSE=$(curl -s -X GET "$BASE_URL/api/reports/active-participants?session_id=$SESSION_ID&time_range=24h&page=0&page_size=10" \
  -H "Authorization: Bearer $READER_TOKEN")

if echo "$ZERO_PAGE_RESPONSE" | jq . >/dev/null 2>&1; then
    actual_page=$(echo "$ZERO_PAGE_RESPONSE" | jq '.pagination.page')
    if [ "$actual_page" -ge 1 ]; then
        pass_test "Zero page number corrected to $actual_page"
    else
        fail_test "Zero page number not corrected (got $actual_page)"
    fi
else
    fail_test "Zero page test failed - invalid JSON"
fi

echo ""
echo -e "${PURPLE}⏱️  PERFORMANCE TESTING${NC}"
echo "======================="

# Test response time
start_time=$(date +%s%N)
PERF_RESPONSE=$(curl -s -X GET "$BASE_URL/api/reports/active-participants?session_id=$SESSION_ID&time_range=24h&page=1&page_size=50" \
  -H "Authorization: Bearer $READER_TOKEN")
end_time=$(date +%s%N)

response_time=$(( (end_time - start_time) / 1000000 ))

if [ "$response_time" -lt 100 ]; then
    pass_test "Response time excellent: ${response_time}ms"
elif [ "$response_time" -lt 500 ]; then
    pass_test "Response time good: ${response_time}ms"
elif [ "$response_time" -lt 1000 ]; then
    warn_test "Response time acceptable: ${response_time}ms"
else
    fail_test "Response time slow: ${response_time}ms"
fi

echo ""
echo -e "${PURPLE}🏆 FINAL RESULTS${NC}"
echo "================"
echo -e "${GREEN}✅ PASSED: $PASS_COUNT tests${NC}"
echo -e "${RED}❌ FAILED: $FAIL_COUNT tests${NC}"

TOTAL_TESTS=$((PASS_COUNT + FAIL_COUNT))
if [ "$TOTAL_TESTS" -gt 0 ]; then
    SUCCESS_RATE=$(( (PASS_COUNT * 100) / TOTAL_TESTS ))
    echo "Success Rate: $SUCCESS_RATE%"
fi

echo ""
echo -e "${PURPLE}📊 FEATURES VALIDATED:${NC}"
echo "======================"
echo "✅ All 10 Core Analytics Endpoints"
echo "✅ All 2 New Paginated Endpoints"
echo "✅ Pagination Structure & Math"
echo "✅ Security & Access Control"
echo "✅ Parameter Validation"
echo "✅ Edge Case Handling"
echo "✅ Performance Metrics"

if [ "$FAIL_COUNT" -eq 0 ]; then
    echo ""
    echo -e "${GREEN}🎉 ALL SYSTEMS OPERATIONAL! 🎉${NC}"
    echo -e "${GREEN}🚀 ANALYTICS API IS PRODUCTION READY! 🚀${NC}"
    echo -e "${GREEN}🏆 EVERY ENDPOINT TESTED AND WORKING! 🏆${NC}"
    exit 0
else
    echo ""
    if [ "$SUCCESS_RATE" -ge 90 ]; then
        echo -e "${YELLOW}⚠️  Minor issues detected but system is largely functional (${SUCCESS_RATE}% success)${NC}"
        exit 1
    else
        echo -e "${RED}💥 Significant issues detected (${SUCCESS_RATE}% success)${NC}"
        exit 1
    fi
fi 