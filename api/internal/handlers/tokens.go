package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"linker/internal/database"
	"linker/internal/middleware"
	"linker/internal/models"
)

type TokensHandler struct {
	db *database.Database
}

func NewTokensHandler(db *database.Database) *TokensHandler {
	return &TokensHandler{db: db}
}

func (h *TokensHandler) CreateToken(c *gin.Context) {
	var req models.CreateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Generate random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	
	tokenString := hex.EncodeToString(tokenBytes)
	
	// Hash the token for storage
	hasher := sha256.New()
	hasher.Write([]byte(tokenString))
	tokenHash := hex.EncodeToString(hasher.Sum(nil))

	apiToken := &models.APIToken{
		UserID:    userID,
		TokenHash: tokenHash,
		Name:      req.Name,
		ExpiresAt: req.ExpiresAt,
	}

	if err := h.db.CreateAPIToken(apiToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	// Return the plain token (only time it's shown)
	response := models.CreateTokenResponse{
		Token:    tokenString,
		APIToken: *apiToken,
	}
	
	c.JSON(http.StatusCreated, response)
}

func (h *TokensHandler) GetTokens(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	tokens, err := h.db.GetUserAPITokens(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tokens": tokens})
}

func (h *TokensHandler) DeleteToken(c *gin.Context) {
	tokenID := c.Param("id")
	if tokenID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token ID"})
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if err := h.db.DeleteAPIToken(tokenID, userID); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Token not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete token"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token deleted successfully"})
}

// ValidateAPIToken validates an API token and returns the associated user ID
func (h *TokensHandler) ValidateAPIToken(tokenString string) (string, error) {
	// Hash the provided token
	hasher := sha256.New()
	hasher.Write([]byte(tokenString))
	tokenHash := hex.EncodeToString(hasher.Sum(nil))

	token, err := h.db.GetAPITokenByHash(tokenHash)
	if err != nil {
		return "", err
	}

	// Check if token is expired
	if token.ExpiresAt != nil && time.Now().After(*token.ExpiresAt) {
		return "", sql.ErrNoRows
	}

	// Update last used timestamp
	h.db.UpdateAPITokenLastUsed(token.ID)

	return token.UserID, nil
}