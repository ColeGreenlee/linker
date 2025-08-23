package api

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"linker/internal/config"
	"linker/internal/database"
	"linker/internal/handlers"
	"linker/internal/middleware"
)

type Server struct {
	config   *config.Config
	db       *database.Database
	router   *gin.Engine
}

func NewServer(config *config.Config, db *database.Database) *Server {
	if config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	
	server := &Server{
		config: config,
		db:     db,
		router: router,
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

	api := s.router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/profile", middleware.AuthMiddleware(s.config.JWTSecret), authHandler.Profile)
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
		}
	}

	s.router.GET("/:shortCode", redirectHandler.Redirect)
	
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.config.Port)
	log.Printf("Server starting on %s", addr)
	return s.router.Run(addr)
}