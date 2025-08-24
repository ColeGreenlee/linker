package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Database struct {
	*sql.DB
}

func Init(databaseURL string) (*Database, error) {
	db, err := sql.Open("sqlite", databaseURL+"?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &Database{db}
	
	if err := database.migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return database, nil
}

func (db *Database) migrate() error {
	// First, ensure the migrations table exists
	migrationFile := filepath.Join("migrations", "000_schema_migrations.sql")
	content, err := ioutil.ReadFile(migrationFile)
	if err != nil {
		return fmt.Errorf("failed to read migration tracking file: %w", err)
	}
	
	if _, err := db.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to create migration tracking table: %w", err)
	}

	migrations := []string{
		"001_init.sql",
		"002_domains.sql", 
		"003_api_tokens.sql",
		"004_uuid_conversion.sql",
		"005_file_sharing.sql",
		"006_fix_short_codes_constraints.sql",
	}

	for _, migration := range migrations {
		// Check if migration has already been applied
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", migration).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to check migration status for %s: %w", migration, err)
		}
		
		if count > 0 {
			log.Printf("Migration %s already applied, skipping", migration)
			continue
		}

		// Apply the migration
		migrationFile := filepath.Join("migrations", migration)
		content, err := ioutil.ReadFile(migrationFile)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", migration, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", migration, err)
		}

		// Record the migration as applied
		_, err = db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", migration)
		if err != nil {
			return fmt.Errorf("failed to record migration %s: %w", migration, err)
		}

		log.Printf("Migration %s completed successfully", migration)
	}

	return nil
}