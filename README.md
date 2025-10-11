# 📄 dafon-cv-api

[![Go Version](https://img.shields.io/badge/Go-1.24.1-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](docker-compose.yml)
[![OpenAI](https://img.shields.io/badge/OpenAI-Integrated-orange.svg)](https://openai.com/)
[![MySQL](https://img.shields.io/badge/MySQL-8.0-blue.svg)](https://mysql.com/)

A comprehensive REST API for professional CV management with advanced AI-powered content generation, built with Go and Clean Architecture principles.

## 🚀 Quick Start

```bash
# Clone and start the application
git clone https://github.com/Daniel-Fonseca-da-Silva/dafon-cv-api.git
cd dafon-cv-api

# Create .env file with your API keys
cp .env.example .env  # Edit with your keys

# Start with Docker
docker compose up -d --build

# Test the API
curl http://localhost:8080/health
```

## 🧭 Table of Contents

1. Purpose and Value Proposition
2. Architecture Overview
3. Detailed Flows and Diagrams
4. Data Model (ERD)
5. API Documentation
6. Configuration and Environment
7. Security and Hardening
8. Observability (Logging, Health, Metrics)
9. DevOps and Deployment (Docker, Compose, CI/CD)
10. Performance and Scalability
11. Troubleshooting and Runbooks
12. Contributing, License, Author

---

## 🎯 Purpose and Value Proposition

**dafon-cv-api** is a complete solution for creating and managing professional CVs, offering:

- ✅ **Static Token Authentication** - Secure API access with configurable tokens
- ✅ **Complete CV Management** - Full CRUD operations for curriculums, works, and education
- ✅ **AI-Powered Content Generation** - 7 specialized AI endpoints for professional content
- ✅ **User & Configuration Management** - Comprehensive user profiles and settings
- ✅ **Email Integration** - Resend-powered email functionality
- ✅ **Robust Data Validation** - Advanced validation for emails, phones, and data integrity
- ✅ **Rate Limiting System** - Redis-based rate limiting with configurable limits
- ✅ **Clean Architecture** - SOLID principles with clear separation of concerns
- ✅ **Production-Ready** - Docker containerization with health checks and monitoring

**Target Users:** Professionals, recruiters, freelancers, and students who need high-quality resumes with AI assistance and structured data management.

---

## 🏗️ Architecture Overview

The project follows **Clean Architecture** principles with clear separation of concerns:

```
├── cmd/api/                    # Application entry point
├── internal/
│   ├── config/                 # Application configuration & environment
│   ├── database/               # Database configuration (GORM, migrations)
│   ├── dto/                    # Data Transfer Objects (request/response)
│   ├── errors/                 # Custom error types and wrappers
│   ├── handlers/               # Presentation layer (HTTP handlers)
│   │   ├── configuration_handler.go
│   │   ├── curriculum_handler.go
│   │   ├── email_handler.go
│   │   ├── generate_*_ai_handler.go  # 7 AI generation handlers
│   │   └── user_handler.go
│   ├── middleware/             # Middlewares (Static Token Auth, CORS)
│   ├── models/                 # Domain entities (GORM models)
│   ├── ratelimit/             # Rate limiting system (Redis-based)
│   ├── redis/                 # Redis connection and configuration
│   ├── repositories/           # Data access layer (interfaces + implementations)
│   ├── routes/                 # Route definitions and setup
│   │   ├── configuration_routes.go
│   │   ├── curriculum_routes.go
│   │   ├── email_routes.go
│   │   ├── generate_*_ai_routes.go  # 7 AI generation routes
│   │   └── user_routes.go
│   ├── usecases/               # Business logic layer
│   ├── validation/              # Custom validation rules
│   └── validators/             # Validation utilities
```

### High-level Architecture Diagram

```mermaid
flowchart LR
  subgraph Client
    A[Web/Mobile Client]
  end

  subgraph API
    B[GIN Router]
    C[Middleware\nStatic Token Auth]
    D[Handlers Layer]
    E[Use Cases Layer]
    F[Repositories Layer]
  end

  subgraph External Services
    G[(MySQL 8.0)]
    H[OpenAI API]
    I[Resend Email]
    J[(Redis Cache)]
  end

  A -->|HTTP/JSON| B --> C --> D --> E --> F --> G
  E --> H
  E --> I
  C --> J
```

---

## 🔁 Detailed Flows and Diagrams

### Static Token Authentication - Flow

```mermaid
sequenceDiagram
  autonumber
  participant C as Client
  participant API as dafon-cv-api
  participant DB as MySQL

  C->>API: Request with Static Token
  API->>API: Validate Static Token
  alt Valid Token
    API->>DB: Process Request
    API-->>C: 200 OK + Response
  else Invalid Token
    API-->>C: 401 Unauthorized
  end
```

### AI Content Generation - Flow

```mermaid
flowchart TD
  A[Request: POST /ai/generate-intro] --> B[Rate Limiting Check]
  B -->|Allowed| C[Validate payload]
  B -->|Rate Limited| D[429 Too Many Requests]
  C --> E[Build prompt]
  E --> F[OpenAI Completion]
  F -->|Success| G[Sanitize and shape response]
  F -->|Error| H[Return error with context]
  G --> I[200 OK: filtered content]
```

### Rate Limiting - Flow

```mermaid
sequenceDiagram
  participant C as Client
  participant API as dafon-cv-api
  participant R as Redis

  C->>API: HTTP Request
  API->>R: Check rate limit (IP-based)
  alt Rate limit OK
    R-->>API: Allow request
    API->>API: Process request
    API-->>C: 200 OK + Response
  else Rate limit exceeded
    R-->>API: Block request
    API-->>C: 429 Too Many Requests
  end
```

---

## 🗂️ Data Model (ERD)

```mermaid
erDiagram
  USERS ||--o{ SESSIONS : has
  USERS ||--o{ CURRICULUMS : owns
  USERS ||--|| CONFIGURATIONS : has
  CURRICULUMS ||--o{ WORKS : includes
  CURRICULUMS ||--o{ EDUCATIONS : includes

  USERS {
    uuid id PK
    string name
    string email
    datetime created_at
    datetime updated_at
  }

  SESSIONS {
    uuid id PK
    uuid user_id FK
    string token
    datetime expires_at
  }

  CONFIGURATIONS {
    uuid id PK
    uuid user_id FK
    string language
    bool newsletter
    bool receive_emails
  }

  CURRICULUMS {
    uuid id PK
    uuid user_id FK
    string full_name
    string email
    string phone
    string intro
    string skills
    string languages
    string courses
    string social_links
  }

  WORKS {
    uuid id PK
    uuid curriculum_id FK
    string position
    string company
    string description
    date start_date
    date end_date
  }

  EDUCATIONS {
    uuid id PK
    uuid curriculum_id FK
    string institution
    string degree
    date start_date
    date end_date
    string description
  }
```

---

## 📚 API Documentation

All endpoints require static token authentication via the `Authorization: Bearer <STATIC_TOKEN>` header.

### Health Check

```http
GET /health                  # Application health status
```

### User Management

```http
POST   /api/v1/user                    # Create user
GET    /api/v1/user/all                # Get all users
GET    /api/v1/user/:id                # Get user by ID
PATCH  /api/v1/user/:id                # Update user
DELETE /api/v1/user/:id                # Delete user
```

### Curriculum Management

```http
POST   /api/v1/curriculums                           # Create curriculum
GET    /api/v1/curriculums/get-all-by-user/:user_id   # Get all curriculums by user
GET    /api/v1/curriculums/:curriculum_id            # Get curriculum by ID
GET    /api/v1/curriculums/get-body/:curriculum_id   # Get curriculum body
DELETE /api/v1/curriculums/:curriculum_id            # Delete curriculum
```

### AI Content Generation

```http
POST /api/v1/generate-intro-ai        # Generate professional introduction
POST /api/v1/generate-courses-ai       # Generate course recommendations
POST /api/v1/generate-academic-ai     # Generate academic content
POST /api/v1/generate-task-ai         # Generate task descriptions
POST /api/v1/generate-skill-ai        # Generate skill recommendations
POST /api/v1/generate-analyze-ai      # Analyze and filter content
POST /api/v1/generate-translation-ai # Translate content
```

### Configuration Management

```http
GET    /api/v1/configuration/:user_id    # Get user configuration
PATCH  /api/v1/configuration/:user_id     # Update configuration
DELETE /api/v1/configuration/:user_id     # Delete configuration
```

### Email Services

```http
POST /api/v1/send-email              # Send email via Resend
```

> **Authentication:** All endpoints require the `Authorization: Bearer <STATIC_TOKEN>` header.

---

## ⚙️ Configuration and Environment

Create a `.env` file in the project root:

### Required Environment Variables

```bash
# OpenAI Configuration
OPENAI_API_KEY=your_openai_api_key_here

# Email Service (Resend)
RESEND_API_KEY=your_resend_api_key_here
MAIL_FROM=your_email@domain.com

# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=dafon_cv
DB_SSL_MODE=disable

# MySQL Root Configuration
MYSQL_ROOT_PASSWORD=your_mysql_root_password
MYSQL_DATABASE=dafon_cv
MYSQL_USER=your_mysql_user
MYSQL_PASSWORD=your_mysql_password

# Redis Configuration
REDIS_HOST=
REDIS_PORT=
REDIS_PASSWORD=
REDIS_DB=

# Rate Limiting Configuration
RATE_LIMIT=
RATE_WINDOW_MINUTES=
AI_RATE_LIMIT=
AI_RATE_WINDOW_MINUTES=

# Application Configuration
BACKEND_APIKEY=your_static_token_here
APP_URL=http://localhost:3000
```

### Optional Environment Variables

```bash
# Server Configuration
PORT=
GIN_MODE=release

# Database Host (for Docker)
DB_HOST=

# Redis Host (for Docker)
REDIS_HOST=
```

> **Security Note:** Never commit `.env` files. They are automatically ignored via `.gitignore`.

---

## 🔒 Security and Hardening

- **Static Token Authentication** - Secure API access with configurable tokens
- **HTTPS in Production** - Terminate TLS at reverse proxy or load balancer
- **CORS Configuration** - Configure origins via `APP_URL` environment variable
- **Input Validation** - Comprehensive validation for emails, phones, and data integrity
- **Container Security** - Distroless non-root base image for minimal attack surface
- **Environment Security** - Never commit secrets; use environment variables
- **Rate Limiting** - Redis-based rate limiting with configurable limits per endpoint
- **IP-based Protection** - Rate limiting by client IP with intelligent detection
- **Secret Rotation** - Regularly rotate API keys and tokens
- **Database Security** - Use SSL connections in production (`DB_SSL_MODE=require`)

---

## 📈 Observability (Logging, Health, Metrics)

- Structured logging via Zap
- Health endpoint: `GET /health`
- Container health checks configured for API and MySQL
- Add metrics (suggestion): Prometheus + Grafana (future enhancement)

---

## 🐳 DevOps and Deployment

### Dockerfile (Multi-stage + Distroless)

- Build with Go Alpine, output static binary
- Final image: Distroless nonroot for minimal attack surface

### docker-compose

- MySQL 8.0 + persistent volume
- Redis 7 Alpine + persistent volume
- API service depends on DB and Redis health, exposes `8080`
- Health checks configured for all services

Run:

```bash
docker compose up -d --build
docker compose logs -f api
```

### CI/CD (Suggestion)

- Lint + test on PRs
- Build image, scan vulnerabilities (Trivy/Grype)
- Push to registry, deploy (Railway/Render/Fly.io)

### Deployment Diagram

```mermaid
flowchart LR
  Dev[Developer] --> CI[CI Pipeline]
  CI --> IMG[Container Registry]
  IMG --> PROD[Production Runtime]
  subgraph PROD
    API[API Service]
    DB[(Managed MySQL)]
  end
  API <--> DB
  Users[End Users] --> API
```

---

## 🚀 Performance and Scalability

### Current Optimizations

- **Connection Pooling** - GORM with MySQL connection pooling
- **Distroless Container** - Minimal attack surface and fast startup
- **Structured Logging** - Zap logger for performance monitoring
- **Health Checks** - Container and application health monitoring
- **Rate Limiting** - Redis-based rate limiting with configurable limits
- **Caching Layer** - Redis for rate limiting and future caching needs

### Recommended Enhancements

- **Advanced Caching** - Cache AI responses and frequent data in Redis
- **Pagination** - Implement pagination for list endpoints
- **Async Processing** - Queue system for heavy AI requests
- **Database Indexing** - Optimize queries with proper indexes
- **Load Balancing** - Multiple API instances behind load balancer
- **CDN Integration** - Static assets and API responses caching
- **Rate Limit Analytics** - Monitor and analyze rate limiting patterns

---

## 🧯 Troubleshooting and Runbooks

### Health Checks

```bash
# Application health
curl http://localhost:8080/health

# Container status
docker compose ps

# Application logs
docker compose logs -f api
docker compose logs -f mysql
docker compose logs -f redis
```

### Common Issues & Solutions

| Issue | Symptoms | Solution |
|-------|----------|----------|
| **401 Unauthorized** | Missing/invalid static token | Check `BACKEND_APIKEY` in `.env` |
| **400 Validation Error** | Invalid request data | Verify DTO constraints (email, phone, UUID) |
| **429 Too Many Requests** | Rate limit exceeded | Check rate limiting configuration |
| **OpenAI Errors** | AI generation fails | Check `OPENAI_API_KEY` and quota |
| **Database Connection** | Connection refused | Verify `DB_HOST`, `DB_PORT`, credentials |
| **Redis Connection** | Redis connection failed | Verify `REDIS_HOST`, `REDIS_PORT` |
| **Email Service** | Email sending fails | Check `RESEND_API_KEY` and `MAIL_FROM` |

### Debug Commands

```bash
# Check environment variables
docker compose exec api env | grep -E "(DB_|REDIS_|RATE_|OPENAI_|RESEND_)"

# Test database connection
docker compose exec mysql mysql -u root -p -e "SHOW DATABASES;"

# Test Redis connection
docker compose exec redis redis-cli ping

# View application configuration
docker compose exec api cat /app/.env

# Check rate limiting status
docker compose exec redis redis-cli keys "*rate*"
```

---

## 🚦 Rate Limiting System

The application implements a comprehensive rate limiting system using Redis to protect against abuse and ensure fair usage.

### Rate Limiting Configuration

- **Global Rate Limiting**: Applied to all routes except `/health`
  - Default: 100 requests per minute per IP
  - Configurable via `RATE_LIMIT` and `RATE_WINDOW_MINUTES`

- **AI Endpoints Rate Limiting**: Stricter limits for AI-powered endpoints
  - Default: 10 requests per minute per IP
  - Configurable via `AI_RATE_LIMIT` and `AI_RATE_WINDOW_MINUTES`

### Rate Limiting Features

- **IP-based Detection**: Intelligent IP detection supporting load balancers and proxies
- **Redis Backend**: High-performance rate limiting using Redis
- **Configurable Limits**: Environment-based configuration for different environments
- **Graceful Responses**: JSON error responses with clear messaging
- **Logging**: Comprehensive logging for monitoring and debugging

### Rate Limiting Response

When rate limits are exceeded, the API returns:

```json
{
    "error": "Too Many Requests",
    "message": "Rate limit exceeded. Please try again later."
}
```

**HTTP Status**: `429 Too Many Requests`

### Configuration Examples

```bash
# Development (more permissive)
RATE_LIMIT=1000
AI_RATE_LIMIT=50

# Production (conservative)
RATE_LIMIT=100
AI_RATE_LIMIT=10

# High-traffic (balanced)
RATE_LIMIT=500
AI_RATE_LIMIT=25
```

---

## 🛠️ Installation and Local Setup

### Prerequisites

- **Go 1.24.1+** - [Download Go](https://golang.org/dl/)
- **Docker & Docker Compose** - [Install Docker](https://docs.docker.com/get-docker/)
- **OpenAI API Key** - [Get OpenAI API Key](https://platform.openai.com/api-keys)
- **Resend API Key** - [Get Resend API Key](https://resend.com/api-keys)

### Quick Start with Docker (Recommended)

1. **Clone the repository:**
   ```bash
   git clone https://github.com/Daniel-Fonseca-da-Silva/dafon-cv-api.git
   cd dafon-cv-api
   ```

2. **Create environment file:**
   ```bash
   cp .env.example .env  # If available, or create manually
   # Edit .env with your API keys and configuration
   ```

3. **Start the application:**
   ```bash
   docker compose up -d --build
   ```

4. **Check application status:**
   ```bash
   docker compose logs -f api
   curl http://localhost:8080/health
   ```

### Local Development Setup

1. **Install dependencies:**
   ```bash
   go mod download
   ```

2. **Set up environment variables:**
   ```bash
   # Create .env file with your configuration
   export OPENAI_API_KEY="your_key_here"
   export RESEND_API_KEY="your_key_here"
   # ... other variables
   ```

3. **Run the application:**
   ```bash
   go run cmd/api/main.go
   ```

### Verification

- **Health Check:** `GET http://localhost:8080/health`
- **API Documentation:** Available at `/api/v1/` endpoints
- **Database:** MySQL running on port 3306

---

## 🤝 Contributing

1. Fork the project
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👨‍💻 Author

**Daniel Fonseca da Silva**
- GitHub: [@Daniel-Fonseca-da-Silva](https://github.com/Daniel-Fonseca-da-Silva)

---

⭐ If this project was helpful to you, consider giving it a star!