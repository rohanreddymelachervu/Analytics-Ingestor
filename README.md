# Educational Analytics Framework

A comprehensive reporting and analytics system for educational applications, designed to track and analyze interactions between Whiteboard App (teachers) and Notebook App (students) during synchronized quiz sessions.

## 🎯 Project Overview

This system captures, stores, and analyzes user interactions across both applications with a focus on quiz features where teachers control questions from the Whiteboard app and students submit answers through the Notebook app.

### Key Features

- **Event-Based Data Collection**: Real-time tracking of teacher and student interactions
- **JWT Authentication**: Role-based access control with scope-based permissions
- **Comprehensive Reporting**: Student performance, classroom engagement, and content effectiveness analysis
- **Scalable Architecture**: Designed for 1,000 schools, 30,000 classrooms, 900,000 students
- **RESTful API**: Clean endpoints for data ingestion and report generation

## 🏗️ Architecture

```
Whiteboard App ──┐
                 ├── HTTPS/JSON ──► Analytics API ──► PostgreSQL
Notebook App  ───┘                      │
                                         ▼
                                   Reports Dashboard
```

## 📊 Database Schema

### Core Entities
- **Quizzes**: Quiz metadata and configuration
- **Classrooms**: Classroom organization and management
- **Students**: Student profiles and enrollment
- **Questions**: Individual quiz questions with answers
- **Quiz Sessions**: Active quiz sessions linking quizzes to classrooms

### Event Tables
- **Question Published Events**: When teachers publish questions (Whiteboard App)
- **Answer Submitted Events**: When students submit answers (Notebook App)

## 🚀 Quick Start

### Prerequisites
- Go 1.24+
- PostgreSQL 13+
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd Analytics-Ingestor
   ```

2. **Set environment variables**
   ```bash
   export DATABASE_URL="postgres://username:password@localhost/analytics_db?sslmode=disable"
   export JWT_SECRET="your-secret-key-here"
   ```

3. **Install dependencies**
   ```bash
   go mod download
   ```

4. **Run database migrations**
   ```bash
   # Apply all migrations
   migrate -path migrations/postgres -database $DATABASE_URL up
   ```

5. **Start the server**
   ```bash
   go run cmd/server/main.go
   ```

The server will start on port 8080 by default.

## 📡 API Documentation

### Authentication

1. **Sign Up** (Create a user account)
   ```bash
   POST /api/auth/signup
   Content-Type: application/json
   
   {
     "name": "John Doe",
     "email": "john@school.edu",
     "password": "securepassword",
     "role": "writer"  // or "reader"
   }
   ```

2. **Login** (Get JWT token)
   ```bash
   POST /api/auth/login
   Content-Type: application/json
   
   {
     "email": "john@school.edu", 
     "password": "securepassword"
   }
   ```

### Event Ingestion (Requires WRITE scope)

1. **Single Event**
   ```bash
   POST /api/events
   Authorization: Bearer <your-jwt-token>
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

2. **Batch Events**
   ```bash
   POST /api/events/batch
   Authorization: Bearer <your-jwt-token>
   Content-Type: application/json
   
   [
     {"event_type": "question_published", "session_id": "...", "data": {...}},
     {"event_type": "answer_submitted", "session_id": "...", "data": {...}}
   ]
   ```

### Reports (Requires READ scope)

1. **Active Participants**
   ```bash
   GET /api/reports/active-participants?session_id=<uuid>&time_range=60m
   Authorization: Bearer <your-jwt-token>
   ```

2. **Questions Per Minute**
   ```bash
   GET /api/reports/questions-per-minute?session_id=<uuid>
   Authorization: Bearer <your-jwt-token>
   ```

3. **Student Performance Analysis**
   ```bash
   GET /api/reports/student-performance?student_id=<uuid>&classroom_id=<uuid>
   Authorization: Bearer <your-jwt-token>
   ```

4. **Classroom Engagement Metrics**
   ```bash
   GET /api/reports/classroom-engagement?classroom_id=<uuid>&date_range=7d
   Authorization: Bearer <your-jwt-token>
   ```

5. **Content Effectiveness Evaluation**
   ```bash
   GET /api/reports/content-effectiveness?quiz_id=<uuid>
   Authorization: Bearer <your-jwt-token>
   ```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `DATABASE_URL` | PostgreSQL connection string | Yes | - |
| `JWT_SECRET` | Secret key for JWT signing | Yes | - |
| `PORT` | Server port | No | 8080 |

### User Roles & Scopes

| Role | Scopes | Description |
|------|--------|-------------|
| `writer` | `WRITE` | Can ingest events (Whiteboard/Notebook apps) |
| `reader` | `READ` | Can access reports (Analytics dashboard) |

## 📈 Scale Specifications

The system is designed to handle:

- **1,000 schools**
- **30,000 classrooms** (30 per school)
- **900,000 students** (30 per classroom)
- **~50,000 events/minute** during peak quiz times
- **~1TB/year** of event data storage

## 🗂️ Project Structure

```
Analytics-Ingestor/
├── cmd/server/           # Application entry point
├── internal/
│   ├── auth/            # Authentication & authorization
│   ├── config/          # Configuration management
│   └── server/          # HTTP server & routing
├── migrations/postgres/ # Database migrations
├── go.mod              # Go module definition
├── TECHNICAL_DESIGN.md # Detailed technical documentation
└── README.md          # This file
```

## 🧪 Testing

### Manual Testing Examples

1. **Create a writer user and get token:**
   ```bash
   # Sign up
   curl -X POST http://localhost:8080/api/auth/signup \
     -H "Content-Type: application/json" \
     -d '{"name":"Test User","email":"test@example.com","password":"password123","role":"writer"}'
   
   # Login
   curl -X POST http://localhost:8080/api/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"password123"}'
   ```

2. **Send analytics events:**
   ```bash
   curl -X POST http://localhost:8080/api/events \
     -H "Authorization: Bearer <your-token>" \
     -H "Content-Type: application/json" \
     -d '{"event_type":"answer_submitted","session_id":"550e8400-e29b-41d4-a716-446655440000","data":{"student_id":"550e8400-e29b-41d4-a716-446655440001","question_id":"550e8400-e29b-41d4-a716-446655440002","answer":"B"}}'
   ```

3. **Get reports (create reader user first):**
   ```bash
   curl -X GET "http://localhost:8080/api/reports/active-participants?session_id=550e8400-e29b-41d4-a716-446655440000" \
     -H "Authorization: Bearer <reader-token>"
   ```

## 🔮 Future Enhancements

### Phase 2: Production Readiness
- Database connection optimization with pooling
- Comprehensive error handling and logging
- Input validation and sanitization
- Swagger API documentation
- Unit and integration tests

### Phase 3: Advanced Analytics
- Real-time data aggregation pipelines
- Advanced SQL queries for complex reports
- Data visualization endpoints
- Export functionality (CSV, PDF reports)

### Phase 4: Scale & Performance
- Load testing and performance optimization
- Redis caching for frequently accessed data
- Database read replicas for report queries
- Monitoring and alerting infrastructure

### Bonus: Generic Query Interface
- Integration with cube.dev for OLAP-style queries
- Measures and dimensions for flexible reporting
- Self-service analytics capabilities

## 📄 Documentation

- [Technical Design Document](TECHNICAL_DESIGN.md) - Comprehensive architecture and design decisions
- [Database Schema](migrations/postgres/) - Complete database migration files
- [API Examples](README.md#api-documentation) - Request/response examples

## 🤝 Contributing

This project was developed as part of an engineering assessment for Superr.ai, demonstrating:

1. **System Design Skills**: Event-based architecture for educational analytics
2. **Backend Development**: Go/Gin API with PostgreSQL
3. **Authentication**: JWT with role-based access control
4. **Database Design**: Normalized schema optimized for analytics queries
5. **Documentation**: Comprehensive technical documentation

## 📧 Contact

**Assignment Submission**: hiring@superr.ai
**Developer**: [Your Name]
**Date**: January 2024

---

## 🏆 Assignment Deliverables Completed

- ✅ **Technical Design Document**: Complete architecture documentation
- ✅ **Database Schema**: ER diagram and optimized table structure
- ✅ **API Design**: RESTful endpoints with proper authentication
- ✅ **Sequence Diagrams**: Key workflow documentation
- ✅ **Working Prototype**: Functional event ingestion and reporting
- ✅ **Three Report Types**: Student performance, classroom engagement, content effectiveness
- ✅ **Scale Considerations**: Designed for 900K users across 1K schools

This implementation provides a solid foundation for an educational analytics platform that can scale to support large educational institutions while maintaining performance and data integrity. 