#!/bin/bash

echo "🚀 POSTGRESQL SETUP FOR ANALYTICS INGESTOR 🚀"
echo "=============================================="

echo "Step 1: Starting PostgreSQL Database"
echo "====================================="
echo "Starting PostgreSQL container..."

# Start PostgreSQL with docker-compose
docker-compose -f docker-compose.postgres.yml up -d

echo "Waiting for PostgreSQL to be ready..."
sleep 5

echo "Step 2: Verifying PostgreSQL Setup"
echo "=================================="
echo "Checking container status..."
docker ps --filter "name=analytics-postgres" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

echo "Testing database connection..."
docker exec analytics-postgres pg_isready -U postgres -d ingestor

if [ $? -eq 0 ]; then
    echo "✅ PostgreSQL setup complete!"
    echo ""
    echo "🔧 Next Steps:"
    echo "1. Run migrations: migrate -path migrations/postgres -database \"postgres://postgres:root@localhost:5432/ingestor\" up"
    echo "2. Start the server: export DATABASE_URL='postgres://postgres:root@localhost:5432/ingestor' && ./server"
    echo "3. Start the consumer: export DATABASE_URL='postgres://postgres:root@localhost:5432/ingestor' && ./consumer"
    echo ""
    echo "📊 Database Details:"
    echo "Host: localhost:5432"
    echo "Database: ingestor"
    echo "Username: postgres"
    echo "Password: root"
    echo ""
    echo "🔗 Connection String:"
    echo "DATABASE_URL='postgres://postgres:root@localhost:5432/ingestor'"
    echo ""
    echo "🛑 To stop PostgreSQL:"
    echo "docker-compose -f docker-compose.postgres.yml down"
else
    echo "❌ PostgreSQL setup failed!"
    echo "Check docker logs: docker logs analytics-postgres"
    exit 1
fi 