package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port           string   `json:"port"`
	DatabaseURL    string   `json:"database_url"`
	DefaultDomain  string   `json:"default_domain"`
	AllowedDomains []string `json:"allowed_domains"`
	UnifiedPrefix  string   `json:"unified_prefix,omitempty"`
	LinkPrefix     string   `json:"link_prefix,omitempty"`
	FilePrefix     string   `json:"file_prefix,omitempty"`
	JWTSecret      string   `json:"jwt_secret"`
	Analytics      bool     `json:"analytics"`
	Environment    string   `json:"environment"`
	S3             S3Config `json:"s3"`
}

type S3Config struct {
	Enabled          bool   `json:"enabled"`
	Endpoint         string `json:"endpoint"`
	Region           string `json:"region"`
	AccessKeyID      string `json:"access_key_id"`
	SecretAccessKey  string `json:"secret_access_key"`
	BucketName       string `json:"bucket_name"`
	UseSSL           bool   `json:"use_ssl"`
	MaxFileSize      int64  `json:"max_file_size_mb"`
	AllowedMimeTypes []string `json:"allowed_mime_types"`
}

func Load() *Config {
	// Try to load from JSON file first
	if config := loadFromJSON(); config != nil {
		return config
	}
	
	// Fall back to environment variables
	return loadFromEnv()
}

func loadFromJSON() *Config {
	configFile := getEnv("CONFIG_FILE", "config.json")
	
	// Check if config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return nil
	}
	
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Warning: Failed to read config file %s: %v\n", configFile, err)
		return nil
	}
	
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Printf("Warning: Failed to parse config file %s: %v\n", configFile, err)
		return nil
	}
	
	// Apply defaults for missing values
	if config.Port == "" {
		config.Port = "8080"
	}
	if config.DatabaseURL == "" {
		config.DatabaseURL = "./linker.db"
	}
	if config.DefaultDomain == "" {
		config.DefaultDomain = "localhost:8080"
	}
	if len(config.AllowedDomains) == 0 {
		config.AllowedDomains = []string{config.DefaultDomain}
	}
	// Handle unified prefix logic
	if config.UnifiedPrefix != "" {
		if config.LinkPrefix == "" {
			config.LinkPrefix = config.UnifiedPrefix
		}
		if config.FilePrefix == "" {
			config.FilePrefix = config.UnifiedPrefix
		}
	} else {
		// Set defaults for individual prefixes
		if config.LinkPrefix == "" {
			config.LinkPrefix = "s"
		}
		if config.FilePrefix == "" {
			config.FilePrefix = "f"
		}
	}
	if config.JWTSecret == "" {
		config.JWTSecret = "your-secret-key-change-this"
	}
	if config.Environment == "" {
		config.Environment = "development"
	}
	
	// Set S3 defaults
	if config.S3.Region == "" {
		config.S3.Region = "us-east-1"
	}
	if config.S3.MaxFileSize == 0 {
		config.S3.MaxFileSize = 100 // 100MB default
	}
	if len(config.S3.AllowedMimeTypes) == 0 {
		config.S3.AllowedMimeTypes = []string{
			"image/jpeg", "image/png", "image/gif", "image/webp",
			"application/pdf", "text/plain", "text/csv",
			"application/zip", "application/json",
			"video/mp4", "video/webm",
			"audio/mpeg", "audio/wav",
		}
	}
	
	// Ensure default domain is in allowed domains
	found := false
	for _, domain := range config.AllowedDomains {
		if domain == config.DefaultDomain {
			found = true
			break
		}
	}
	if !found {
		config.AllowedDomains = append([]string{config.DefaultDomain}, config.AllowedDomains...)
	}
	
	fmt.Printf("Loaded configuration from %s\n", configFile)
	return &config
}

func loadFromEnv() *Config {
	defaultDomain := getEnv("DEFAULT_DOMAIN", "localhost:8080")
	allowedDomains := []string{defaultDomain}
	
	if extraDomains := getEnv("ALLOWED_DOMAINS", ""); extraDomains != "" {
		for _, domain := range splitAndTrim(extraDomains, ",") {
			allowedDomains = append(allowedDomains, domain)
		}
	}
	
	// Handle unified prefix from environment
	unifiedPrefix := getEnv("UNIFIED_PREFIX", "")
	linkPrefix := getEnv("LINK_PREFIX", "")
	filePrefix := getEnv("FILE_PREFIX", "")
	
	if unifiedPrefix != "" {
		if linkPrefix == "" {
			linkPrefix = unifiedPrefix
		}
		if filePrefix == "" {
			filePrefix = unifiedPrefix
		}
	} else {
		if linkPrefix == "" {
			linkPrefix = "s"
		}
		if filePrefix == "" {
			filePrefix = "f"
		}
	}
	
	fmt.Println("Loaded configuration from environment variables")
	return &Config{
		Port:           getEnv("PORT", "8080"),
		DatabaseURL:    getEnv("DATABASE_URL", "./linker.db"),
		DefaultDomain:  defaultDomain,
		AllowedDomains: allowedDomains,
		LinkPrefix:     linkPrefix,
		FilePrefix:     filePrefix,
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key-change-this"),
		Analytics:      getEnvBool("ANALYTICS", true),
		Environment:    getEnv("ENVIRONMENT", "development"),
		S3: S3Config{
			Enabled:         getEnvBool("S3_ENABLED", false),
			Endpoint:        getEnv("S3_ENDPOINT", ""),
			Region:          getEnv("S3_REGION", "us-east-1"),
			AccessKeyID:     getEnv("S3_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("S3_SECRET_ACCESS_KEY", ""),
			BucketName:      getEnv("S3_BUCKET_NAME", "linker-files"),
			UseSSL:          getEnvBool("S3_USE_SSL", true),
			MaxFileSize:     getEnvInt64("S3_MAX_FILE_SIZE_MB", 100),
			AllowedMimeTypes: []string{
				"image/jpeg", "image/png", "image/gif", "image/webp",
				"application/pdf", "text/plain", "text/csv",
				"application/zip", "application/json",
				"video/mp4", "video/webm",
				"audio/mpeg", "audio/wav",
			},
		},
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	if val := os.Getenv(key); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			return b
		}
	}
	return defaultVal
}

func getEnvInt64(key string, defaultVal int64) int64 {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			return i
		}
	}
	return defaultVal
}

func splitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}