package tests

import (
	"os"
	"testing"

	"linker/internal/config"
)

func TestLoadDefaultConfig(t *testing.T) {
	config := config.Load()
	
	if config.Port != "8080" {
		t.Fatalf("Expected port 8080, got %s", config.Port)
	}
	
	if config.DatabaseURL != "./linker.db" {
		t.Fatalf("Expected database URL ./linker.db, got %s", config.DatabaseURL)
	}
	
	if config.DefaultDomain != "localhost:8080" {
		t.Fatalf("Expected default domain localhost:8080, got %s", config.DefaultDomain)
	}
	
	if len(config.AllowedDomains) != 1 || config.AllowedDomains[0] != "localhost:8080" {
		t.Fatalf("Expected allowed domains [localhost:8080], got %v", config.AllowedDomains)
	}
	
	if config.Analytics != true {
		t.Fatal("Expected analytics to be true by default")
	}
	
	if config.Environment != "development" {
		t.Fatalf("Expected environment development, got %s", config.Environment)
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	os.Setenv("PORT", "3000")
	os.Setenv("DATABASE_URL", "/tmp/test.db")
	os.Setenv("DEFAULT_DOMAIN", "example.com")
	os.Setenv("ALLOWED_DOMAINS", "example.com,test.com,another.com")
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("ANALYTICS", "false")
	os.Setenv("ENVIRONMENT", "production")
	
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("DEFAULT_DOMAIN")
		os.Unsetenv("ALLOWED_DOMAINS")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("ANALYTICS")
		os.Unsetenv("ENVIRONMENT")
	}()
	
	config := config.Load()
	
	if config.Port != "3000" {
		t.Fatalf("Expected port 3000, got %s", config.Port)
	}
	
	if config.DatabaseURL != "/tmp/test.db" {
		t.Fatalf("Expected database URL /tmp/test.db, got %s", config.DatabaseURL)
	}
	
	if config.DefaultDomain != "example.com" {
		t.Fatalf("Expected default domain example.com, got %s", config.DefaultDomain)
	}
	
	expectedDomains := []string{"example.com", "test.com", "another.com"}
	if len(config.AllowedDomains) != 4 {
		t.Fatalf("Expected 4 allowed domains (default + extras), got %d", len(config.AllowedDomains))
	}
	
	if config.AllowedDomains[0] != "example.com" {
		t.Fatalf("Expected default domain example.com at index 0, got %s", config.AllowedDomains[0])
	}
	
	for _, expected := range expectedDomains {
		found := false
		for _, actual := range config.AllowedDomains {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Expected domain %s to be in allowed domains", expected)
		}
	}
	
	if config.JWTSecret != "test-secret" {
		t.Fatalf("Expected JWT secret test-secret, got %s", config.JWTSecret)
	}
	
	if config.Analytics != false {
		t.Fatal("Expected analytics to be false")
	}
	
	if config.Environment != "production" {
		t.Fatalf("Expected environment production, got %s", config.Environment)
	}
}