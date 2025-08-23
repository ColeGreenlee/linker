package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port           string
	DatabaseURL    string
	DefaultDomain  string
	AllowedDomains []string
	JWTSecret      string
	Analytics      bool
	Environment    string
}

func Load() *Config {
	defaultDomain := getEnv("DEFAULT_DOMAIN", "localhost:8080")
	allowedDomains := []string{defaultDomain}
	
	if extraDomains := getEnv("ALLOWED_DOMAINS", ""); extraDomains != "" {
		for _, domain := range splitAndTrim(extraDomains, ",") {
			allowedDomains = append(allowedDomains, domain)
		}
	}
	
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