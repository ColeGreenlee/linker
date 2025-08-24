# Linker - URL Shortener & File Sharing Platform

A modern, full-stack URL shortener and file sharing platform built with Go (API) and vanilla JavaScript (UI). Features custom short URLs, secure file sharing, analytics, and ShareX integration.

## 🚀 Quick Start

### Development Environment

```bash
# Start all services (API, UI, MinIO)
docker-compose -f docker-compose.dev.yml up -d

# Access the application
open http://localhost:3000
```

### Production Deployment

```bash
# Deploy with pre-built images from GHCR
docker-compose up -d

# Access the application
open http://localhost
```

## 🛠️ Development

### API Development
```bash
cd api
go mod tidy
go run main.go
```

### UI Development
The UI is a vanilla JavaScript SPA that can be served with any static web server:
```bash
cd ui/src
python -m http.server 8000
# or
npx serve .
```

### Testing
```bash
# Run API tests
cd api && go test -v ./tests

# Build and test containers
docker-compose -f docker-compose.dev.yml build
```

## ✨ Features

### Core Functionality
- **URL Shortening**: Create custom short URLs with analytics
- **File Sharing**: Secure file uploads with S3/MinIO storage  
- **User Authentication**: JWT-based auth with registration/login
- **Analytics Dashboard**: Track clicks, downloads, and usage statistics

### Advanced Features
- **Custom Short Codes**: Create personalized short URLs and file codes
- **Multiple Short Codes**: Each link/file supports multiple aliases
- **Unified URL Prefixes**: Configurable prefixes for links and files (`/s/`, `/go/`, etc.)
- **API Token Authentication**: Secure programmatic access with Bearer tokens
- **Password Protection**: Secure links and files with passwords
- **File Expiration**: Set automatic expiration dates
- **ShareX Integration**: Built-in support for ShareX screenshot uploads
- **Rate Limiting**: Built-in upload and API rate limiting
- **Multi-domain Support**: Configure multiple domains
- **Real-time Analytics**: Track clicks, downloads, and usage statistics
- **Mobile Responsive**: Clean, modern UI that works on all devices

## 🔧 Configuration

### Configuration File (Recommended)

Create a `config.json` file in the root directory:

```json
{
  "port": "8080",
  "database_url": "./linker.db",
  "default_domain": "yourdomain.com",
  "allowed_domains": [
    "yourdomain.com",
    "short.yourdomain.com"
  ],
  "unified_prefix": "s",
  "jwt_secret": "your-production-secret-key-change-this",
  "analytics": true,
  "environment": "production",
  "s3": {
    "enabled": true,
    "endpoint": "minio:9000",
    "access_key_id": "minioadmin",
    "secret_access_key": "minioadmin",
    "bucket_name": "linker-files",
    "max_file_size_mb": 100
  }
}
```

### URL Prefix Configuration

Linker supports flexible URL prefix configuration:

#### Unified Prefix (Recommended)
```json
{
  "unified_prefix": "s"
}
```
**Result**: Both links and files use `/s/` prefix
- Links: `https://yourdomain.com/s/my-link`
- Files: `https://yourdomain.com/s/my-file`

#### Separate Prefixes
```json
{
  "link_prefix": "go",
  "file_prefix": "dl"
}
```
**Result**: Different prefixes for links and files
- Links: `https://yourdomain.com/go/my-link`
- Files: `https://yourdomain.com/dl/my-file`

#### Mixed Configuration
```json
{
  "unified_prefix": "s",
  "file_prefix": "files"
}
```
**Result**: Override files while keeping links unified
- Links: `https://yourdomain.com/s/my-link`
- Files: `https://yourdomain.com/files/my-file`

### Environment Variables (Alternative)

**API Configuration:**
```env
PORT=8080
DATABASE_URL=/app/data/linker.db
DEFAULT_DOMAIN=yourdomain.com
UNIFIED_PREFIX=s
LINK_PREFIX=s
FILE_PREFIX=f
JWT_SECRET=your-secret-key
ANALYTICS=true
ENVIRONMENT=production

# S3/MinIO Configuration
S3_ENABLED=true
S3_ENDPOINT=minio:9000
S3_ACCESS_KEY_ID=minioadmin
S3_SECRET_ACCESS_KEY=minioadmin
S3_BUCKET_NAME=linker-files
S3_MAX_FILE_SIZE_MB=100
```

### Docker Compose Files

- **`docker-compose.dev.yml`**: Development environment with building
- **`docker-compose.yml`**: Production deployment with GHCR images

## 📤 ShareX Integration

Linker includes built-in ShareX support for seamless screenshot and file uploads.

### Setup Instructions

1. **Create an API token**:
   - Login to your Linker instance
   - Navigate to the API Tokens section
   - Create a new token and copy it

2. **Configure ShareX**:
   - Use the included `sharex-config.sxcu` file as a template
   - Edit the file and replace `YOUR_API_TOKEN_HERE` with your actual API token
   - Update the `RequestURL` to match your Linker instance URL
   - Import the configuration in ShareX: `Destinations` → `Custom uploader settings` → `Import` → `From file`

3. **Set as default uploader**:
   - Go to `Destinations` in ShareX
   - Select your imported Linker uploader for Image uploader and/or File uploader

### ShareX Configuration Template

```json
{
  "Version": "18.0.1",
  "Name": "Linker",
  "DestinationType": "ImageUploader, FileUploader",
  "RequestMethod": "POST",
  "RequestURL": "https://yourdomain.com/api/v1/files",
  "Headers": {
    "Authorization": "Bearer YOUR_API_TOKEN_HERE"
  },
  "Body": "MultipartFormData",
  "Arguments": {
    "analytics": "true",
    "public": "true"
  },
  "FileFormName": "file",
  "URL": "{json:url}",
  "ThumbnailURL": "{json:data.thumb}",
  "DeletionURL": "{json:data.delete_url}",
  "ErrorMessage": "{json:error}"
}
```

### Testing Your Configuration

Test the API manually with curl:
```bash
curl -X POST \
  -H "Authorization: Bearer YOUR_API_TOKEN" \
  -F "file=@test-file.png" \
  -F "analytics=true" \
  -F "public=true" \
  https://yourdomain.com/api/v1/files
```

### Advanced ShareX Features

You can customize the ShareX configuration for additional features:

**Password Protection:**
```json
"Arguments": {
  "analytics": "true",
  "public": "true",
  "password": "your-password"
}
```

**Custom Expiration:**
```json
"Arguments": {
  "analytics": "true",
  "public": "true",
  "expires_at": "2024-12-31T23:59:59Z"
}
```

**Custom Short Code:**
```json
"Arguments": {
  "analytics": "true",
  "public": "true",
  "short_code": "my-custom-code"
}
```

## 🐳 Container Images

Pre-built images are available on GitHub Container Registry:

- **API**: `ghcr.io/colegreenlee/linker-api:latest`
- **UI**: `ghcr.io/colegreenlee/linker-ui:latest`

## 🚀 Deployment Options

### Local Development
```bash
# Clone the repository
git clone https://github.com/colegreenlee/linker.git
cd linker

# Start all services (API, UI, MinIO)
docker-compose -f docker-compose.dev.yml up -d

# Access the application
open http://localhost:3000
```

### Production (Docker)
```bash
# Copy and configure settings
cp config.example.json config.json
# Edit config.json with your production settings

# Deploy with pre-built images
docker-compose up -d

# Access the application
open http://localhost
```

### Standalone API Development
```bash
cd api
go mod tidy
go run main.go
```

### Standalone UI Development
```bash
cd ui/src
python -m http.server 8000
```

---

## 📡 API Reference

### Base URL
All API endpoints are prefixed with `/api/v1/`

### Authentication

#### Register User
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "string (required, 3-50 chars)",
  "email": "string (required, valid email)",
  "password": "string (required, min 6 chars)"
}
```

Returns: `AuthResponse` with JWT token and user data

#### Login User
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "string (required)",
  "password": "string (required)"
}
```

Returns: `AuthResponse` with JWT token and user data

#### Get Profile
```http
GET /api/v1/auth/profile
Authorization: Bearer <token>
```

Returns: `User` object

---

### Links Management

#### Create Short Link
```http
POST /api/v1/links
Authorization: Bearer <token>
Content-Type: application/json

{
  "original_url": "string (required, valid URL)",
  "short_codes": ["string"] (optional),
  "domain_id": "string" (optional),
  "title": "string" (optional),
  "description": "string" (optional),
  "analytics": boolean (default: true),
  "expires_at": "ISO8601 datetime" (optional)
}
```

Returns: `Link` object with generated short codes

#### Get User Links
```http
GET /api/v1/links
Authorization: Bearer <token>
```

Returns: Array of `Link` objects

#### Get Specific Link
```http
GET /api/v1/links/:id
Authorization: Bearer <token>
```

Returns: `Link` object

#### Update Link
```http
PUT /api/v1/links/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "original_url": "string" (optional, valid URL),
  "title": "string" (optional),
  "description": "string" (optional),
  "analytics": boolean,
  "expires_at": "ISO8601 datetime" (optional)
}
```

Returns: Updated `Link` object

#### Delete Link
```http
DELETE /api/v1/links/:id
Authorization: Bearer <token>
```

Returns: Success message

---

### File Management

#### Upload File
```http
POST /api/v1/files
Authorization: Bearer <token>
Content-Type: multipart/form-data

file: <file data>
short_codes: ["string"] (optional)
domain_id: "string" (optional)
title: "string" (optional)
description: "string" (optional)
analytics: boolean (default: true)
is_public: boolean (default: true)
password: "string" (optional)
expires_at: "ISO8601 datetime" (optional)
```

Returns: `File` object with download URL

#### Get User Files
```http
GET /api/v1/files
Authorization: Bearer <token>
```

Returns: Array of `File` objects

#### Get Specific File
```http
GET /api/v1/files/:id
Authorization: Bearer <token>
```

Returns: `File` object

#### Update File Metadata
```http
PUT /api/v1/files/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "string" (optional),
  "description": "string" (optional),
  "analytics": boolean,
  "is_public": boolean,
  "password": "string" (optional),
  "expires_at": "ISO8601 datetime" (optional)
}
```

Returns: Updated `File` object

#### Delete File
```http
DELETE /api/v1/files/:id
Authorization: Bearer <token>
```

Returns: Success message

---

### Public Access

#### Access Short Link/File
```http
GET /{prefix}/:shortCode
```

Redirects to original URL (for links) or serves file (for files)

#### Get File Info
```http
GET /{prefix}/:shortCode?info=true
```

Returns: File metadata without downloading

#### Access Protected Content
```http
GET /{prefix}/:shortCode?password=SECRET
```

Access password-protected links or files

---

### Analytics

#### Get User Analytics
```http
GET /api/v1/analytics/user
Authorization: Bearer <token>
```

Returns: `UserAnalytics` object with comprehensive statistics

#### Get Link Analytics
```http
GET /api/v1/analytics/links/:id
Authorization: Bearer <token>
```

Returns: Detailed analytics for specific link

#### Get File Analytics
```http
GET /api/v1/analytics/files/:id/summary
Authorization: Bearer <token>
```

Returns: `FileAnalyticsSummary` object

---

### API Tokens

#### Create API Token
```http
POST /api/v1/tokens
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "string" (optional),
  "expires_at": "ISO8601 datetime" (optional)
}
```

Returns: `CreateTokenResponse` with plain-text token (only shown once)

#### List User Tokens
```http
GET /api/v1/tokens
Authorization: Bearer <token>
```

Returns: Array of `APIToken` objects (without token values)

#### Delete API Token
```http
DELETE /api/v1/tokens/:id
Authorization: Bearer <token>
```

Returns: Success message

---

### Data Models

#### User
```json
{
  "id": "string (UUID)",
  "username": "string",
  "email": "string",
  "created_at": "ISO8601 datetime",
  "updated_at": "ISO8601 datetime"
}
```

#### Link
```json
{
  "id": "string (UUID)",
  "user_id": "string (UUID)",
  "domain_id": "string (UUID, optional)",
  "domain": "Domain object (optional)",
  "short_codes": ["ShortCode objects"],
  "original_url": "string",
  "title": "string (optional)",
  "description": "string (optional)",
  "clicks": "integer",
  "analytics": "boolean",
  "expires_at": "ISO8601 datetime (optional)",
  "created_at": "ISO8601 datetime",
  "updated_at": "ISO8601 datetime"
}
```

#### File
```json
{
  "id": "string (UUID)",
  "user_id": "string (UUID)",
  "domain_id": "string (UUID, optional)",
  "domain": "Domain object (optional)",
  "short_codes": ["ShortCode objects"],
  "filename": "string (generated)",
  "original_name": "string (user filename)",
  "mime_type": "string",
  "file_size": "integer (bytes)",
  "title": "string (optional)",
  "description": "string (optional)",
  "downloads": "integer",
  "analytics": "boolean",
  "is_public": "boolean",
  "expires_at": "ISO8601 datetime (optional)",
  "created_at": "ISO8601 datetime",
  "updated_at": "ISO8601 datetime"
}
```

#### ShortCode
```json
{
  "id": "string (UUID)",
  "link_id": "string (UUID, if for link)",
  "file_id": "string (UUID, if for file)",
  "short_code": "string",
  "is_primary": "boolean",
  "created_at": "ISO8601 datetime"
}
```

#### Domain
```json
{
  "id": "string (UUID)",
  "domain": "string",
  "is_default": "boolean",
  "enabled": "boolean",
  "created_at": "ISO8601 datetime",
  "updated_at": "ISO8601 datetime"
}
```

#### APIToken
```json
{
  "id": "string (UUID)",
  "user_id": "string (UUID)",
  "name": "string (optional)",
  "last_used_at": "ISO8601 datetime (optional)",
  "expires_at": "ISO8601 datetime (optional)",
  "created_at": "ISO8601 datetime"
}
```

#### UserAnalytics
```json
{
  "user_id": "string (UUID)",
  "total_links": "integer",
  "total_clicks": "integer",
  "clicks_today": "integer",
  "clicks_this_week": "integer",
  "clicks_this_month": "integer",
  "top_links": ["LinkAnalyticsSummary objects"],
  "recent_clicks": ["Click objects"],
  "clicks_by_date": ["ClicksByDate objects"],
  "top_referrers": ["ReferrerStats objects"],
  "top_countries": ["CountryStats objects"],
  "top_user_agents": ["UserAgentStats objects"]
}
```

#### Click (Analytics)
```json
{
  "id": "string (UUID)",
  "link_id": "string (UUID)",
  "ip_address": "string",
  "user_agent": "string",
  "referer": "string (optional)",
  "country": "string (optional)",
  "created_at": "ISO8601 datetime"
}
```

#### FileDownload (Analytics)
```json
{
  "id": "string (UUID)",
  "file_id": "string (UUID)",
  "ip_address": "string",
  "user_agent": "string",
  "referer": "string (optional)",
  "country": "string (optional)",
  "created_at": "ISO8601 datetime"
}
```

### Error Responses

All API endpoints return consistent error responses:

```json
{
  "error": "string (error message)",
  "code": "string (error code, optional)",
  "details": "object (additional context, optional)"
}
```

**Common HTTP Status Codes:**
- `200` - Success
- `201` - Created
- `400` - Bad Request (validation errors)
- `401` - Unauthorized (invalid/missing token)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found
- `409` - Conflict (duplicate short code, etc.)
- `413` - Payload Too Large (file size exceeded)
- `429` - Too Many Requests (rate limited)
- `500` - Internal Server Error

### Rate Limiting

API endpoints are rate-limited per IP address:
- File uploads: 10 requests per minute
- General API calls: 100 requests per minute

Rate limit headers are included in responses:
- `X-RateLimit-Limit` - Request limit per window
- `X-RateLimit-Remaining` - Requests remaining in window  
- `X-RateLimit-Reset` - Window reset time (Unix timestamp)