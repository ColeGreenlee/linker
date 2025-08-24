# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

This is a monorepo for Linker, a URL shortener and file sharing platform with Go API backend and vanilla JavaScript frontend.

## Architecture

### Monorepo Structure
```
api/          # Go backend (Gin framework, SQLite, S3/MinIO)
ui/           # Vanilla JavaScript frontend (served by Nginx)
```

### Core Components

**API Backend (`api/`):**
- **Framework**: Gin HTTP framework with middleware for CORS, auth, rate limiting
- **Database**: SQLite with automatic migrations, foreign key constraints enabled
- **Storage**: S3/MinIO integration for file uploads with configurable backends
- **Auth**: JWT tokens with bcrypt password hashing
- **Models**: UUID-based primary keys, supports multiple short codes per resource

**Configuration System:**
- Dual configuration: JSON file (`config.json`) takes precedence over environment variables
- S3 config embedded within main config for MinIO/S3 backend switching

**Database Design:**
- `users`, `links`, `files` tables with UUID primary keys
- `short_codes` table supports both links AND files (polymorphic via link_id/file_id)
- `clicks`/`file_downloads` tables for analytics with IP, user agent, referrer tracking
- Migration system with `schema_migrations` table

**Security Features:**
- Rate limiting middleware for uploads and API calls
- Password protection for individual files/links
- Input validation middleware
- File type restrictions via MIME type filtering

### Development Commands

**API Development:**
```bash
cd api
go mod tidy
go run main.go
```

**Run API Tests:**
```bash
cd api
go test -v ./tests                    # All tests
go test -v ./tests/file_sharing_test.go  # Specific test file
```

**Docker Development Environment:**
```bash
docker-compose -f docker-compose.dev.yml up -d    # Start all services (API, UI, MinIO)
docker-compose -f docker-compose.dev.yml build    # Rebuild containers
```

**Production Deployment:**
```bash
docker-compose up -d    # Uses pre-built GHCR images
```

**UI Development:**
```bash
cd ui/src
python -m http.server 8000    # Or any static server
```

## Key Implementation Details

### Configuration Loading
The config system loads from `config.json` first, then falls back to environment variables. S3 configuration is embedded and supports both development (MinIO) and production (AWS S3) scenarios.

### Short Code System
Both links and files share the same short code system via the `short_codes` table. Each resource can have multiple aliases, with one marked as primary.

### File Upload Flow
1. Multipart form handling with size/type validation
2. S3/MinIO storage with unique filename generation  
3. Database record creation with short code assignment
4. Analytics tracking on downloads

### Database Migrations
Uses a custom migration system with `schema_migrations` table tracking. Migrations are in `api/migrations/` and run automatically on startup.

### Authentication Flow
JWT tokens with 24-hour expiration. All protected endpoints use the auth middleware to validate tokens and extract user context.

### Testing Strategy
Tests use in-memory SQLite databases and mock S3 clients. File upload tests create temporary files and validate both storage and database operations.

## GitHub Actions

The workflow builds and publishes separate Docker images:
- `ghcr.io/colegreenlee/linker-api:latest`
- `ghcr.io/colegreenlee/linker-ui:latest`

Tests run before building the API image. Both images are referenced in the production docker-compose.yml.