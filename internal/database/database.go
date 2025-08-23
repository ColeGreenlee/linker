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
	migrations := []string{
		"001_init.sql",
		"002_domains.sql", 
		"003_api_tokens.sql",
	}

	for _, migration := range migrations {
		migrationFile := filepath.Join("migrations", migration)
		content, err := ioutil.ReadFile(migrationFile)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", migration, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", migration, err)
		}

		log.Printf("Migration %s completed successfully", migration)
	}

	return nil
}