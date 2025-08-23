package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"linker/internal/database"
	"linker/internal/middleware"
	"linker/internal/models"
)

type LinksHandler struct {
	db *database.Database
}

func NewLinksHandler(db *database.Database) *LinksHandler {
	return &LinksHandler{db: db}
}

func (h *LinksHandler) CreateLink(c *gin.Context) {
	var req models.CreateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if req.ShortCode == "" {
		req.ShortCode = generateShortCode()
	}

	if _, err := h.db.GetLinkByShortCode(req.ShortCode); err != sql.ErrNoRows {
		c.JSON(http.StatusConflict, gin.H{"error": "Short code already exists"})
		return
	}

	link := &models.Link{
		UserID:      userID,
		DomainID:    req.DomainID,
		ShortCode:   req.ShortCode,
		OriginalURL: req.OriginalURL,
		Title:       req.Title,
		Description: req.Description,
		Analytics:   req.Analytics,
		ExpiresAt:   req.ExpiresAt,
	}

	if err := h.db.CreateLink(link); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create link"})
		return
	}

	c.JSON(http.StatusCreated, link)
}

func (h *LinksHandler) GetUserLinks(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	limit := 20
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	links, err := h.db.GetUserLinks(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve links"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"links": links})
}

func (h *LinksHandler) GetLink(c *gin.Context) {
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

	links, err := h.db.GetUserLinks(userID, 1, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve link"})
		return
	}

	var link *models.Link
	for _, l := range links {
		if l.ID == linkID {
			link = &l
			break
		}
	}

	if link == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		return
	}

	c.JSON(http.StatusOK, link)
}

func (h *LinksHandler) UpdateLink(c *gin.Context) {
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

	var req models.UpdateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.UpdateLink(linkID, userID, &req); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update link"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Link updated successfully"})
}

func (h *LinksHandler) DeleteLink(c *gin.Context) {
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

	if err := h.db.DeleteLink(linkID, userID); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete link"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Link deleted successfully"})
}

func generateShortCode() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}