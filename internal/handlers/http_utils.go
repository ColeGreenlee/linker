package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"linker/internal/auth"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func getUserID(r *http.Request) (int, error) {
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		return 0, http.ErrNoCookie
	}
	return strconv.Atoi(userIDStr)
}

func getUsername(r *http.Request) string {
	return r.Header.Get("X-Username")
}

func getPathParam(path, prefix string) string {
	if !strings.HasPrefix(path, prefix) {
		return ""
	}
	return strings.TrimPrefix(path, prefix)
}

func ValidateToken(tokenString, secret string) (*auth.Claims, error) {
	return auth.ValidateToken(tokenString, secret)
}