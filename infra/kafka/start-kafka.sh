#!/bin/bash

echo "🚀 KAFKA SETUP FOR ANALYTICS INGESTOR 🚀"
echo "========================================"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo ""
echo -e "${BLUE}Step 1: Starting Kafka Infrastructure${NC}"
echo "======================================"

# Start Kafka services
echo "Starting Zookeeper and Kafka..."
docker-compose -f docker-compose.kafka.yml up -d

echo ""
echo -e "${YELLOW}Waiting for Kafka to be ready...${NC}"
sleep 15

echo ""
echo -e "${BLUE}Step 2: Verifying Kafka Setup${NC}"
echo "============================="

# Check if containers are running
echo "Checking container status..."
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" | grep analytics

echo ""
echo "Checking if topic was created..."
docker exec analytics-kafka kafka-topics --list --bootstrap-server localhost:9092

echo ""
echo -e "${GREEN}✅ Kafka setup complete!${NC}"
echo ""
echo -e "${YELLOW}🔧 Next Steps:${NC}"
echo "1. Start the consumer: ./consumer"
echo "2. Start the server with Kafka: KAFKA_ENABLED=true ./server"
echo "3. Test events via API: ./test-kafka.sh"
echo ""
echo -e "${YELLOW}📊 Environment Variables:${NC}"
echo "KAFKA_ENABLED=true     - Enable Kafka mode"
echo "KAFKA_BROKERS=localhost:9092  - Kafka brokers"
echo "KAFKA_TOPIC=quiz-events       - Topic name"
echo ""
echo -e "${YELLOW}🛑 To stop Kafka:${NC}"
echo "docker-compose -f docker-compose.kafka.yml down" 