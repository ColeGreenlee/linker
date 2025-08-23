package models

import (
	"time"
)

type User struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Domain struct {
	ID        int       `json:"id" db:"id"`
	Domain    string    `json:"domain" db:"domain"`
	IsDefault bool      `json:"is_default" db:"is_default"`
	Enabled   bool      `json:"enabled" db:"enabled"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Link struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	DomainID    *int      `json:"domain_id,omitempty" db:"domain_id"`
	Domain      *Domain   `json:"domain,omitempty" db:"-"`
	ShortCode   string    `json:"short_code" db:"short_code"`
	OriginalURL string    `json:"original_url" db:"original_url"`
	Title       string    `json:"title,omitempty" db:"title"`
	Description string    `json:"description,omitempty" db:"description"`
	Clicks      int       `json:"clicks" db:"clicks"`
	Analytics   bool      `json:"analytics" db:"analytics"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Click struct {
	ID        int       `json:"id" db:"id"`
	LinkID    int       `json:"link_id" db:"link_id"`
	IPAddress string    `json:"ip_address" db:"ip_address"`
	UserAgent string    `json:"user_agent" db:"user_agent"`
	Referer   string    `json:"referer,omitempty" db:"referer"`
	Country   string    `json:"country,omitempty" db:"country"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CreateLinkRequest struct {
	OriginalURL string     `json:"original_url" binding:"required,url"`
	ShortCode   string     `json:"short_code,omitempty"`
	DomainID    *int       `json:"domain_id,omitempty"`
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	Analytics   bool       `json:"analytics"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

type UpdateLinkRequest struct {
	OriginalURL string     `json:"original_url,omitempty" binding:"omitempty,url"`
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	Analytics   bool       `json:"analytics"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type APIToken struct {
	ID         int        `json:"id" db:"id"`
	UserID     int        `json:"user_id" db:"user_id"`
	TokenHash  string     `json:"-" db:"token_hash"`
	Name       *string    `json:"name,omitempty" db:"name"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty" db:"last_used_at"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
}

type CreateTokenRequest struct {
	Name      *string    `json:"name,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type CreateTokenResponse struct {
	Token    string    `json:"token"`
	APIToken APIToken  `json:"api_token"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}