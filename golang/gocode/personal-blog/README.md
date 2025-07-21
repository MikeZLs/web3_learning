# Personal Blog System Backend

A complete backend implementation for a personal blog system built with Go, Gin framework, GORM, and MySQL.

## Features

- **User Authentication**: Registration, login with JWT tokens
- **Blog Posts**: CRUD operations for blog posts
- **Comments**: Users can comment on posts
- **Authorization**: Only post authors can edit/delete their posts
- **Error Handling**: Comprehensive error handling with proper HTTP status codes
- **Logging**: Using Zap logger with file output to E drive

## Tech Stack

- **Go 1.24+**
- **Gin Framework**: HTTP web framework
- **GORM**: ORM library for database operations
- **MySQL**: Version 8.0.26
- **JWT**: Authentication tokens
- **Zap**: Structured, leveled logging
- **bcrypt**: Password hashing

## Installation & Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd personal-blog
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Setup MySQL Database**
    - Create a MySQL database named `blog_db`
    - Update database connection in `config/config.go` if needed

4. **Environment Variables** (Optional)
   ```bash
   export DATABASE_URL="root:root@tcp(localhost:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
   export JWT_SECRET="your-secret-key"
   export PORT="8080"
   export LOG_LEVEL="INFO"  # DEBUG, INFO, WARN, ERROR
   ```

5. **Create log directory**
    - The application will automatically create `E:/logs/` directory
    - Or set custom path with `LOG_PATH` environment variable

6. **Run the application**
   ```bash
   go run main.go
   ```

## Logging Configuration

The application uses Zap logger with the following configuration:

### Log Files Location
- **Main Log**: `E:/logs/blog-app.log` - All application logs
- **Error Log**: `E:/logs/blog-error.log` - Error-level logs only
- **Console Output**: Logs also appear in terminal

### Log Levels
Set via `LOG_LEVEL` environment variable:
- `DEBUG`: Most verbose, includes all logs
- `INFO`: General information (default)
- `WARN`: Warning messages
- `ERROR`: Error messages only

### Log Format
Logs are structured in JSON format with the following fields:
- `timestamp`: ISO8601 formatted time
- `level`: Log level (DEBUG, INFO, WARN, ERROR)
- `logger`: Logger name
- `caller`: Source file and line number
- `message`: Log message
- `stacktrace`: Stack trace for errors

### Example Log Entry
```json
{
  "timestamp": "2024-01-15T10:30:45.123Z",
  "level": "INFO",
  "logger": "personal-blog",
  "caller": "controllers/auth.go:45",
  "message": "User registered successfully",
  "username": "john_doe"
}
```

## API Endpoints

### Authentication
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login

### Posts
- `GET /api/posts` - Get all posts
- `GET /api/posts/:id` - Get single post
- `POST /api/posts` - Create post (authenticated)
- `PUT /api/posts/:id` - Update post (authenticated, owner only)
- `DELETE /api/posts/:id` - Delete post (authenticated, owner only)

### Comments
- `GET /api/posts/:post_id/comments` - Get post comments
- `POST /api/posts/:post_id/comments` - Create comment (authenticated)

### Health Check
- `GET /health` - Service health status

## Database Schema

### Users Table
- `id` (Primary Key)
- `username` (Unique)
- `email` (Unique)
- `password` (Hashed)
- `created_at`
- `updated_at`
- `deleted_at`

### Posts Table
- `id` (Primary Key)
- `title`
- `content`
- `user_id` (Foreign Key to Users)
- `created_at`
- `updated_at`
- `deleted_at`

### Comments Table
- `id` (Primary Key)
- `content`
- `user_id` (Foreign Key to Users)
- `post_id` (Foreign Key to Posts)
- `created_at`

## Project Structure

```
personal-blog/
├── main.go                 # Application entry point
├── go.mod                  # Go module dependencies
├── README.md              # Project documentation
│
├── config/                # Configuration management
│   └── config.go          # Database, JWT, port configurations
│
├── models/                # Data models
│   ├── user.go           # User model and structs
│   ├── post.go           # Post model and structs
│   └── comment.go        # Comment model and structs
│
├── database/              # Database related
│   └── database.go       # Database initialization and migration
│
├── middleware/            # Middleware
│   └── auth.go           # JWT authentication middleware
│
├── controllers/           # Controllers (business logic)
│   ├── auth.go           # Authentication controller
│   ├── post.go           # Post management controller
│   └── comment.go        # Comment management controller
│
└── routes/                # Route configuration
    └── routes.go         # API routes setup

Logs will be stored in:
E:/logs/
├── blog-app.log          # Main application logs
└── blog-error.log        # Error logs only
```

## Authentication

The API uses JWT (JSON Web Tokens) for authentication:

1. Register or login to get a JWT token
2. Include the token in the Authorization header: `Bearer <token>`
3. Tokens expire after 24 hours

## Error Handling

The API returns standard HTTP status codes:
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `409` - Conflict
- `500` - Internal Server Error

## Development

### Running in Development Mode
```bash
# Set Gin to debug mode
export GIN_MODE=debug

# Set log level to debug
export LOG_LEVEL=DEBUG

go run main.go
```

### Testing the API
You can use tools like Postman, curl, or any HTTP client to test the endpoints.

Example registration:
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'
```