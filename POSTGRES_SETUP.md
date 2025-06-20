# PostgreSQL Setup for Analytics Ingestor

This document provides complete setup instructions for PostgreSQL database used by the Analytics Ingestor project.

## Quick Start

### Option 1: Automated Setup (Recommended)
```bash
# One-command setup - starts PostgreSQL and runs migrations
./setup-database.sh
```

### Option 2: Manual Setup
```bash
# Start PostgreSQL container
./infra/db/start-postgres.sh

# Run migrations separately
export DATABASE_URL="postgres://postgres:root@localhost:5432/ingestor"
migrate -path migrations/postgres -database "$DATABASE_URL" up
```

### Option 3: Docker Compose Only
```bash
# Start PostgreSQL only
docker-compose -f docker-compose.postgres.yml up -d

# Stop PostgreSQL
docker-compose -f docker-compose.postgres.yml down
```

## Database Configuration

### Connection Details
- **Host**: `localhost:5432`
- **Database**: `ingestor`
- **Username**: `postgres`
- **Password**: `root`
- **Connection String**: `postgres://postgres:root@localhost:5432/ingestor`

### Environment Variables
```bash
export DATABASE_URL="postgres://postgres:root@localhost:5432/ingestor"
export JWT_SECRET="62c23d514144fc4fd1dd75fdfed51791f4b9ee14f153db00411ef0eb0bb62aca"
```

## Docker Compose Configuration

The `docker-compose.postgres.yml` file provides:

### PostgreSQL Service
- **Image**: PostgreSQL 15.4 (stable production version)
- **Container Name**: `analytics-postgres`
- **Port**: 5432 (mapped to host)
- **Network**: `analytics-network` (shared with Kafka)

### Performance Tuning
The PostgreSQL container is configured with optimized settings:
- **Max Connections**: 200
- **Shared Buffers**: 256MB
- **Effective Cache Size**: 1GB
- **Work Memory**: 4MB
- **WAL Configuration**: Optimized for performance
- **Logging**: All statements logged (duration > 1000ms)

### Health Checks
- Automatic health monitoring with `pg_isready`
- 10-second intervals with 5-second timeout
- 5 retry attempts before marking as unhealthy

## Database Schema

The project uses 9 tables created through migrations:

### Core Tables
1. **users** - Authentication and user management
2. **quizzes** - Quiz definitions
3. **classrooms** - Classroom information
4. **students** - Student records
5. **questions** - Question definitions
6. **quiz_sessions** - Active quiz sessions
7. **classroom_students** - Many-to-many relationship

### Event Tables
8. **question_published_events** - Published question events
9. **answer_submitted_events** - Student answer submissions

## Migrations

### Running Migrations
```bash
# Install golang-migrate (if not installed)
brew install golang-migrate  # macOS
# or
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path migrations/postgres -database "postgres://postgres:root@localhost:5432/ingestor" up
```

### Migration Files Location
- **Path**: `migrations/postgres/`
- **Format**: `000001_init_schema.up.sql` / `000001_init_schema.down.sql`
- **Count**: 9 migration pairs (up/down)

## Usage with Applications

### Server Application
```bash
export DATABASE_URL="postgres://postgres:root@localhost:5432/ingestor"
export JWT_SECRET="62c23d514144fc4fd1dd75fdfed51791f4b9ee14f153db00411ef0eb0bb62aca"
./server
```

### Consumer Application
```bash
export DATABASE_URL="postgres://postgres:root@localhost:5432/ingestor"
./consumer
```

### With Kafka (Full Setup)
```bash
# Start PostgreSQL
./setup-database.sh

# Start Kafka (in another terminal)
./infra/kafka/start-kafka.sh

# Start applications with Kafka enabled
export DATABASE_URL="postgres://postgres:root@localhost:5432/ingestor"
export JWT_SECRET="62c23d514144fc4fd1dd75fdfed51791f4b9ee14f153db00411ef0eb0bb62aca"
export KAFKA_ENABLED=true
./server &
./consumer &
```

## Troubleshooting

### Common Issues

#### 1. Port 5432 Already in Use
```bash
# Check what's using port 5432
lsof -i :5432

# Stop local PostgreSQL if running
brew services stop postgresql  # macOS
sudo systemctl stop postgresql  # Linux

# Or use different port in docker-compose.postgres.yml
ports:
  - "5433:5432"  # Change to 5433
```

#### 2. Database Connection Errors
```bash
# Check container status
docker ps --filter "name=analytics-postgres"

# Check container logs
docker logs analytics-postgres

# Test connection manually
docker exec -it analytics-postgres psql -U postgres -d ingestor
```

#### 3. Migration Failures
```bash
# Check current migration version
migrate -path migrations/postgres -database "postgres://postgres:root@localhost:5432/ingestor" version

# Force migration version (if needed)
migrate -path migrations/postgres -database "postgres://postgres:root@localhost:5432/ingestor" force 9

# Drop database and recreate (CAUTION: loses all data)
docker exec analytics-postgres dropdb -U postgres ingestor
docker exec analytics-postgres createdb -U postgres ingestor
migrate -path migrations/postgres -database "postgres://postgres:root@localhost:5432/ingestor" up
```

#### 4. Performance Issues
```bash
# Monitor PostgreSQL performance
docker exec analytics-postgres psql -U postgres -d ingestor -c "
SELECT 
  query,
  calls,
  total_time,
  mean_time,
  rows
FROM pg_stat_statements 
ORDER BY total_time DESC 
LIMIT 10;"

# Check active connections
docker exec analytics-postgres psql -U postgres -d ingestor -c "
SELECT count(*) as active_connections 
FROM pg_stat_activity 
WHERE state = 'active';"
```

## Database Administration

### Backup Database
```bash
# Create backup
docker exec analytics-postgres pg_dump -U postgres ingestor > backup.sql

# Restore backup
docker exec -i analytics-postgres psql -U postgres ingestor < backup.sql
```

### Access Database Shell
```bash
# Connect to database
docker exec -it analytics-postgres psql -U postgres -d ingestor

# Useful commands in psql:
# \dt          - List tables
# \d tablename - Describe table
# \q           - Quit
```

### Monitor Database
```bash
# Real-time monitoring
docker exec analytics-postgres psql -U postgres -d ingestor -c "
SELECT 
  schemaname,
  tablename,
  n_tup_ins as inserts,
  n_tup_upd as updates,
  n_tup_del as deletes
FROM pg_stat_user_tables;"
```

## Network Configuration

The PostgreSQL container uses the `analytics-network` Docker network, which is shared with:
- Kafka services (`analytics-kafka`, `analytics-zookeeper`)
- Any future services in the analytics ecosystem

This allows containers to communicate using service names as hostnames.

## Security Considerations

### Production Deployment
For production, update these security settings:

1. **Change default password**:
   ```yaml
   environment:
     POSTGRES_PASSWORD: "your-secure-password"
   ```

2. **Enable SSL**:
   ```yaml
   command: >
     postgres
     -c ssl=on
     -c ssl_cert_file=/path/to/cert.pem
   ```

3. **Restrict connections**:
   ```yaml
   command: >
     postgres
     -c listen_addresses=localhost
   ```

4. **Use secrets management**:
   ```yaml
   environment:
     POSTGRES_PASSWORD_FILE: /run/secrets/postgres_password
   secrets:
     - postgres_password
   ```

## Integration with Analytics Ingestor

This PostgreSQL setup is specifically designed for the Analytics Ingestor project:

- **Schema**: Optimized for educational analytics workloads
- **Indexing**: Strategic indexes for analytical queries
- **Performance**: Tuned for high-volume event ingestion
- **Scalability**: Ready for 900,000+ students across 1,000 schools
- **Analytics**: Support for real-time and historical reporting

The database serves as the primary data store for:
- User authentication and authorization
- Quiz and classroom management
- Real-time event processing
- Analytics and reporting queries
- Performance metrics and insights 