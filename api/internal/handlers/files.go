package handlers

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"linker/internal/auth"
	"linker/internal/database"
	"linker/internal/middleware"
	"linker/internal/models"
	"linker/internal/storage"
)

type FilesHandler struct {
	db       *database.Database
	s3Client *storage.S3Client
}

func NewFilesHandler(db *database.Database, s3Client *storage.S3Client) *FilesHandler {
	return &FilesHandler{
		db:       db,
		s3Client: s3Client,
	}
}

func (h *FilesHandler) UploadFile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Check if S3 is enabled
	if h.s3Client == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "File upload service is not enabled"})
		return
	}

	// Parse multipart form
	err := c.Request.ParseMultipartForm(100 << 20) // 100MB max
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form"})
		return
	}

	// Get file from form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}
	defer file.Close()

	// Get metadata from form
	var req models.CreateFileRequest
	if title := c.PostForm("title"); title != "" {
		req.Title = title
	}
	if description := c.PostForm("description"); description != "" {
		req.Description = description
	}
	if analytics := c.PostForm("analytics"); analytics == "true" {
		req.Analytics = true
	}
	if isPublic := c.PostForm("is_public"); isPublic == "false" {
		req.IsPublic = false
	} else {
		req.IsPublic = true // Default to public
	}
	if password := c.PostForm("password"); password != "" {
		req.Password = &password
	}
	if domainID := c.PostForm("domain_id"); domainID != "" {
		req.DomainID = &domainID
	}

	// Parse short codes
	shortCodes := c.PostFormArray("short_codes")
	if len(shortCodes) == 0 {
		// Generate a default short code
		shortCodes = []string{generateFileShortCode()}
	}
	req.ShortCodes = shortCodes

	// Check if short codes are already in use
	for _, shortCode := range req.ShortCodes {
		if _, err := h.db.GetLinkByShortCode(shortCode); err != sql.ErrNoRows {
			c.JSON(http.StatusConflict, gin.H{"error": fmt.Sprintf("Short code '%s' already exists", shortCode)})
			return
		}
		if _, err := h.db.GetFileByShortCode(shortCode); err != sql.ErrNoRows {
			c.JSON(http.StatusConflict, gin.H{"error": fmt.Sprintf("Short code '%s' already exists", shortCode)})
			return
		}
	}

	// Determine MIME type
	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = storage.GetMimeTypeFromExtension(storage.GetFileExtension(header.Filename))
	}

	// Upload to S3
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	uploadResult, err := h.s3Client.Upload(ctx, header.Filename, file, mimeType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Upload failed: %v", err)})
		return
	}

	// Hash password if provided
	var hashedPassword *string
	if req.Password != nil {
		hashed, err := auth.HashPassword(*req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		hashedPassword = &hashed
	}

	// Create file record
	fileRecord := &models.File{
		UserID:       userID,
		DomainID:     req.DomainID,
		Filename:     generateUniqueFilename(header.Filename),
		OriginalName: header.Filename,
		MimeType:     mimeType,
		FileSize:     uploadResult.Size,
		S3Key:        uploadResult.Key,
		S3Bucket:     "linker-files", // TODO: get from config
		Title:        req.Title,
		Description:  req.Description,
		Analytics:    req.Analytics,
		IsPublic:     req.IsPublic,
		Password:     hashedPassword,
		ExpiresAt:    req.ExpiresAt,
	}

	err = h.db.CreateFile(fileRecord)
	if err != nil {
		// Cleanup uploaded file on database error
		h.s3Client.Delete(ctx, uploadResult.Key)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file record"})
		return
	}

	// Create short codes
	for i, shortCode := range req.ShortCodes {
		isPrimary := i == 0
		err := h.db.CreateFileShortCode(fileRecord.ID, shortCode, isPrimary)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create short code"})
			return
		}
	}

	// Load short codes back into file for response
	shortCodeRecords, err := h.db.GetShortCodesByFileID(fileRecord.ID)
	if err == nil {
		fileRecord.ShortCodes = shortCodeRecords
	}

	c.JSON(http.StatusCreated, fileRecord)
}

func (h *FilesHandler) GetUserFiles(c *gin.Context) {
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

	files, err := h.db.GetUserFiles(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve files"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": files})
}

func (h *FilesHandler) GetFile(c *gin.Context) {
	fileID := c.Param("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	file, err := h.db.GetFileByID(fileID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve file"})
		}
		return
	}

	c.JSON(http.StatusOK, file)
}

func (h *FilesHandler) UpdateFile(c *gin.Context) {
	fileID := c.Param("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.UpdateFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password if provided
	if req.Password != nil {
		hashed, err := auth.HashPassword(*req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		req.Password = &hashed
	}

	err := h.db.UpdateFile(fileID, userID, &req)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update file"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File updated successfully"})
}

func (h *FilesHandler) DeleteFile(c *gin.Context) {
	fileID := c.Param("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get file info before deletion for S3 cleanup
	file, err := h.db.GetFileByID(fileID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve file"})
		}
		return
	}

	// Delete from database
	err = h.db.DeleteFile(fileID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	// Delete from S3
	if h.s3Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		h.s3Client.Delete(ctx, file.S3Key) // Don't fail if S3 deletion fails
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}

func (h *FilesHandler) DownloadFile(c *gin.Context) {
	shortCode := c.Param("shortCode")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Short code required"})
		return
	}

	// Check password if provided
	password := c.Query("password")

	file, err := h.db.GetFileByShortCode(shortCode)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Check if file is expired
	if file.ExpiresAt != nil && time.Now().After(*file.ExpiresAt) {
		c.JSON(http.StatusGone, gin.H{"error": "File has expired"})
		return
	}

	// Check if file is public or password protected
	if !file.IsPublic {
		if file.Password == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "File is private"})
			return
		}
		
		if password == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Password required", 
				"password_required": true,
			})
			return
		}

		if !auth.CheckPassword(password, *file.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			return
		}
	}

	// Check if we should return file info only
	if c.Query("info") == "true" {
		c.JSON(http.StatusOK, gin.H{
			"id":            file.ID,
			"filename":      file.Filename,
			"original_name": file.OriginalName,
			"mime_type":     file.MimeType,
			"file_size":     file.FileSize,
			"title":         file.Title,
			"description":   file.Description,
			"downloads":     file.Downloads,
			"created_at":    file.CreatedAt,
		})
		return
	}

	// Track download
	if err := h.db.IncrementFileDownloads(file.ID); err != nil {
		// Log error but don't fail the download
	}

	if file.Analytics {
		download := &models.FileDownload{
			FileID:    file.ID,
			IPAddress: getClientIP(c),
			UserAgent: c.GetHeader("User-Agent"),
			Referer:   c.GetHeader("Referer"),
		}
		
		if err := h.db.CreateFileDownload(download); err != nil {
			// Log error but don't fail the download
		}
	}

	// Stream file from S3
	if h.s3Client == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "File download service is not available"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	reader, err := h.s3Client.Download(ctx, file.S3Key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download file"})
		return
	}
	defer reader.Close()

	// Set headers for file download
	c.Header("Content-Type", file.MimeType)
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, file.OriginalName))
	c.Header("Content-Length", fmt.Sprintf("%d", file.FileSize))

	// Stream the file
	io.Copy(c.Writer, reader)
}

func (h *FilesHandler) GetFileAnalytics(c *gin.Context) {
	fileID := c.Param("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	downloads, err := h.db.GetFileAnalytics(fileID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve analytics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"file_id":   fileID,
		"downloads": downloads,
		"total":     len(downloads),
	})
}

func (h *FilesHandler) GetUserFileAnalytics(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	analytics, err := h.db.GetUserFileAnalytics(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user file analytics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"files":   analytics,
		"total":   len(analytics),
	})
}

func (h *FilesHandler) GetFileAnalyticsSummary(c *gin.Context) {
	fileID := c.Param("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	summary, err := h.db.GetFileAnalyticsSummary(fileID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve analytics summary"})
		}
		return
	}

	c.JSON(http.StatusOK, summary)
}

// Helper functions
func generateFileShortCode() string {
	bytes := make([]byte, 6)
	rand.Read(bytes)
	return "f-" + hex.EncodeToString(bytes)[:8] // "f-" prefix for files
}

func generateUniqueFilename(originalFilename string) string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes) + "_" + originalFilename
}

func getClientIP(c *gin.Context) string {
	// This is the same logic as in redirect.go - could be extracted to a shared utility
	clientIP := c.ClientIP()
	
	if forwarded := c.GetHeader("X-Forwarded-For"); forwarded != "" {
		return forwarded
	}
	
	if realIP := c.GetHeader("X-Real-IP"); realIP != "" {
		return realIP
	}
	
	return clientIP
}