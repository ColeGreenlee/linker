package handlers

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"linker/internal/database"
	"linker/internal/models"
)

type RedirectHandler struct {
	db        *database.Database
	analytics bool
}

func NewRedirectHandler(db *database.Database, analytics bool) *RedirectHandler {
	return &RedirectHandler{
		db:        db,
		analytics: analytics,
	}
}

func (h *RedirectHandler) Redirect(c *gin.Context) {
	shortCode := c.Param("shortCode")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Short code required"})
		return
	}

	link, err := h.db.GetLinkByShortCode(shortCode)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	if link.ExpiresAt != nil && time.Now().After(*link.ExpiresAt) {
		c.JSON(http.StatusGone, gin.H{"error": "Link has expired"})
		return
	}

	if err := h.db.IncrementLinkClicks(link.ID); err != nil {
		// Log error but don't fail the redirect
	}

	if h.analytics && link.Analytics {
		click := &models.Click{
			LinkID:    link.ID,
			IPAddress: h.getClientIP(c),
			UserAgent: c.GetHeader("User-Agent"),
			Referer:   c.GetHeader("Referer"),
		}
		
		if err := h.db.CreateClick(click); err != nil {
			// Log error but don't fail the redirect
		}
	}

	c.Redirect(http.StatusFound, link.OriginalURL)
}

func (h *RedirectHandler) getClientIP(c *gin.Context) string {
	clientIP := c.ClientIP()
	
	if forwarded := c.GetHeader("X-Forwarded-For"); forwarded != "" {
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	
	if realIP := c.GetHeader("X-Real-IP"); realIP != "" {
		return realIP
	}
	
	return clientIP
}