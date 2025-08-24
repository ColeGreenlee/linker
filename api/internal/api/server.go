package api

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"linker/internal/config"
	"linker/internal/database"
	"linker/internal/handlers"
	"linker/internal/middleware"
	"linker/internal/storage"
)

type Server struct {
	config      *config.Config
	db          *database.Database
	router      *gin.Engine
	rateLimiter *middleware.RateLimiter
}

func NewServer(config *config.Config, db *database.Database) *Server {
	if config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	
	server := &Server{
		config:      config,
		db:          db,
		router:      router,
		rateLimiter: middleware.NewRateLimiter(),
	}
	
	server.setupMiddleware()
	server.setupRoutes()
	
	return server
}

func (s *Server) setupMiddleware() {
	s.router.Use(middleware.CORSMiddleware())
	s.router.Use(gin.Recovery())
	s.router.Use(gin.Logger())
}

func (s *Server) setupRoutes() {
	authHandler := handlers.NewAuthHandler(s.db, s.config.JWTSecret)
	linksHandler := handlers.NewLinksHandler(s.db)
	redirectHandler := handlers.NewRedirectHandler(s.db, s.config.Analytics)
	analyticsHandler := handlers.NewAnalyticsHandler(s.db)
	tokensHandler := handlers.NewTokensHandler(s.db)
	
	// Initialize S3 client if configured
	var s3Client *storage.S3Client
	if s.config.S3.Enabled {
		var err error
		s3Client, err = storage.NewS3Client(&s.config.S3)
		if err != nil {
			log.Printf("Failed to initialize S3 client: %v", err)
			s3Client = nil
		}
	}
	
	filesHandler := handlers.NewFilesHandler(s.db, s3Client, s.config)

	api := s.router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/profile", middleware.AuthMiddleware(s.config.JWTSecret), authHandler.Profile)
		}

		tokens := api.Group("/tokens")
		tokens.Use(middleware.AuthMiddleware(s.config.JWTSecret))
		{
			tokens.POST("", tokensHandler.CreateToken)
			tokens.GET("", tokensHandler.GetTokens)
			tokens.DELETE("/:id", tokensHandler.DeleteToken)
		}

		links := api.Group("/links")
		links.Use(middleware.AuthMiddleware(s.config.JWTSecret))
		{
			links.POST("", linksHandler.CreateLink)
			links.GET("", linksHandler.GetUserLinks)
			links.GET("/:id", linksHandler.GetLink)
			links.PUT("/:id", linksHandler.UpdateLink)
			links.DELETE("/:id", linksHandler.DeleteLink)
		}

		analytics := api.Group("/analytics")
		analytics.Use(middleware.AuthMiddleware(s.config.JWTSecret))
		{
			analytics.GET("/links/:id", analyticsHandler.GetLinkAnalytics)
			analytics.GET("/user", analyticsHandler.GetUserAnalytics)
			analytics.GET("/files", filesHandler.GetUserFileAnalytics)
			analytics.GET("/files/:id/summary", filesHandler.GetFileAnalyticsSummary)
		}

		files := api.Group("/files")
		files.Use(middleware.AuthMiddlewareWithAPITokens(s.config.JWTSecret, s.db))
		{
			files.POST("", 
				s.rateLimiter.FileUploadMiddleware(10, 1*time.Hour), // 10 uploads per hour
				middleware.FileUploadValidationMiddleware(),
				filesHandler.UploadFile)
			files.GET("", filesHandler.GetUserFiles)
			files.GET("/:id", filesHandler.GetFile)
			files.PUT("/:id", filesHandler.UpdateFile)
			files.DELETE("/:id", filesHandler.DeleteFile)
			files.GET("/:id/analytics", filesHandler.GetFileAnalytics)
		}
	}

	// Setup redirect route with configurable prefix
	prefixPattern := fmt.Sprintf("/%s/:shortCode", s.config.LinkPrefix)
	s.router.GET(prefixPattern, redirectHandler.Redirect)
	
	// Setup public file download route with configurable prefix
	filePrefixPattern := fmt.Sprintf("/%s/:shortCode", s.config.FilePrefix)
	s.router.GET(filePrefixPattern, filesHandler.DownloadFile)
	
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.config.Port)
	log.Printf("Server starting on %s", addr)
	return s.router.Run(addr)
}

func (s *Server) Shutdown() {
	if s.rateLimiter != nil {
		s.rateLimiter.Stop()
	}
}