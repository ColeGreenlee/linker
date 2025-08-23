package tests

import (
	"io/ioutil"
	"os"
	"testing"

	"linker/internal/config"
)

func TestLoadJSONConfig(t *testing.T) {
	// Create a temporary config file
	configJSON := `{
		"port": "3000",
		"database_url": "/tmp/test.db",
		"default_domain": "example.com",
		"allowed_domains": ["example.com", "test.com"],
		"jwt_secret": "test-secret-json",
		"analytics": false,
		"environment": "test"
	}`
	
	tmpfile, err := ioutil.TempFile("", "config*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	
	if _, err := tmpfile.Write([]byte(configJSON)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()
	
	// Set environment variable to use our temp file
	os.Setenv("CONFIG_FILE", tmpfile.Name())
	defer os.Unsetenv("CONFIG_FILE")
	
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
	
	if len(config.AllowedDomains) != 2 {
		t.Fatalf("Expected 2 allowed domains, got %d", len(config.AllowedDomains))
	}
	
	if config.JWTSecret != "test-secret-json" {
		t.Fatalf("Expected JWT secret test-secret-json, got %s", config.JWTSecret)
	}
	
	if config.Analytics != false {
		t.Fatal("Expected analytics to be false")
	}
	
	if config.Environment != "test" {
		t.Fatalf("Expected environment test, got %s", config.Environment)
	}
}

func TestJSONConfigWithDefaults(t *testing.T) {
	// Create a minimal config file
	configJSON := `{
		"default_domain": "test.local"
	}`
	
	tmpfile, err := ioutil.TempFile("", "config*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	
	if _, err := tmpfile.Write([]byte(configJSON)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()
	
	os.Setenv("CONFIG_FILE", tmpfile.Name())
	defer os.Unsetenv("CONFIG_FILE")
	
	config := config.Load()
	
	// Check defaults are applied
	if config.Port != "8080" {
		t.Fatalf("Expected default port 8080, got %s", config.Port)
	}
	
	if config.DatabaseURL != "./linker.db" {
		t.Fatalf("Expected default database URL ./linker.db, got %s", config.DatabaseURL)
	}
	
	if config.DefaultDomain != "test.local" {
		t.Fatalf("Expected default domain test.local, got %s", config.DefaultDomain)
	}
	
	// Default domain should be added to allowed domains
	if len(config.AllowedDomains) != 1 || config.AllowedDomains[0] != "test.local" {
		t.Fatalf("Expected allowed domains [test.local], got %v", config.AllowedDomains)
	}
}

func TestFallbackToEnvWhenNoJSONFile(t *testing.T) {
	// Set config file to non-existent file
	os.Setenv("CONFIG_FILE", "nonexistent.json")
	os.Setenv("PORT", "9000")
	os.Setenv("JWT_SECRET", "env-secret")
	
	defer func() {
		os.Unsetenv("CONFIG_FILE")
		os.Unsetenv("PORT")
		os.Unsetenv("JWT_SECRET")
	}()
	
	config := config.Load()
	
	if config.Port != "9000" {
		t.Fatalf("Expected port from env 9000, got %s", config.Port)
	}
	
	if config.JWTSecret != "env-secret" {
		t.Fatalf("Expected JWT secret from env env-secret, got %s", config.JWTSecret)
	}
}