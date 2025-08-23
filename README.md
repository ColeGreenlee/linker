# Linker - URL Shortener Service

A lightweight, fast, and scalable URL shortener service built with Go, designed for minimal resource usage and high performance.

## Features

- **Multi-domain support**: Configure multiple domains for short URLs
- **User authentication**: JWT-based authentication with user registration
- **Analytics**: Optional click tracking and analytics (toggleable per link)
- **Link management**: Create, update, delete, and organize your links
- **API-driven**: Complete REST API for all functionality
- **SQLite backend**: Lightweight database with automatic migrations
- **Docker ready**: Optimized for containerized deployment
- **Secure**: Password hashing, JWT tokens, input validation

## Quick Start

### Docker (Recommended)

```bash
# Clone and run with docker-compose
git clone git@github.com:ColeGreenlee/linker.git
cd linker
docker-compose up -d
```

### Local Development

```bash
# Clone repository
git clone git@github.com:ColeGreenlee/linker.git
cd linker

# Install dependencies
go mod tidy

# Option 1: Use JSON config (recommended)
cp config.example.json config.json
# Edit config.json with your settings
go run .

# Option 2: Use environment variables
export JWT_SECRET=your-secret-key
export DEFAULT_DOMAIN=localhost:8080
go run .
```

## API Endpoints

### Authentication

- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login user
- `GET /api/v1/auth/profile` - Get user profile (requires auth)

### Links

- `POST /api/v1/links` - Create new short link
- `GET /api/v1/links` - Get user's links (paginated)
- `GET /api/v1/links/:id` - Get specific link
- `PUT /api/v1/links/:id` - Update link
- `DELETE /api/v1/links/:id` - Delete link

### Analytics

- `GET /api/v1/analytics/links/:id` - Get link click analytics

### API Tokens

- `POST /api/v1/tokens` - Create new API token (requires auth)
- `GET /api/v1/tokens` - List user's API tokens (requires auth)
- `DELETE /api/v1/tokens/:id` - Delete API token (requires auth)

### Redirect

- `GET /:shortCode` - Redirect to original URL

### Health

- `GET /health` - Health check endpoint

## Configuration

Linker supports two configuration methods: **JSON file** (recommended) or **environment variables**.

### JSON Configuration (Recommended)

Create a `config.json` file in your project directory:

```json
{
  "port": "8080",
  "database_url": "./linker.db",
  "default_domain": "localhost:8080",
  "allowed_domains": [
    "localhost:8080",
    "short.local",
    "yourdomain.com"
  ],
  "jwt_secret": "your-production-secret-key-change-this",
  "analytics": true,
  "environment": "development"
}
```

You can also specify a custom config file location:
```bash
export CONFIG_FILE=/path/to/your/config.json
./linker
```

### Environment Variables

If no `config.json` is found, Linker will fall back to environment variables:

```bash
PORT=8080                                    # Server port
DATABASE_URL=./linker.db                     # SQLite database path
DEFAULT_DOMAIN=localhost:8080                # Default domain for short URLs
ALLOWED_DOMAINS=domain1.com,domain2.com      # Additional allowed domains (comma-separated)
JWT_SECRET=your-secret-key                   # JWT signing secret
ANALYTICS=true                               # Enable analytics by default
ENVIRONMENT=development                      # Environment (development/production)
CONFIG_FILE=config.json                      # Custom config file path (optional)
```

### Configuration Priority

1. **JSON file** (if exists) - highest priority
2. **Environment variables** - fallback
3. **Default values** - final fallback

## API Usage Examples

### Register User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \\
  -H "Content-Type: application/json" \\
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "securepassword"
  }'
```

### Create Short Link

```bash
curl -X POST http://localhost:8080/api/v1/links \\
  -H "Content-Type: application/json" \\
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \\
  -d '{
    "original_url": "https://example.com/very/long/url",
    "short_code": "custom",
    "title": "My Link",
    "analytics": true
  }'
```

### Create API Token

```bash
curl -X POST http://localhost:8080/api/v1/tokens \\
  -H "Content-Type: application/json" \\
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \\
  -d '{
    "name": "My API Token",
    "expires_at": "2024-12-31T23:59:59Z"
  }'
```

### Use Short Link

```bash
curl http://localhost:8080/custom
# Redirects to https://example.com/very/long/url
```

## Deployment

### Production Docker

```dockerfile
# Build optimized image
docker build -t linker .

# Run with persistent data
docker run -d \\
  -p 8080:8080 \\
  -v linker_data:/app/data \\
  -e JWT_SECRET=your-production-secret \\
  -e DEFAULT_DOMAIN=yourdomain.com \\
  linker
```

### GitHub Container Registry

This project automatically builds and publishes to GitHub Container Registry on push to main:

```bash
docker pull ghcr.io/colegreenlee/linker:main
```

## Development

### Running Tests

```bash
# Run unit tests
go test ./tests/auth_test.go ./tests/config_test.go

# Note: Database tests require CGO for SQLite
```

### Building

```bash
# Build binary
go build -o linker .

# Build with static linking (for Docker)
CGO_ENABLED=1 go build -a -ldflags '-linkmode external -extldflags "-static"' -o linker .
```

## Architecture

- **Gin Framework**: Lightweight HTTP router and middleware
- **SQLite**: Embedded database with foreign keys enabled
- **JWT**: Stateless authentication
- **Bcrypt**: Password hashing
- **Multi-stage Docker**: Minimal production image

## Performance

Designed for lightweight deployment:
- Alpine-based Docker image (~20MB)
- SQLite for minimal memory usage
- Efficient Go routines for concurrent requests
- Static binary with no external dependencies

## Security

- Passwords hashed with bcrypt
- JWT tokens for stateless auth
- Input validation and sanitization
- CORS headers configured
- Non-root user in Docker container

## License

MIT License - see LICENSE file for details.