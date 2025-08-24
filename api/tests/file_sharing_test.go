package tests

import (
	"bytes"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"
	"time"

	"linker/internal/auth"
	"linker/internal/database"
	"linker/internal/models"
	"linker/internal/storage"
)

func setupTestDB(t *testing.T) *database.Database {
	// Change to project root where migrations are located
	originalDir, _ := os.Getwd()
	projectRoot := filepath.Join(originalDir, "..")
	os.Chdir(projectRoot)
	
	testDBPath := "./test_file_sharing.db"
	os.Remove(testDBPath)
	
	db, err := database.Init(testDBPath)
	if err != nil {
		os.Chdir(originalDir) // Restore directory
		t.Fatalf("Failed to initialize test database: %v", err)
	}
	
	t.Cleanup(func() {
		db.Close()
		os.Remove(testDBPath)
		os.Chdir(originalDir) // Restore directory
	})
	
	return db
}

func setupTestFileSystem(t *testing.T) (*database.Database, *storage.S3Client) {
	db := setupTestDB(t)
	
	// For unit tests, we'll mock S3 operations or skip them
	var s3Client *storage.S3Client = nil
	
	return db, s3Client
}

func createTestUser(t *testing.T, db *database.Database, username, email string) *models.User {
	hashedPassword, err := auth.HashPassword("testpassword123")
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	
	user := &models.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}
	
	err = db.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	
	return user
}

func TestCreateFile(t *testing.T) {
	db, _ := setupTestFileSystem(t)
	user := createTestUser(t, db, "fileuser", "fileuser@example.com")
	
	file := &models.File{
		UserID:       user.ID,
		Filename:     "test_image.jpg",
		OriginalName: "original.jpg",
		MimeType:     "image/jpeg",
		FileSize:     1024,
		S3Key:        "2024/01/01/uuid-test.jpg",
		S3Bucket:     "test-bucket",
		Title:        "Test Image",
		Description:  "A test image file",
		Analytics:    true,
		IsPublic:     true,
	}
	
	err := db.CreateFile(file)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	
	if file.ID == "" {
		t.Fatal("File ID should be set after creation")
	}
	
	// Test retrieving the file
	retrievedFile, err := db.GetFileByID(file.ID, user.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve file: %v", err)
	}
	
	if retrievedFile.Filename != file.Filename {
		t.Errorf("Expected filename %s, got %s", file.Filename, retrievedFile.Filename)
	}
	
	if retrievedFile.UserID != user.ID {
		t.Errorf("Expected user ID %s, got %s", user.ID, retrievedFile.UserID)
	}
}

func TestCreateFileShortCode(t *testing.T) {
	db, _ := setupTestFileSystem(t)
	user := createTestUser(t, db, "shortcodeuser", "shortcode@example.com")
	
	// Create a file
	file := &models.File{
		UserID:       user.ID,
		Filename:     "test.pdf",
		OriginalName: "document.pdf",
		MimeType:     "application/pdf",
		FileSize:     2048,
		S3Key:        "2024/01/01/uuid-test.pdf",
		S3Bucket:     "test-bucket",
		IsPublic:     true,
	}
	
	err := db.CreateFile(file)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	
	// Create short codes for the file
	shortCodes := []string{"testpdf", "mydoc"}
	for i, shortCode := range shortCodes {
		isPrimary := i == 0
		err := db.CreateFileShortCode(file.ID, shortCode, isPrimary)
		if err != nil {
			t.Fatalf("Failed to create short code %s: %v", shortCode, err)
		}
	}
	
	// Test retrieving file by short code
	retrievedFile, err := db.GetFileByShortCode("testpdf")
	if err != nil {
		t.Fatalf("Failed to retrieve file by short code: %v", err)
	}
	
	if retrievedFile.ID != file.ID {
		t.Errorf("Expected file ID %s, got %s", file.ID, retrievedFile.ID)
	}
	
	// Test retrieving short codes for the file
	shortCodeRecords, err := db.GetShortCodesByFileID(file.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve short codes: %v", err)
	}
	
	if len(shortCodeRecords) != 2 {
		t.Errorf("Expected 2 short codes, got %d", len(shortCodeRecords))
	}
	
	// Verify primary short code
	primaryFound := false
	for _, sc := range shortCodeRecords {
		if sc.IsPrimary && sc.ShortCode == "testpdf" {
			primaryFound = true
		}
	}
	
	if !primaryFound {
		t.Error("Primary short code not found or incorrect")
	}
}

func TestFileDownloadTracking(t *testing.T) {
	db, _ := setupTestFileSystem(t)
	user := createTestUser(t, db, "analyticsuser", "analytics@example.com")
	
	// Create a file with analytics enabled
	file := &models.File{
		UserID:       user.ID,
		Filename:     "tracked.txt",
		OriginalName: "tracked_file.txt",
		MimeType:     "text/plain",
		FileSize:     512,
		S3Key:        "2024/01/01/uuid-tracked.txt",
		S3Bucket:     "test-bucket",
		Analytics:    true,
		IsPublic:     true,
	}
	
	err := db.CreateFile(file)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	
	// Simulate file downloads
	downloads := []models.FileDownload{
		{
			FileID:    file.ID,
			IPAddress: "192.168.1.1",
			UserAgent: "Mozilla/5.0 Test Browser",
			Referer:   "https://example.com",
		},
		{
			FileID:    file.ID,
			IPAddress: "192.168.1.2",
			UserAgent: "Mozilla/5.0 Another Browser",
			Referer:   "https://test.com",
		},
	}
	
	for _, download := range downloads {
		err := db.CreateFileDownload(&download)
		if err != nil {
			t.Fatalf("Failed to create download record: %v", err)
		}
	}
	
	// Increment download counter
	err = db.IncrementFileDownloads(file.ID)
	if err != nil {
		t.Fatalf("Failed to increment downloads: %v", err)
	}
	err = db.IncrementFileDownloads(file.ID)
	if err != nil {
		t.Fatalf("Failed to increment downloads: %v", err)
	}
	
	// Test analytics retrieval
	analytics, err := db.GetFileAnalytics(file.ID, user.ID)
	if err != nil {
		t.Fatalf("Failed to get file analytics: %v", err)
	}
	
	if len(analytics) != 2 {
		t.Errorf("Expected 2 download records, got %d", len(analytics))
	}
	
	// Test analytics summary
	summary, err := db.GetFileAnalyticsSummary(file.ID, user.ID)
	if err != nil {
		t.Fatalf("Failed to get analytics summary: %v", err)
	}
	
	if summary.TotalDownloads != 2 {
		t.Errorf("Expected 2 total downloads, got %d", summary.TotalDownloads)
	}
	
	if len(summary.TopReferrers) == 0 {
		t.Error("Expected referrer data, got none")
	}
}

func TestFilePasswordProtection(t *testing.T) {
	db, _ := setupTestFileSystem(t)
	user := createTestUser(t, db, "pwduser", "pwd@example.com")
	
	// Hash a password
	password := "secretfile123"
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	
	// Create a password-protected file
	file := &models.File{
		UserID:       user.ID,
		Filename:     "secret.txt",
		OriginalName: "secret_document.txt",
		MimeType:     "text/plain",
		FileSize:     256,
		S3Key:        "2024/01/01/uuid-secret.txt",
		S3Bucket:     "test-bucket",
		IsPublic:     false,
		Password:     &hashedPassword,
	}
	
	err = db.CreateFile(file)
	if err != nil {
		t.Fatalf("Failed to create password-protected file: %v", err)
	}
	
	// Create short code
	err = db.CreateFileShortCode(file.ID, "secretfile", true)
	if err != nil {
		t.Fatalf("Failed to create short code: %v", err)
	}
	
	// Test password verification
	retrievedFile, err := db.GetFileByShortCode("secretfile")
	if err != nil {
		t.Fatalf("Failed to retrieve file: %v", err)
	}
	
	if retrievedFile.Password == nil {
		t.Fatal("Expected password to be set")
	}
	
	// Verify correct password
	if !auth.CheckPassword(password, *retrievedFile.Password) {
		t.Error("Password verification failed for correct password")
	}
	
	// Verify incorrect password
	if auth.CheckPassword("wrongpassword", *retrievedFile.Password) {
		t.Error("Password verification should fail for incorrect password")
	}
}

func TestFileExpiration(t *testing.T) {
	db, _ := setupTestFileSystem(t)
	user := createTestUser(t, db, "expireuser", "expire@example.com")
	
	// Create a file that expires in the past
	pastTime := time.Now().Add(-24 * time.Hour)
	file := &models.File{
		UserID:       user.ID,
		Filename:     "expired.txt",
		OriginalName: "expired_document.txt",
		MimeType:     "text/plain",
		FileSize:     128,
		S3Key:        "2024/01/01/uuid-expired.txt",
		S3Bucket:     "test-bucket",
		IsPublic:     true,
		ExpiresAt:    &pastTime,
	}
	
	err := db.CreateFile(file)
	if err != nil {
		t.Fatalf("Failed to create expiring file: %v", err)
	}
	
	// Create short code
	err = db.CreateFileShortCode(file.ID, "expiredfile", true)
	if err != nil {
		t.Fatalf("Failed to create short code: %v", err)
	}
	
	// Test that we can retrieve the file (expiration logic is handled at the handler level)
	retrievedFile, err := db.GetFileByShortCode("expiredfile")
	if err != nil {
		t.Fatalf("Failed to retrieve expired file: %v", err)
	}
	
	// Verify expiration time
	if retrievedFile.ExpiresAt == nil {
		t.Fatal("Expected expiration time to be set")
	}
	
	if !retrievedFile.ExpiresAt.Before(time.Now()) {
		t.Error("File should be expired")
	}
}

func TestUserFileAnalytics(t *testing.T) {
	db, _ := setupTestFileSystem(t)
	user := createTestUser(t, db, "statsuser", "stats@example.com")
	
	// Create multiple files
	files := []*models.File{
		{
			UserID:       user.ID,
			Filename:     "doc1.pdf",
			OriginalName: "document1.pdf",
			MimeType:     "application/pdf",
			FileSize:     1024,
			S3Key:        "2024/01/01/uuid-doc1.pdf",
			S3Bucket:     "test-bucket",
			IsPublic:     true,
		},
		{
			UserID:       user.ID,
			Filename:     "img1.jpg",
			OriginalName: "image1.jpg",
			MimeType:     "image/jpeg",
			FileSize:     2048,
			S3Key:        "2024/01/01/uuid-img1.jpg",
			S3Bucket:     "test-bucket",
			IsPublic:     true,
		},
	}
	
	for _, file := range files {
		err := db.CreateFile(file)
		if err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
		
		// Simulate some downloads
		for i := 0; i < 3; i++ {
			err = db.IncrementFileDownloads(file.ID)
			if err != nil {
				t.Fatalf("Failed to increment downloads: %v", err)
			}
		}
	}
	
	// Test user file analytics
	userStats, err := db.GetUserFileAnalytics(user.ID)
	if err != nil {
		t.Fatalf("Failed to get user file analytics: %v", err)
	}
	
	if len(userStats) != 2 {
		t.Errorf("Expected 2 files in user stats, got %d", len(userStats))
	}
	
	for _, stat := range userStats {
		if stat.TotalDownloads != 3 {
			t.Errorf("Expected 3 downloads for file %s, got %d", stat.Filename, stat.TotalDownloads)
		}
	}
}

func TestFileUpdateAndDelete(t *testing.T) {
	db, _ := setupTestFileSystem(t)
	user := createTestUser(t, db, "updateuser", "update@example.com")
	
	// Create a file
	file := &models.File{
		UserID:       user.ID,
		Filename:     "update_me.txt",
		OriginalName: "update_document.txt",
		MimeType:     "text/plain",
		FileSize:     512,
		S3Key:        "2024/01/01/uuid-update.txt",
		S3Bucket:     "test-bucket",
		Title:        "Original Title",
		Description:  "Original Description",
		IsPublic:     true,
	}
	
	err := db.CreateFile(file)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	
	// Test file update
	updateReq := &models.UpdateFileRequest{
		Title:       "Updated Title",
		Description: "Updated Description",
		IsPublic:    false,
	}
	
	err = db.UpdateFile(file.ID, user.ID, updateReq)
	if err != nil {
		t.Fatalf("Failed to update file: %v", err)
	}
	
	// Verify update
	updatedFile, err := db.GetFileByID(file.ID, user.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve updated file: %v", err)
	}
	
	if updatedFile.Title != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got '%s'", updatedFile.Title)
	}
	
	if updatedFile.IsPublic {
		t.Error("Expected file to be private after update")
	}
	
	// Test file deletion
	err = db.DeleteFile(file.ID, user.ID)
	if err != nil {
		t.Fatalf("Failed to delete file: %v", err)
	}
	
	// Verify deletion
	_, err = db.GetFileByID(file.ID, user.ID)
	if err == nil {
		t.Error("Expected error when retrieving deleted file")
	}
}

// Helper function to create a multipart form for file uploads
func createMultipartForm(filename string, content []byte, fields map[string]string) (*bytes.Buffer, string) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	
	// Add file
	fileWriter, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, ""
	}
	fileWriter.Write(content)
	
	// Add form fields
	for key, value := range fields {
		writer.WriteField(key, value)
	}
	
	writer.Close()
	return &buf, writer.FormDataContentType()
}

func TestFileShortCodeUniqueness(t *testing.T) {
	db, _ := setupTestFileSystem(t)
	user := createTestUser(t, db, "uniqueuser", "unique@example.com")
	
	// Create first file
	file1 := &models.File{
		UserID:       user.ID,
		Filename:     "file1.txt",
		OriginalName: "file1.txt",
		MimeType:     "text/plain",
		FileSize:     100,
		S3Key:        "2024/01/01/uuid-file1.txt",
		S3Bucket:     "test-bucket",
		IsPublic:     true,
	}
	
	err := db.CreateFile(file1)
	if err != nil {
		t.Fatalf("Failed to create first file: %v", err)
	}
	
	// Create short code for first file
	err = db.CreateFileShortCode(file1.ID, "uniquecode", true)
	if err != nil {
		t.Fatalf("Failed to create short code for first file: %v", err)
	}
	
	// Create second file
	file2 := &models.File{
		UserID:       user.ID,
		Filename:     "file2.txt",
		OriginalName: "file2.txt",
		MimeType:     "text/plain",
		FileSize:     200,
		S3Key:        "2024/01/01/uuid-file2.txt",
		S3Bucket:     "test-bucket",
		IsPublic:     true,
	}
	
	err = db.CreateFile(file2)
	if err != nil {
		t.Fatalf("Failed to create second file: %v", err)
	}
	
	// Try to create same short code for second file - should fail
	err = db.CreateFileShortCode(file2.ID, "uniquecode", true)
	if err == nil {
		t.Error("Expected error when creating duplicate short code")
	}
}