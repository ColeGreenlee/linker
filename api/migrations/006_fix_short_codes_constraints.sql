-- Fix short_codes table to allow NULL link_id for file short codes
-- SQLite doesn't support modifying column constraints directly, so we need to recreate the table

-- Create new short_codes table with proper constraints
CREATE TABLE short_codes_fixed (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
    link_id TEXT,
    file_id TEXT,
    short_code TEXT NOT NULL UNIQUE,
    is_primary BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (link_id) REFERENCES links (id) ON DELETE CASCADE,
    FOREIGN KEY (file_id) REFERENCES files (id) ON DELETE CASCADE
);

-- Copy existing data
INSERT INTO short_codes_fixed (id, link_id, short_code, is_primary, created_at)
SELECT id, link_id, short_code, is_primary, created_at FROM short_codes;

-- Drop old table and rename new one
DROP TABLE short_codes;
ALTER TABLE short_codes_fixed RENAME TO short_codes;

-- Recreate indexes
CREATE INDEX IF NOT EXISTS idx_short_codes_link_id ON short_codes (link_id);
CREATE INDEX IF NOT EXISTS idx_short_codes_file_id ON short_codes (file_id);
CREATE INDEX IF NOT EXISTS idx_short_codes_code ON short_codes (short_code);

-- Update trigger to ensure short codes reference either a link OR a file, but not both
CREATE TRIGGER validate_short_code_reference
    BEFORE INSERT ON short_codes
    WHEN (NEW.link_id IS NULL AND NEW.file_id IS NULL) OR (NEW.link_id IS NOT NULL AND NEW.file_id IS NOT NULL)
BEGIN
    SELECT RAISE(FAIL, 'Short code must reference either a link or a file, but not both or neither');
END;

CREATE TRIGGER validate_short_code_reference_update
    BEFORE UPDATE ON short_codes
    WHEN (NEW.link_id IS NULL AND NEW.file_id IS NULL) OR (NEW.link_id IS NOT NULL AND NEW.file_id IS NOT NULL)
BEGIN
    SELECT RAISE(FAIL, 'Short code must reference either a link or a file, but not both or neither');
END;