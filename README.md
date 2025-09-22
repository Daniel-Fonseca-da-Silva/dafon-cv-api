# 📄 dafon-cv-api

[![Go Version](https://img.shields.io/badge/Go-1.24.1-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](docker-compose.yml)

A robust REST API for CV management with AI features for automatic professional content generation.

## 🎯 Purpose

**dafon-cv-api** is a complete solution for creating and managing professional CVs, offering:

- ✅ Secure email-based authentication system
- ✅ Complete CV CRUD operations
- ✅ Automatic content generation with AI (OpenAI)
- ✅ Robust data validation
- ✅ Well-documented REST API
- ✅ Clean and scalable architecture

## 🏗️ Architecture

The project follows **Clean Architecture** principles with clear separation of concerns:

```
├── cmd/api/           # Application entry point
├── internal/
│   ├── config/        # Application configuration
│   ├── database/      # Database configuration
│   ├── dto/           # Data Transfer Objects
│   ├── errors/        # Custom error handling
│   ├── handlers/      # Presentation layer (HTTP)
│   ├── middleware/    # Middlewares (Auth, CORS, etc.)
│   ├── models/        # Domain entities
│   ├── repositories/  # Data access layer
│   ├── routes/        # Route definitions
│   ├── usecases/      # Business logic
│   └── utils/         # Utilities and validations
```

## 🚀 Technologies

### Backend
- **Go 1.24.1** - Main language
- **Gin** - Web framework
- **GORM** - Database ORM
- **MySQL 8.0** - Database
- **JWT** - Authentication
- **Zap** - Structured logging

### AI and Integrations
- **OpenAI GPT-4o-mini** - Content generation
- **Resend** - Email sending
- **Google UUID** - Unique ID generation

### DevOps
- **Docker & Docker Compose** - Containerization
- **Multi-stage builds** - Image optimization
- **Health checks** - Health monitoring

## 📋 Features

### 🔐 Authentication
- **Passwordless login**: Email-based authentication with temporary tokens
- **User registration**: Account creation with validation
- **JWT**: Secure tokens for sessions
- **Logout**: Session invalidation

### 📄 CV Management
- **Complete CRUD**: Create, read, update and delete CVs
- **Robust validation**: Data validated with specialized libraries
- **Relationships**: Users, CVs, works and configurations

### 🤖 AI Content Generation
- **Professional introductions**: Automatic presentation generation
- **Course lists**: Course and certification suggestions
- **Task descriptions**: Professional description improvements
- **Multiple languages**: Support for Portuguese, English and Spanish

### ⚙️ Configuration
- **User profile**: Customizable settings
- **Preferences**: Application settings
- **Validation**: Real-time data validation

## 🛠️ Installation and Setup

### Prerequisites
- Go 1.24.1+
- Docker and Docker Compose
- MySQL 8.0+ (or use Docker)
- OpenAI API key
- Resend API key

### 1. Clone the repository
```bash
git clone https://github.com/Daniel-Fonseca-da-Silva/dafon-cv-api.git
cd dafon-cv-api
```

### 2. Configure environment variables
Create a `.env` file in the project root with the following variables:

**Required variables:**
- `OPENAI_API_KEY` - OpenAI API key
- `RESEND_API_KEY` - Resend API key for emails  
- `JWT_SECRET` - JWT secret key (generate a secure string)
- `DB_PASSWORD` - Database password
- `MYSQL_ROOT_PASSWORD` - MySQL root password

**Optional variables:**
- `PORT` - Server port (default: 8080)
- `GIN_MODE` - Gin mode (default: release)
- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 3306)
- `DB_NAME` - Database name (default: dafon_cv)
- `JWT_DURATION` - JWT duration (default: 24h)
- `SESSION_DURATION` - Session duration (default: 1h)
- `APP_URL` - Application URL (default: http://localhost:8080)

> ⚠️ **Important**: Never commit the `.env` file to the repository. It's already in `.gitignore`.

### 3. Run with Docker (Recommended)
```bash
# Build and run
docker compose up -d --build

# Just run (if already built)
docker compose up -d

# Stop services
docker compose down
```

### 4. Run locally (Development)
```bash
# Install dependencies
go mod download

# Run migrations
go run cmd/api/main.go

# Or run directly
go run cmd/api/main.go
```

## 📚 API Endpoints

### 🔐 Authentication
```http
POST /auth/register          # Register user
POST /auth/login             # Login (sends token via email)
GET  /auth/login-with-token  # Login with token
POST /auth/logout            # Logout
```

### 👤 Users
```http
GET    /users/profile        # Get user profile
PUT    /users/profile        # Update profile
DELETE /users/account        # Delete account
```

### 📄 CVs
```http
GET    /curriculums          # List CVs
POST   /curriculums          # Create CV
GET    /curriculums/:id      # Get CV
PUT    /curriculums/:id      # Update CV
DELETE /curriculums/:id      # Delete CV
```

### 🤖 AI - Content Generation
```http
POST /ai/generate-intro      # Generate professional introduction
POST /ai/generate-courses    # Generate course list
POST /ai/generate-tasks      # Generate task descriptions
```

### ⚙️ Configuration
```http
GET /configurations          # Get configurations
PUT /configurations          # Update configurations
```

### 🏥 Health Check
```http
GET /health                  # Application status
```

## 🔧 Usage Examples

### Register a user
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Silva",
    "email": "john@example.com"
  }'
```

### Login (will receive token via email)
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com"
  }'
```

### Create a CV
```bash
curl -X POST http://localhost:8080/curriculums \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "full_name": "John Silva",
    "email": "john@example.com",
    "phone": "+5511999999999",
    "intro": "Experienced developer...",
    "technologies": "Go, Python, JavaScript",
    "languages": "Portuguese, English",
    "level_education": "Bachelor Degree"
  }'
```

### Generate content with AI
```bash
curl -X POST http://localhost:8080/ai/generate-intro \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Go developer with 5 years of experience"
  }'
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific tests
go test ./internal/usecases/...
```

## 📊 Monitoring

### Health Check
```bash
curl http://localhost:8080/health
```

### Logs
```bash
# View container logs
docker compose logs -f api

# View database logs
docker compose logs -f mysql
```

## 🔒 Security

- ✅ **JWT** for authentication
- ✅ **CORS** configured
- ✅ **Input validation**
- ✅ **Data sanitization**
- ✅ **Rate limiting** (recommended for production)
- ✅ **HTTPS** (recommended for production)

## 🚀 Deployment

### Production
1. Configure production environment variables
2. Use `GIN_MODE=release`
3. Configure HTTPS
4. Implement rate limiting
5. Configure monitoring
6. Use managed database

### Docker
```bash
# Build for production
docker build -t dafon-cv-api:latest .

# Run in production
docker run -d \
  --name dafon-cv-api \
  -p 8080:8080 \
  --env-file .env.production \
  dafon-cv-api:latest
```

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

## 🙏 Acknowledgments

- [Gin](https://gin-gonic.com/) - Web framework
- [GORM](https://gorm.io/) - ORM
- [OpenAI](https://openai.com/) - AI for content generation
- [Resend](https://resend.com/) - Email service

---

⭐ **If this project was helpful to you, consider giving it a star!**