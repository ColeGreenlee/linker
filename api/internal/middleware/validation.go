package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// FileUploadValidationMiddleware validates file upload requests
func FileUploadValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check Content-Type for multipart/form-data
		contentType := c.GetHeader("Content-Type")
		if !strings.Contains(contentType, "multipart/form-data") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Content-Type must be multipart/form-data",
			})
			c.Abort()
			return
		}

		// Validate short codes if provided
		shortCodes := c.PostFormArray("short_codes")
		for _, shortCode := range shortCodes {
			if !isValidShortCode(shortCode) {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid short code format. Short codes must be 3-32 characters long and contain only letters, numbers, hyphens, and underscores.",
				})
				c.Abort()
				return
			}
		}

		// Validate optional fields
		if title := c.PostForm("title"); title != "" && len(title) > 255 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Title must be less than 255 characters",
			})
			c.Abort()
			return
		}

		if description := c.PostForm("description"); description != "" && len(description) > 1000 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Description must be less than 1000 characters",
			})
			c.Abort()
			return
		}

		if password := c.PostForm("password"); password != "" && len(password) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Password must be at least 6 characters long",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// isValidShortCode checks if a short code meets the requirements
func isValidShortCode(shortCode string) bool {
	if len(shortCode) < 3 || len(shortCode) > 32 {
		return false
	}
	
	for _, r := range shortCode {
		if !((r >= 'a' && r <= 'z') || 
			 (r >= 'A' && r <= 'Z') || 
			 (r >= '0' && r <= '9') || 
			 r == '-' || r == '_') {
			return false
		}
	}
	
	return true
}