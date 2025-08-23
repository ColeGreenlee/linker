package tests

import (
	"testing"

	"linker/internal/auth"
	"linker/internal/models"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"
	hash, err := auth.HashPassword(password)
	
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	
	if hash == "" {
		t.Fatal("Hash should not be empty")
	}
	
	if hash == password {
		t.Fatal("Hash should not equal original password")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "testpassword123"
	wrongPassword := "wrongpassword"
	
	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	
	if !auth.CheckPassword(password, hash) {
		t.Fatal("Password check should pass for correct password")
	}
	
	if auth.CheckPassword(wrongPassword, hash) {
		t.Fatal("Password check should fail for incorrect password")
	}
}

func TestGenerateToken(t *testing.T) {
	user := &models.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}
	
	secret := "test-secret"
	token, err := auth.GenerateToken(user, secret)
	
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	if token == "" {
		t.Fatal("Token should not be empty")
	}
}

func TestValidateToken(t *testing.T) {
	user := &models.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}
	
	secret := "test-secret"
	token, err := auth.GenerateToken(user, secret)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	claims, err := auth.ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}
	
	if claims.UserID != user.ID {
		t.Fatalf("Expected user ID %d, got %d", user.ID, claims.UserID)
	}
	
	if claims.Username != user.Username {
		t.Fatalf("Expected username %s, got %s", user.Username, claims.Username)
	}
}

func TestValidateTokenInvalid(t *testing.T) {
	secret := "test-secret"
	invalidToken := "invalid.token.here"
	
	_, err := auth.ValidateToken(invalidToken, secret)
	if err == nil {
		t.Fatal("Should fail to validate invalid token")
	}
}

func TestValidateTokenWrongSecret(t *testing.T) {
	user := &models.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}
	
	secret := "test-secret"
	wrongSecret := "wrong-secret"
	
	token, err := auth.GenerateToken(user, secret)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	_, err = auth.ValidateToken(token, wrongSecret)
	if err == nil {
		t.Fatal("Should fail to validate token with wrong secret")
	}
}