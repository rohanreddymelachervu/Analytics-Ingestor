#!/bin/bash

echo "🧪 KAFKA INTEGRATION TEST 🧪"
echo "============================="

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

# Get authentication token
echo -e "${BLUE}Step 1: Getting authentication token${NC}"
echo "===================================="

WRITER_TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "writer.test@example.com", "password": "test123"}' | \
  grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$WRITER_TOKEN" ]; then
    echo -e "${RED}❌ Failed to get authentication token${NC}"
    echo "Make sure the server is running and database is set up"
    exit 1
fi

echo -e "${GREEN}✅ Got authentication token${NC}"

echo ""
echo -e "${BLUE}Step 2: Testing Kafka event publishing${NC}"
echo "======================================"

# Test event data - using proper UUIDs
EVENT_DATA='{
  "event_id": "900e8400-e29b-41d4-a716-446655440099",
  "event_type": "QUESTION_PUBLISHED",
  "timestamp": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'",
  "session_id": "900e8400-e29b-41d4-a716-446655440000",
  "quiz_id": "900e8400-e29b-41d4-a716-446655440010",
  "classroom_id": "900e8400-e29b-41d4-a716-446655440020",
  "question_id": "900e8400-e29b-41d4-a716-446655440030",
  "timer_sec": 60
}'

echo "Sending test event to Kafka..."
RESPONSE=$(curl -s -w "HTTPSTATUS:%{http_code}" \
  -X POST http://localhost:8080/api/events \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $WRITER_TOKEN" \
  -d "$EVENT_DATA")

HTTP_STATUS=$(echo $RESPONSE | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
RESPONSE_BODY=$(echo $RESPONSE | sed -e 's/HTTPSTATUS\:.*//g')

echo "HTTP Status: $HTTP_STATUS"
echo "Response: $RESPONSE_BODY"

if [ "$HTTP_STATUS" = "201" ]; then
    if echo "$RESPONSE_BODY" | grep -q "kafka"; then
        echo -e "${GREEN}✅ Event successfully published to Kafka${NC}"
    else
        echo -e "${YELLOW}⚠️  Event processed but not via Kafka (direct mode)${NC}"
    fi
else
    echo -e "${RED}❌ Failed to publish event (HTTP $HTTP_STATUS)${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}🎉 Kafka integration test completed!${NC}" 