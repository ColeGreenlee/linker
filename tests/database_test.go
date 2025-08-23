package tests

import (
	"database/sql"
	"os"
	"testing"

	"linker/internal/auth"
	"linker/internal/database"
	"linker/internal/models"
)

func setupTestDB(t *testing.T) *database.Database {
	testDBPath := "./test.db"
	os.Remove(testDBPath)
	
	db, err := database.Init(testDBPath)
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}
	
	t.Cleanup(func() {
		db.Close()
		os.Remove(testDBPath)
	})
	
	return db
}

func TestCreateAndGetUser(t *testing.T) {
	db := setupTestDB(t)
	
	hashedPassword, err := auth.HashPassword("testpassword")
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: hashedPassword,
	}
	
	err = db.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	if user.ID == 0 {
		t.Fatal("User ID should be set after creation")
	}
	
	retrievedUser, err := db.GetUserByUsername("testuser")
	if err != nil {
		t.Fatalf("Failed to get user by username: %v", err)
	}
	
	if retrievedUser.ID != user.ID {
		t.Fatalf("Expected user ID %d, got %d", user.ID, retrievedUser.ID)
	}
	
	if retrievedUser.Username != user.Username {
		t.Fatalf("Expected username %s, got %s", user.Username, retrievedUser.Username)
	}
	
	retrievedUser, err = db.GetUserByEmail("test@example.com")
	if err != nil {
		t.Fatalf("Failed to get user by email: %v", err)
	}
	
	if retrievedUser.ID != user.ID {
		t.Fatalf("Expected user ID %d, got %d", user.ID, retrievedUser.ID)
	}
}

func TestCreateAndGetLink(t *testing.T) {
	db := setupTestDB(t)
	
	hashedPassword, _ := auth.HashPassword("testpassword")
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: hashedPassword,
	}
	err := db.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	link := &models.Link{
		UserID:      user.ID,
		ShortCode:   "test123",
		OriginalURL: "https://example.com",
		Title:       "Test Link",
		Description: "A test link",
		Analytics:   true,
	}
	
	err = db.CreateLink(link)
	if err != nil {
		t.Fatalf("Failed to create link: %v", err)
	}
	
	if link.ID == 0 {
		t.Fatal("Link ID should be set after creation")
	}
	
	retrievedLink, err := db.GetLinkByShortCode("test123")
	if err != nil {
		t.Fatalf("Failed to get link by short code: %v", err)
	}
	
	if retrievedLink.ID != link.ID {
		t.Fatalf("Expected link ID %d, got %d", link.ID, retrievedLink.ID)
	}
	
	if retrievedLink.OriginalURL != link.OriginalURL {
		t.Fatalf("Expected URL %s, got %s", link.OriginalURL, retrievedLink.OriginalURL)
	}
}

func TestGetUserLinks(t *testing.T) {
	db := setupTestDB(t)
	
	hashedPassword, _ := auth.HashPassword("testpassword")
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: hashedPassword,
	}
	err := db.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	link1 := &models.Link{
		UserID:      user.ID,
		ShortCode:   "test1",
		OriginalURL: "https://example1.com",
		Analytics:   true,
	}
	
	link2 := &models.Link{
		UserID:      user.ID,
		ShortCode:   "test2",
		OriginalURL: "https://example2.com",
		Analytics:   false,
	}
	
	err = db.CreateLink(link1)
	if err != nil {
		t.Fatalf("Failed to create link1: %v", err)
	}
	
	err = db.CreateLink(link2)
	if err != nil {
		t.Fatalf("Failed to create link2: %v", err)
	}
	
	links, err := db.GetUserLinks(user.ID, 10, 0)
	if err != nil {
		t.Fatalf("Failed to get user links: %v", err)
	}
	
	if len(links) != 2 {
		t.Fatalf("Expected 2 links, got %d", len(links))
	}
}

func TestUpdateLink(t *testing.T) {
	db := setupTestDB(t)
	
	hashedPassword, _ := auth.HashPassword("testpassword")
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: hashedPassword,
	}
	err := db.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	link := &models.Link{
		UserID:      user.ID,
		ShortCode:   "test123",
		OriginalURL: "https://example.com",
		Title:       "Original Title",
		Analytics:   true,
	}
	
	err = db.CreateLink(link)
	if err != nil {
		t.Fatalf("Failed to create link: %v", err)
	}
	
	updates := &models.UpdateLinkRequest{
		OriginalURL: "https://updated.com",
		Title:       "Updated Title",
		Analytics:   false,
	}
	
	err = db.UpdateLink(link.ID, updates)
	if err != nil {
		t.Fatalf("Failed to update link: %v", err)
	}
	
	updatedLink, err := db.GetLinkByShortCode("test123")
	if err != nil {
		t.Fatalf("Failed to get updated link: %v", err)
	}
	
	if updatedLink.OriginalURL != updates.OriginalURL {
		t.Fatalf("Expected URL %s, got %s", updates.OriginalURL, updatedLink.OriginalURL)
	}
	
	if updatedLink.Title != updates.Title {
		t.Fatalf("Expected title %s, got %s", updates.Title, updatedLink.Title)
	}
	
	if updatedLink.Analytics != updates.Analytics {
		t.Fatalf("Expected analytics %t, got %t", updates.Analytics, updatedLink.Analytics)
	}
}

func TestDeleteLink(t *testing.T) {
	db := setupTestDB(t)
	
	hashedPassword, _ := auth.HashPassword("testpassword")
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: hashedPassword,
	}
	err := db.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	link := &models.Link{
		UserID:      user.ID,
		ShortCode:   "test123",
		OriginalURL: "https://example.com",
		Analytics:   true,
	}
	
	err = db.CreateLink(link)
	if err != nil {
		t.Fatalf("Failed to create link: %v", err)
	}
	
	err = db.DeleteLink(link.ID, user.ID)
	if err != nil {
		t.Fatalf("Failed to delete link: %v", err)
	}
	
	_, err = db.GetLinkByShortCode("test123")
	if err != sql.ErrNoRows {
		t.Fatal("Link should be deleted")
	}
}

func TestIncrementLinkClicks(t *testing.T) {
	db := setupTestDB(t)
	
	hashedPassword, _ := auth.HashPassword("testpassword")
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: hashedPassword,
	}
	err := db.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	link := &models.Link{
		UserID:      user.ID,
		ShortCode:   "test123",
		OriginalURL: "https://example.com",
		Analytics:   true,
	}
	
	err = db.CreateLink(link)
	if err != nil {
		t.Fatalf("Failed to create link: %v", err)
	}
	
	err = db.IncrementLinkClicks(link.ID)
	if err != nil {
		t.Fatalf("Failed to increment clicks: %v", err)
	}
	
	updatedLink, err := db.GetLinkByShortCode("test123")
	if err != nil {
		t.Fatalf("Failed to get updated link: %v", err)
	}
	
	if updatedLink.Clicks != 1 {
		t.Fatalf("Expected 1 click, got %d", updatedLink.Clicks)
	}
}

func TestCreateAndGetClick(t *testing.T) {
	db := setupTestDB(t)
	
	hashedPassword, _ := auth.HashPassword("testpassword")
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: hashedPassword,
	}
	err := db.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	link := &models.Link{
		UserID:      user.ID,
		ShortCode:   "test123",
		OriginalURL: "https://example.com",
		Analytics:   true,
	}
	
	err = db.CreateLink(link)
	if err != nil {
		t.Fatalf("Failed to create link: %v", err)
	}
	
	click := &models.Click{
		LinkID:    link.ID,
		IPAddress: "192.168.1.1",
		UserAgent: "Mozilla/5.0",
		Referer:   "https://google.com",
		Country:   "US",
	}
	
	err = db.CreateClick(click)
	if err != nil {
		t.Fatalf("Failed to create click: %v", err)
	}
	
	clicks, err := db.GetLinkAnalytics(link.ID, user.ID)
	if err != nil {
		t.Fatalf("Failed to get analytics: %v", err)
	}
	
	if len(clicks) != 1 {
		t.Fatalf("Expected 1 click, got %d", len(clicks))
	}
	
	if clicks[0].IPAddress != click.IPAddress {
		t.Fatalf("Expected IP %s, got %s", click.IPAddress, clicks[0].IPAddress)
	}
}