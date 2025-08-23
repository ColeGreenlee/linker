package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	*sql.DB
}

func Init(databaseURL string) (*Database, error) {
	db, err := sql.Open("sqlite3", databaseURL+"?_foreign_keys=on")
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
	migrationFile := filepath.Join("migrations", "001_init.sql")
	content, err := ioutil.ReadFile(migrationFile)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	if _, err := db.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}