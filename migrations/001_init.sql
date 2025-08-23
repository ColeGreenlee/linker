-- Users table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Links table  
CREATE TABLE IF NOT EXISTS links (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    short_code TEXT NOT NULL UNIQUE,
    original_url TEXT NOT NULL,
    title TEXT,
    description TEXT,
    clicks INTEGER DEFAULT 0,
    analytics BOOLEAN DEFAULT 1,
    expires_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Clicks table for analytics
CREATE TABLE IF NOT EXISTS clicks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    link_id INTEGER NOT NULL,
    ip_address TEXT,
    user_agent TEXT,
    referer TEXT,
    country TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (link_id) REFERENCES links (id) ON DELETE CASCADE
);

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_links_user_id ON links (user_id);
CREATE INDEX IF NOT EXISTS idx_links_short_code ON links (short_code);
CREATE INDEX IF NOT EXISTS idx_clicks_link_id ON clicks (link_id);
CREATE INDEX IF NOT EXISTS idx_clicks_created_at ON clicks (created_at);

-- Trigger to update updated_at timestamp
CREATE TRIGGER IF NOT EXISTS update_users_timestamp 
    AFTER UPDATE ON users
BEGIN
    UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_links_timestamp 
    AFTER UPDATE ON links
BEGIN
    UPDATE links SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;