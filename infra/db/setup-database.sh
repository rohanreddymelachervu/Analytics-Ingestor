#!/bin/bash

set -e  # Exit on any error

echo "🗄️  ANALYTICS INGESTOR DATABASE SETUP 🗄️"
echo "=========================================="

# Check if docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    echo "❌ docker-compose is not installed. Please install docker-compose and try again."
    exit 1
fi

# Check if migrate tool is available
if ! command -v migrate &> /dev/null; then
    echo "⚠️  migrate tool not found. Installing..."
    echo "Please install golang-migrate:"
    echo "  macOS: brew install golang-migrate"
    echo "  Linux: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate"
    echo ""
    echo "Or run: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
    exit 1
fi

echo "Step 1: Starting PostgreSQL Container"
echo "====================================="

# Stop any existing PostgreSQL container
echo "Stopping existing PostgreSQL containers..."
docker-compose -f docker-compose.postgres.yml down 2>/dev/null || true

# Start PostgreSQL
echo "Starting PostgreSQL..."
docker-compose -f docker-compose.postgres.yml up -d postgres

echo "Waiting for PostgreSQL to be ready..."
for i in {1..30}; do
    if docker exec analytics-postgres pg_isready -U postgres -d ingestor > /dev/null 2>&1; then
        echo "✅ PostgreSQL is ready!"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "❌ PostgreSQL failed to start within 30 seconds"
        docker logs analytics-postgres
        exit 1
    fi
    echo "Waiting... ($i/30)"
    sleep 1
done

echo ""
echo "Step 2: Running Database Migrations"
echo "=================================="

export DATABASE_URL="postgres://postgres:root@localhost:5432/ingestor"

echo "Running migrations..."
migrate -path migrations/postgres -database "$DATABASE_URL" up

if [ $? -eq 0 ]; then
    echo "✅ Migrations completed successfully!"
else
    echo "❌ Migration failed!"
    exit 1
fi

echo ""
echo "Step 3: Verifying Database Setup"
echo "==============================="

# Test database connection
echo "Testing database connection..."
docker exec analytics-postgres psql -U postgres -d ingestor -c "SELECT 'Database connection successful!' as status;" > /dev/null 2>&1

if [ $? -eq 0 ]; then
    echo "✅ Database connection verified!"
else
    echo "❌ Database connection failed!"
    exit 1
fi

# Check tables
echo "Checking database tables..."
TABLES=$(docker exec analytics-postgres psql -U postgres -d ingestor -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';")
TABLES=$(echo $TABLES | xargs)  # Trim whitespace

echo "Found $TABLES tables in the database"

echo ""
echo "🎉 DATABASE SETUP COMPLETE! 🎉"
echo "============================="
echo ""
echo "📊 Database Details:"
echo "  Host: localhost:5432"
echo "  Database: ingestor"  
echo "  Username: postgres"
echo "  Password: root"
echo ""
echo "🔗 Connection String:"
echo "  DATABASE_URL='postgres://postgres:root@localhost:5432/ingestor'"
echo ""
echo "🚀 Next Steps:"
echo "1. Start the server:"
echo "   export DATABASE_URL='postgres://postgres:root@localhost:5432/ingestor'"
echo "   export JWT_SECRET='62c23d514144fc4fd1dd75fdfed51791f4b9ee14f153db00411ef0eb0bb62aca'"
echo "   ./server"
echo ""
echo "2. Start the consumer (in another terminal):"
echo "   export DATABASE_URL='postgres://postgres:root@localhost:5432/ingestor'"
echo "   ./consumer"
echo ""
echo "3. Optional - Start with Kafka:"
echo "   export KAFKA_ENABLED=true"
echo "   ./infra/kafka/start-kafka.sh"
echo ""
echo "🛑 To stop PostgreSQL:"
echo "   docker-compose -f docker-compose.postgres.yml down"
echo ""
echo "📝 To view logs:"
echo "   docker logs analytics-postgres" 