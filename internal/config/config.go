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
	JWTSecret      string   `json:"jwt_secret"`
	Analytics      bool     `json:"analytics"`
	Environment    string   `json:"environment"`
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
	if config.JWTSecret == "" {
		config.JWTSecret = "your-secret-key-change-this"
	}
	if config.Environment == "" {
		config.Environment = "development"
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
	
	fmt.Println("Loaded configuration from environment variables")
	return &Config{
		Port:           getEnv("PORT", "8080"),
		DatabaseURL:    getEnv("DATABASE_URL", "./linker.db"),
		DefaultDomain:  defaultDomain,
		AllowedDomains: allowedDomains,
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key-change-this"),
		Analytics:      getEnvBool("ANALYTICS", true),
		Environment:    getEnv("ENVIRONMENT", "development"),
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