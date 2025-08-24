package middleware

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"linker/internal/auth"
	"linker/internal/database"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		claims, err := auth.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

func AuthMiddlewareWithAPITokens(jwtSecret string, db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		// First try JWT token
		if claims, err := auth.ValidateToken(tokenString, jwtSecret); err == nil {
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("token_type", "jwt")
			c.Next()
			return
		}

		// Then try API token
		userID, err := validateAPIToken(db, tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Get user info
		user, err := db.GetUserByID(userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("username", user.Username)
		c.Set("token_type", "api")
		c.Next()
	}
}

func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}
	return userID.(string), true
}

func validateAPIToken(db *database.Database, tokenString string) (string, error) {
	// Hash the provided token
	hasher := sha256.New()
	hasher.Write([]byte(tokenString))
	tokenHash := hex.EncodeToString(hasher.Sum(nil))

	token, err := db.GetAPITokenByHash(tokenHash)
	if err != nil {
		return "", err
	}

	// Check if token is expired
	if token.ExpiresAt != nil && time.Now().After(*token.ExpiresAt) {
		return "", sql.ErrNoRows
	}

	// Update last used timestamp
	db.UpdateAPITokenLastUsed(token.ID)

	return token.UserID, nil
}