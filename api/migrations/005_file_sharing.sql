-- Add file sharing tables
-- This migration adds support for file uploads and sharing

CREATE TABLE IF NOT EXISTS files (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
    user_id TEXT NOT NULL,
    domain_id TEXT,
    filename TEXT NOT NULL,
    original_name TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    file_size INTEGER NOT NULL,
    s3_key TEXT NOT NULL,
    s3_bucket TEXT NOT NULL,
    title TEXT,
    description TEXT,
    downloads INTEGER DEFAULT 0,
    analytics BOOLEAN DEFAULT 1,
    is_public BOOLEAN DEFAULT 1,
    password TEXT,
    expires_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (domain_id) REFERENCES domains (id)
);

CREATE TABLE IF NOT EXISTS file_downloads (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
    file_id TEXT NOT NULL,
    ip_address TEXT,
    user_agent TEXT,
    referer TEXT,
    country TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (file_id) REFERENCES files (id) ON DELETE CASCADE
);

-- Update short_codes table to support both links and files
-- This will be safe even if the column already exists due to IF NOT EXISTS in 006
-- For now, just add the column, constraints will be handled in 006
ALTER TABLE short_codes ADD COLUMN file_id TEXT;

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_files_user_id ON files (user_id);
CREATE INDEX IF NOT EXISTS idx_files_s3_key ON files (s3_key);
CREATE INDEX IF NOT EXISTS idx_files_domain_id ON files (domain_id);
CREATE INDEX IF NOT EXISTS idx_files_created_at ON files (created_at);
CREATE INDEX IF NOT EXISTS idx_files_expires_at ON files (expires_at);
CREATE INDEX IF NOT EXISTS idx_files_mime_type ON files (mime_type);

CREATE INDEX IF NOT EXISTS idx_file_downloads_file_id ON file_downloads (file_id);
CREATE INDEX IF NOT EXISTS idx_file_downloads_created_at ON file_downloads (created_at);

CREATE INDEX IF NOT EXISTS idx_short_codes_file_id ON short_codes (file_id);

-- Create triggers for updated_at timestamps
CREATE TRIGGER IF NOT EXISTS update_files_timestamp 
    AFTER UPDATE ON files
BEGIN
    UPDATE files SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;