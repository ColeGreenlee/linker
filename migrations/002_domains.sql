-- Domains table for supporting multiple domains
CREATE TABLE IF NOT EXISTS domains (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    domain TEXT NOT NULL UNIQUE,
    is_default BOOLEAN DEFAULT 0,
    enabled BOOLEAN DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Add domain_id to links table
ALTER TABLE links ADD COLUMN domain_id INTEGER REFERENCES domains(id);

-- Create index for domain lookups
CREATE INDEX IF NOT EXISTS idx_domains_domain ON domains (domain);
CREATE INDEX IF NOT EXISTS idx_links_domain_id ON links (domain_id);

-- Trigger to update domains timestamp
CREATE TRIGGER IF NOT EXISTS update_domains_timestamp 
    AFTER UPDATE ON domains
BEGIN
    UPDATE domains SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;