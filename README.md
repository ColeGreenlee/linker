# Linker - URL Shortener & File Sharing Platform

A modern, full-stack URL shortener and file sharing platform built with Go (API) and vanilla JavaScript (UI).

## 🏗️ Monorepo Structure

```
linker/
├── api/                    # Go API backend
│   ├── internal/          # Internal Go packages
│   ├── migrations/        # Database migrations
│   ├── tests/            # Go tests
│   ├── main.go           # API entry point
│   ├── go.mod            # Go dependencies
│   └── Dockerfile        # API container
├── ui/                    # Frontend web application  
│   ├── src/              # Source files
│   │   ├── index.html    # Main HTML
│   │   ├── app.js        # JavaScript application
│   │   └── styles.css    # CSS styles
│   ├── nginx.conf        # Nginx configuration
│   └── Dockerfile        # UI container
├── docker-compose.yml     # Production deployment
├── docker-compose.dev.yml # Development environment
└── README.md             # This file
```

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
- **Multiple Short Codes**: Each link/file supports multiple aliases
- **Password Protection**: Secure links and files with passwords
- **File Expiration**: Set automatic expiration dates
- **Rate Limiting**: Built-in upload and API rate limiting
- **Multi-domain Support**: Configure multiple domains
- **Mobile Responsive**: Works on all devices

## 🔧 Configuration

### Environment Variables

**API Configuration:**
```env
PORT=8080
DATABASE_URL=/app/data/linker.db
DEFAULT_DOMAIN=yourdomain.com
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

## 📡 API Reference

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login user
- `GET /api/v1/auth/profile` - Get user profile

### Links
- `POST /api/v1/links` - Create short link
- `GET /api/v1/links` - Get user's links
- `GET /api/v1/links/:id` - Get specific link
- `PUT /api/v1/links/:id` - Update link
- `DELETE /api/v1/links/:id` - Delete link

### Files  
- `POST /api/v1/files` - Upload file
- `GET /api/v1/files` - Get user's files
- `GET /api/v1/files/:id` - Get specific file
- `PUT /api/v1/files/:id` - Update file metadata
- `DELETE /api/v1/files/:id` - Delete file

### Public Access
- `GET /x/:shortCode` - Redirect to URL
- `GET /f/:shortCode` - Download file
- `GET /f/:shortCode?info=true` - Get file info
- `GET /f/:shortCode?password=SECRET` - Access protected file

### Analytics
- `GET /api/v1/analytics/user` - User analytics
- `GET /api/v1/analytics/links/:id` - Link analytics  
- `GET /api/v1/analytics/files/:id/summary` - File analytics

## 🐳 Container Images

Pre-built images are available on GitHub Container Registry:

- **API**: `ghcr.io/colegreenlee/linker-api:latest`
- **UI**: `ghcr.io/colegreenlee/linker-ui:latest`

## 🏛️ Architecture

### Backend (Go API)
- **Framework**: Gin HTTP framework
- **Database**: SQLite with migrations
- **Storage**: S3/MinIO object storage  
- **Auth**: JWT tokens with bcrypt passwords
- **Security**: Rate limiting, input validation, CORS

### Frontend (JavaScript SPA)
- **Framework**: Vanilla JavaScript (no framework)
- **Build**: Static files served by Nginx
- **API**: Fetch API for backend communication
- **Styling**: Pure CSS with responsive design

### Infrastructure
- **Containers**: Docker with multi-stage builds
- **Proxy**: Nginx reverse proxy for API
- **Storage**: MinIO for development, S3 for production
- **Database**: SQLite (embedded) with automatic migrations

## 🔒 Security Features

- JWT token authentication
- Bcrypt password hashing  
- File password protection
- Rate limiting for uploads
- MIME type validation
- File size limits
- Input validation and sanitization
- CORS configuration
- Security headers (nginx)

## 📊 Monitoring & Analytics

- Click tracking with IP, user agent, referrer
- Download analytics with detailed metrics
- User dashboard with statistics
- File and link performance metrics
- Geographic data (when available)

## 🚀 Deployment Options

### Local Development
```bash
docker-compose -f docker-compose.dev.yml up -d
```

### Production (Docker)
```bash
docker-compose up -d
```

### Production (Manual)
1. Build and push images to registry
2. Update `docker-compose.yml` with your registry URLs
3. Deploy with your container orchestration platform

## 📝 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable  
5. Submit a pull request

## 📞 Support

- 🐛 **Issues**: [GitHub Issues](https://github.com/colegreenlee/linker/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/colegreenlee/linker/discussions)
- 📧 **Contact**: [Your Email](mailto:your-email@example.com)