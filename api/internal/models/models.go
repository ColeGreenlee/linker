package models

import (
	"time"
)

type User struct {
	ID        string    `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Domain struct {
	ID        string    `json:"id" db:"id"`
	Domain    string    `json:"domain" db:"domain"`
	IsDefault bool      `json:"is_default" db:"is_default"`
	Enabled   bool      `json:"enabled" db:"enabled"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Link struct {
	ID          string      `json:"id" db:"id"`
	UserID      string      `json:"user_id" db:"user_id"`
	DomainID    *string     `json:"domain_id,omitempty" db:"domain_id"`
	Domain      *Domain     `json:"domain,omitempty" db:"-"`
	ShortCodes  []ShortCode `json:"short_codes,omitempty" db:"-"`
	OriginalURL string      `json:"original_url" db:"original_url"`
	Title       string      `json:"title,omitempty" db:"title"`
	Description string      `json:"description,omitempty" db:"description"`
	Clicks      int         `json:"clicks" db:"clicks"`
	Analytics   bool        `json:"analytics" db:"analytics"`
	ExpiresAt   *time.Time  `json:"expires_at,omitempty" db:"expires_at"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
}

type ShortCode struct {
	ID        string    `json:"id" db:"id"`
	LinkID    string    `json:"link_id" db:"link_id"`
	ShortCode string    `json:"short_code" db:"short_code"`
	IsPrimary bool      `json:"is_primary" db:"is_primary"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Click struct {
	ID        string    `json:"id" db:"id"`
	LinkID    string    `json:"link_id" db:"link_id"`
	IPAddress string    `json:"ip_address" db:"ip_address"`
	UserAgent string    `json:"user_agent" db:"user_agent"`
	Referer   string    `json:"referer,omitempty" db:"referer"`
	Country   string    `json:"country,omitempty" db:"country"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CreateLinkRequest struct {
	OriginalURL string     `json:"original_url" binding:"required,url"`
	ShortCodes  []string   `json:"short_codes,omitempty"`
	DomainID    *string    `json:"domain_id,omitempty"`
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
	ID         string     `json:"id" db:"id"`
	UserID     string     `json:"user_id" db:"user_id"`
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

type UserAnalytics struct {
	UserID          string                 `json:"user_id"`
	TotalLinks      int                    `json:"total_links"`
	TotalClicks     int                    `json:"total_clicks"`
	ClicksToday     int                    `json:"clicks_today"`
	ClicksThisWeek  int                    `json:"clicks_this_week"`
	ClicksThisMonth int                    `json:"clicks_this_month"`
	TopLinks        []LinkAnalyticsSummary `json:"top_links"`
	RecentClicks    []Click                `json:"recent_clicks"`
	ClicksByDate    []ClicksByDate         `json:"clicks_by_date"`
	TopReferrers    []ReferrerStats        `json:"top_referrers"`
	TopCountries    []CountryStats         `json:"top_countries"`
	TopUserAgents   []UserAgentStats       `json:"top_user_agents"`
}

type LinkAnalyticsSummary struct {
	LinkID      string `json:"link_id"`
	OriginalURL string `json:"original_url"`
	Title       string `json:"title"`
	ShortCode   string `json:"short_code"`
	TotalClicks int    `json:"total_clicks"`
}

type ClicksByDate struct {
	Date   string `json:"date"`
	Clicks int    `json:"clicks"`
}

type ReferrerStats struct {
	Referer string `json:"referer"`
	Clicks  int    `json:"clicks"`
}

type CountryStats struct {
	Country string `json:"country"`
	Clicks  int    `json:"clicks"`
}

type UserAgentStats struct {
	UserAgent string `json:"user_agent"`
	Clicks    int    `json:"clicks"`
}

// File sharing models
type File struct {
	ID           string     `json:"id" db:"id"`
	UserID       string     `json:"user_id" db:"user_id"`
	DomainID     *string    `json:"domain_id,omitempty" db:"domain_id"`
	Domain       *Domain    `json:"domain,omitempty" db:"-"`
	ShortCodes   []ShortCode `json:"short_codes,omitempty" db:"-"`
	Filename     string     `json:"filename" db:"filename"`
	OriginalName string     `json:"original_name" db:"original_name"`
	MimeType     string     `json:"mime_type" db:"mime_type"`
	FileSize     int64      `json:"file_size" db:"file_size"`
	S3Key        string     `json:"-" db:"s3_key"`
	S3Bucket     string     `json:"-" db:"s3_bucket"`
	Title        string     `json:"title,omitempty" db:"title"`
	Description  string     `json:"description,omitempty" db:"description"`
	Downloads    int        `json:"downloads" db:"downloads"`
	Analytics    bool       `json:"analytics" db:"analytics"`
	IsPublic     bool       `json:"is_public" db:"is_public"`
	Password     *string    `json:"-" db:"password"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

type FileDownload struct {
	ID        string    `json:"id" db:"id"`
	FileID    string    `json:"file_id" db:"file_id"`
	IPAddress string    `json:"ip_address" db:"ip_address"`
	UserAgent string    `json:"user_agent" db:"user_agent"`
	Referer   string    `json:"referer,omitempty" db:"referer"`
	Country   string    `json:"country,omitempty" db:"country"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CreateFileRequest struct {
	ShortCodes  []string   `json:"short_codes,omitempty"`
	DomainID    *string    `json:"domain_id,omitempty"`
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	Analytics   bool       `json:"analytics"`
	IsPublic    bool       `json:"is_public"`
	Password    *string    `json:"password,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

type UpdateFileRequest struct {
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	Analytics   bool       `json:"analytics"`
	IsPublic    bool       `json:"is_public"`
	Password    *string    `json:"password,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

type FileAnalyticsSummary struct {
	FileID             string         `json:"file_id"`
	TotalDownloads     int            `json:"total_downloads"`
	DownloadsToday     int            `json:"downloads_today"`
	DownloadsThisWeek  int            `json:"downloads_this_week"`
	DownloadsThisMonth int            `json:"downloads_this_month"`
	UniqueVisitors     int            `json:"unique_visitors"`
	TopReferrers       []ReferrerStat `json:"top_referrers"`
}

type ReferrerStat struct {
	Referer string `json:"referer"`
	Count   int    `json:"count"`
}

type UserFileStats struct {
	FileID          string    `json:"file_id"`
	Filename        string    `json:"filename"`
	OriginalName    string    `json:"original_name"`
	MimeType        string    `json:"mime_type"`
	FileSize        int64     `json:"file_size"`
	TotalDownloads  int       `json:"total_downloads"`
	RecentDownloads int       `json:"recent_downloads"`
	CreatedAt       time.Time `json:"created_at"`
}

type UserFileAnalytics struct {
	UserID           string                 `json:"user_id"`
	TotalFiles       int                    `json:"total_files"`
	TotalDownloads   int                    `json:"total_downloads"`
	TotalFileSize    int64                  `json:"total_file_size_bytes"`
	DownloadsToday   int                    `json:"downloads_today"`
	DownloadsThisWeek int                   `json:"downloads_this_week"`
	DownloadsThisMonth int                 `json:"downloads_this_month"`
	TopFiles         []FileAnalyticsSummary `json:"top_files"`
	RecentDownloads  []FileDownload         `json:"recent_downloads"`
	DownloadsByDate  []ClicksByDate         `json:"downloads_by_date"`
	TopFileTypes     []FileTypeStats        `json:"top_file_types"`
}

type FileTypeStats struct {
	MimeType  string `json:"mime_type"`
	Extension string `json:"extension"`
	Count     int    `json:"count"`
	TotalSize int64  `json:"total_size_bytes"`
}