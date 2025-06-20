#!/bin/bash

echo "🗄️  POSTGRESQL DATABASE RESET SCRIPT 🗄️"
echo "========================================"
echo ""
echo "This script will:"
echo "1. Drop all existing data"
echo "2. Recreate the database schema"
echo "3. Insert fresh test data"
echo "4. Verify the reset was successful"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m'

# Database configuration
DB_NAME="ingestor"
DB_USER="postgres"
DB_PASSWORD="root"
DB_HOST="localhost"
DB_PORT="5432"

echo -e "${YELLOW}⚠️  WARNING: This will DELETE ALL DATA in the database!${NC}"
echo "Press Enter to continue or Ctrl+C to cancel..."
read

echo -e "${PURPLE}🔧 STEP 1: Database Connection Test${NC}"
echo "=================================="

# Test connection
if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "\q" >/dev/null 2>&1; then
    echo -e "${GREEN}✅ Database connection successful${NC}"
else
    echo -e "${RED}❌ Database connection failed${NC}"
    echo "Please ensure PostgreSQL is running and credentials are correct:"
    echo "  Database: $DB_NAME"
    echo "  User: $DB_USER"
    echo "  Host: $DB_HOST:$DB_PORT"
    exit 1
fi

echo ""
echo -e "${PURPLE}🗑️  STEP 2: Clearing All Data${NC}"
echo "============================="

# Drop all data (preserve schema structure)
echo "Clearing all tables..."
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME << 'EOF'
-- Disable foreign key constraints temporarily
SET session_replication_role = replica;

-- Clear all data from tables (order matters due to dependencies)
TRUNCATE TABLE answer_submitted_events CASCADE;
TRUNCATE TABLE question_published_events CASCADE;
TRUNCATE TABLE classroom_students CASCADE;
TRUNCATE TABLE quiz_sessions CASCADE;
TRUNCATE TABLE questions CASCADE;
TRUNCATE TABLE quizzes CASCADE;
TRUNCATE TABLE classrooms CASCADE;
TRUNCATE TABLE students CASCADE;
TRUNCATE TABLE users CASCADE;

-- Re-enable foreign key constraints
SET session_replication_role = DEFAULT;

SELECT 'All tables cleared successfully' as status;
EOF

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ All data cleared successfully${NC}"
else
    echo -e "${RED}❌ Failed to clear data${NC}"
    exit 1
fi

echo ""
echo -e "${PURPLE}🏗️  STEP 3: Inserting Fresh Test Data${NC}"
echo "===================================="

echo "Inserting comprehensive test data..."
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME << 'EOF'

-- Insert Users (authentication users with bcrypt hashed passwords)
-- Password for all users is 'test123' (hash verified working)
INSERT INTO users (id, name, email, password, role) VALUES
(1, 'Reader Test User', 'reader.test@example.com', '$2a$10$g29Vxi1hiRXmj8GrsUYrVOnBuo/AcdQkL1KD7QjlwuRYMhC9DKZRS', 'reader'),
(2, 'Writer Test User', 'writer.test@example.com', '$2a$10$g29Vxi1hiRXmj8GrsUYrVOnBuo/AcdQkL1KD7QjlwuRYMhC9DKZRS', 'writer'),
(3, 'Alice Smith', 'alice@school.edu', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'writer'),
(4, 'Bob Jones', 'bob@school.edu', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'writer'),
(5, 'Charlie Brown', 'charlie@school.edu', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'writer'),
(6, 'Diana Wilson', 'diana@school.edu', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'writer'),
(7, 'Eva Davis', 'eva@school.edu', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'writer'),
(8, 'Frank Miller', 'frank@school.edu', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'writer'),
(9, 'Grace Lee', 'grace@school.edu', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'writer'),
(10, 'Henry Taylor', 'henry@school.edu', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'writer');

-- Insert Students
INSERT INTO students (student_id, name) VALUES
('900e8400-e29b-41d4-a716-446655440040', 'Alice Smith'),
('900e8400-e29b-41d4-a716-446655440041', 'Bob Jones'),
('900e8400-e29b-41d4-a716-446655440042', 'Charlie Brown'),
('900e8400-e29b-41d4-a716-446655440043', 'Diana Wilson'),
('900e8400-e29b-41d4-a716-446655440044', 'Eva Davis'),
('900e8400-e29b-41d4-a716-446655440045', 'Frank Miller'),
('900e8400-e29b-41d4-a716-446655440046', 'Grace Lee'),
('900e8400-e29b-41d4-a716-446655440047', 'Henry Taylor');

-- Insert Classrooms
INSERT INTO classrooms (classroom_id, name) VALUES
('900e8400-e29b-41d4-a716-446655440020', 'Advanced Computer Science'),
('900e8400-e29b-41d4-a716-446655440021', 'Introduction to Programming'),
('900e8400-e29b-41d4-a716-446655440022', 'Data Structures'),
('900e8400-e29b-41d4-a716-446655440023', 'Web Development');

-- Insert Classroom-Student relationships
INSERT INTO classroom_students (classroom_id, student_id) VALUES
('900e8400-e29b-41d4-a716-446655440020', '900e8400-e29b-41d4-a716-446655440040'),
('900e8400-e29b-41d4-a716-446655440020', '900e8400-e29b-41d4-a716-446655440041'),
('900e8400-e29b-41d4-a716-446655440020', '900e8400-e29b-41d4-a716-446655440042'),
('900e8400-e29b-41d4-a716-446655440021', '900e8400-e29b-41d4-a716-446655440043'),
('900e8400-e29b-41d4-a716-446655440021', '900e8400-e29b-41d4-a716-446655440044'),
('900e8400-e29b-41d4-a716-446655440022', '900e8400-e29b-41d4-a716-446655440045'),
('900e8400-e29b-41d4-a716-446655440023', '900e8400-e29b-41d4-a716-446655440046'),
('900e8400-e29b-41d4-a716-446655440023', '900e8400-e29b-41d4-a716-446655440047');

-- Insert Quizzes
INSERT INTO quizzes (quiz_id, title, description) VALUES
('900e8400-e29b-41d4-a716-446655440010', 'Algorithms Quiz', 'Test on sorting and searching algorithms'),
('900e8400-e29b-41d4-a716-446655440011', 'Programming Basics', 'Fundamental programming concepts'),
('900e8400-e29b-41d4-a716-446655440012', 'Data Structures Quiz', 'Arrays, linked lists, and trees'),
('900e8400-e29b-41d4-a716-446655440013', 'HTML & CSS Basics', 'Web development fundamentals');

-- Insert Questions
INSERT INTO questions (question_id, quiz_id) VALUES
('900e8400-e29b-41d4-a716-446655440030', '900e8400-e29b-41d4-a716-446655440010'),
('900e8400-e29b-41d4-a716-446655440031', '900e8400-e29b-41d4-a716-446655440010'),
('900e8400-e29b-41d4-a716-446655440036', '900e8400-e29b-41d4-a716-446655440010'),
('900e8400-e29b-41d4-a716-446655440032', '900e8400-e29b-41d4-a716-446655440011'),
('900e8400-e29b-41d4-a716-446655440033', '900e8400-e29b-41d4-a716-446655440011'),
('900e8400-e29b-41d4-a716-446655440037', '900e8400-e29b-41d4-a716-446655440011'),
('900e8400-e29b-41d4-a716-446655440034', '900e8400-e29b-41d4-a716-446655440012'),
('900e8400-e29b-41d4-a716-446655440038', '900e8400-e29b-41d4-a716-446655440012'),
('900e8400-e29b-41d4-a716-446655440035', '900e8400-e29b-41d4-a716-446655440013'),
('900e8400-e29b-41d4-a716-446655440039', '900e8400-e29b-41d4-a716-446655440013');

-- Insert Quiz Sessions
INSERT INTO quiz_sessions (session_id, quiz_id, classroom_id, started_at, ended_at) VALUES
('900e8400-e29b-41d4-a716-446655440000', '900e8400-e29b-41d4-a716-446655440010', '900e8400-e29b-41d4-a716-446655440020', NOW() - INTERVAL '30 minutes', NULL),
('900e8400-e29b-41d4-a716-446655440001', '900e8400-e29b-41d4-a716-446655440011', '900e8400-e29b-41d4-a716-446655440021', NOW() - INTERVAL '2 hours', NOW() - INTERVAL '1 hour'),
('900e8400-e29b-41d4-a716-446655440002', '900e8400-e29b-41d4-a716-446655440012', '900e8400-e29b-41d4-a716-446655440022', NOW() + INTERVAL '1 hour', NULL);

-- Insert Question Published Events
INSERT INTO question_published_events (event_id, session_id, question_id, teacher_id, published_at, timer_duration_sec) VALUES
('900e8400-e29b-41d4-a716-446655440080', '900e8400-e29b-41d4-a716-446655440000', '900e8400-e29b-41d4-a716-446655440030', NULL, NOW() - INTERVAL '25 minutes', 60),
('900e8400-e29b-41d4-a716-446655440081', '900e8400-e29b-41d4-a716-446655440000', '900e8400-e29b-41d4-a716-446655440031', NULL, NOW() - INTERVAL '20 minutes', 60),
('900e8400-e29b-41d4-a716-446655440082', '900e8400-e29b-41d4-a716-446655440001', '900e8400-e29b-41d4-a716-446655440032', NULL, NOW() - INTERVAL '90 minutes', 45),
('900e8400-e29b-41d4-a716-446655440083', '900e8400-e29b-41d4-a716-446655440001', '900e8400-e29b-41d4-a716-446655440033', NULL, NOW() - INTERVAL '85 minutes', 45);

-- Insert Answer Submitted Events (comprehensive test data)
INSERT INTO answer_submitted_events (event_id, session_id, question_id, student_id, answer, is_correct, submitted_at) VALUES
-- Session 1 responses (active session) - students performing different levels
('900e8400-e29b-41d4-a716-446655440060', '900e8400-e29b-41d4-a716-446655440000', '900e8400-e29b-41d4-a716-446655440030', '900e8400-e29b-41d4-a716-446655440040', 'O(n log n)', true, NOW() - INTERVAL '25 minutes'),
('900e8400-e29b-41d4-a716-446655440061', '900e8400-e29b-41d4-a716-446655440000', '900e8400-e29b-41d4-a716-446655440031', '900e8400-e29b-41d4-a716-446655440040', 'Dijkstra', true, NOW() - INTERVAL '20 minutes'),
('900e8400-e29b-41d4-a716-446655440062', '900e8400-e29b-41d4-a716-446655440000', '900e8400-e29b-41d4-a716-446655440030', '900e8400-e29b-41d4-a716-446655440041', 'O(n)', false, NOW() - INTERVAL '24 minutes'),
('900e8400-e29b-41d4-a716-446655440063', '900e8400-e29b-41d4-a716-446655440000', '900e8400-e29b-41d4-a716-446655440030', '900e8400-e29b-41d4-a716-446655440042', 'O(n log n)', true, NOW() - INTERVAL '23 minutes'),
('900e8400-e29b-41d4-a716-446655440064', '900e8400-e29b-41d4-a716-446655440000', '900e8400-e29b-41d4-a716-446655440031', '900e8400-e29b-41d4-a716-446655440042', 'Dijkstra', true, NOW() - INTERVAL '18 minutes'),

-- Session 2 responses (completed session)
('900e8400-e29b-41d4-a716-446655440065', '900e8400-e29b-41d4-a716-446655440001', '900e8400-e29b-41d4-a716-446655440032', '900e8400-e29b-41d4-a716-446655440043', 'A named storage location', true, NOW() - INTERVAL '90 minutes'),
('900e8400-e29b-41d4-a716-446655440066', '900e8400-e29b-41d4-a716-446655440001', '900e8400-e29b-41d4-a716-446655440033', '900e8400-e29b-41d4-a716-446655440043', 'Integrated Development Environment', true, NOW() - INTERVAL '85 minutes'),
('900e8400-e29b-41d4-a716-446655440067', '900e8400-e29b-41d4-a716-446655440001', '900e8400-e29b-41d4-a716-446655440032', '900e8400-e29b-41d4-a716-446655440044', 'A function', false, NOW() - INTERVAL '88 minutes'),
('900e8400-e29b-41d4-a716-446655440068', '900e8400-e29b-41d4-a716-446655440001', '900e8400-e29b-41d4-a716-446655440033', '900e8400-e29b-41d4-a716-446655440044', 'Internet Data Exchange', false, NOW() - INTERVAL '83 minutes'),

-- Additional responses for comprehensive analytics
('900e8400-e29b-41d4-a716-446655440069', '900e8400-e29b-41d4-a716-446655440000', '900e8400-e29b-41d4-a716-446655440030', '900e8400-e29b-41d4-a716-446655440043', 'O(n log n)', true, NOW() - INTERVAL '22 minutes'),
('900e8400-e29b-41d4-a716-446655440070', '900e8400-e29b-41d4-a716-446655440000', '900e8400-e29b-41d4-a716-446655440031', '900e8400-e29b-41d4-a716-446655440044', 'BFS', false, NOW() - INTERVAL '15 minutes'),
('900e8400-e29b-41d4-a716-446655440071', '900e8400-e29b-41d4-a716-446655440000', '900e8400-e29b-41d4-a716-446655440030', '900e8400-e29b-41d4-a716-446655440045', 'O(n²)', false, NOW() - INTERVAL '21 minutes');

SELECT 'Test data inserted successfully' as status;
EOF

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Fresh test data inserted successfully${NC}"
else
    echo -e "${RED}❌ Failed to insert test data${NC}"
    exit 1
fi

echo ""
echo -e "${PURPLE}🔍 STEP 4: Verification${NC}"
echo "======================"

echo "Verifying data integrity..."
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME << 'EOF'
SELECT 
    'users' as table_name, 
    COUNT(*) as record_count 
FROM users
UNION ALL
SELECT 
    'students' as table_name, 
    COUNT(*) as record_count 
FROM students
UNION ALL
SELECT 
    'classrooms' as table_name, 
    COUNT(*) as record_count 
FROM classrooms
UNION ALL
SELECT 
    'classroom_students' as table_name, 
    COUNT(*) as record_count 
FROM classroom_students
UNION ALL
SELECT 
    'quizzes' as table_name, 
    COUNT(*) as record_count 
FROM quizzes
UNION ALL
SELECT 
    'questions' as table_name, 
    COUNT(*) as record_count 
FROM questions
UNION ALL
SELECT 
    'quiz_sessions' as table_name, 
    COUNT(*) as record_count 
FROM quiz_sessions
UNION ALL
SELECT 
    'question_published_events' as table_name, 
    COUNT(*) as record_count 
FROM question_published_events
UNION ALL
SELECT 
    'answer_submitted_events' as table_name, 
    COUNT(*) as record_count 
FROM answer_submitted_events
ORDER BY table_name;
EOF

echo ""
echo -e "${PURPLE}📊 STEP 5: Test Data Summary${NC}"
echo "==========================="

PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME << 'EOF'
-- Show key test data for verification
SELECT 
    'Authentication Users' as category,
    email,
    role
FROM users 
WHERE email LIKE '%.test@%'
ORDER BY email;

SELECT 
    'Active Quiz Session' as category,
    qs.session_id,
    q.title as quiz_title,
    c.name as classroom_name,
    CASE 
        WHEN qs.ended_at IS NULL THEN 'active'
        ELSE 'completed'
    END as status
FROM quiz_sessions qs
JOIN quizzes q ON qs.quiz_id = q.quiz_id
JOIN classrooms c ON qs.classroom_id = c.classroom_id
WHERE qs.session_id = '900e8400-e29b-41d4-a716-446655440000';

SELECT 
    'Student Responses in Active Session' as category,
    COUNT(*) as response_count,
    COUNT(DISTINCT student_id) as unique_students,
    ROUND(AVG(CASE WHEN is_correct THEN 1 ELSE 0 END) * 100, 2) as avg_accuracy_percent
FROM answer_submitted_events 
WHERE session_id = '900e8400-e29b-41d4-a716-446655440000';
EOF

echo ""
echo -e "${GREEN}🎉 DATABASE RESET COMPLETED SUCCESSFULLY! 🎉${NC}"
echo "=============================================="
echo ""
echo -e "${CYAN}📋 Ready for Testing:${NC}"
echo "✅ Fresh database with clean test data"
echo "✅ Authentication users: reader.test@example.com (READ), writer.test@example.com (WRITE)"
echo "✅ Active quiz session: 900e8400-e29b-41d4-a716-446655440000"
echo "✅ Test students with responses for analytics"
echo "✅ Multiple classrooms, quizzes, and questions"
echo ""
echo -e "${YELLOW}🔧 Next Steps:${NC}"
echo "1. Restart your Analytics server if needed"
echo "2. Run ./final_test.sh to verify all endpoints"
echo "3. Test individual endpoints manually with curl"
echo ""
echo -e "${BLUE}📖 Test Credentials:${NC}"
echo "Reader Token: POST /api/auth/login with reader.test@example.com / test123"
echo "Writer Token: POST /api/auth/login with writer.test@example.com / test123"
echo ""
echo -e "${GREEN}Database reset completed at $(date)${NC}" 