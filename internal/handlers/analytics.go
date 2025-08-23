package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"linker/internal/database"
	"linker/internal/middleware"
)

type AnalyticsHandler struct {
	db *database.Database
}

func NewAnalyticsHandler(db *database.Database) *AnalyticsHandler {
	return &AnalyticsHandler{db: db}
}

func (h *AnalyticsHandler) GetLinkAnalytics(c *gin.Context) {
	linkID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid link ID"})
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	clicks, err := h.db.GetLinkAnalytics(linkID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve analytics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"link_id": linkID,
		"clicks":  clicks,
		"total":   len(clicks),
	})
}