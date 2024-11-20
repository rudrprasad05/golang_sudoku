package database

import (
	"database/sql"
	"fmt"
	"log"
)

type Config struct {
	Username string
	Password string
	Host     string
	Port     int
	DbName   string
}

func InitDB(config Config) (*sql.DB, error) {
	// Connect to MySQL server without a specific database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", config.Username, config.Password, config.Host, config.Port)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL server: %w", err)
	}
	defer db.Close()

	// Check if the database exists
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM information_schema.schemata WHERE schema_name = ?)"
	err = db.QueryRow(query, config.DbName).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check database existence: %w", err)
	}

	if exists {
		log.Printf("Database '%s' already exists. Skipping creation.", config.DbName)
	} else {
		// Create the database if it doesn't exist
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", config.DbName))
		if err != nil {
			return nil, fmt.Errorf("failed to create database '%s': %w", config.DbName, err)
		}
		log.Printf("Database '%s' created successfully.", config.DbName)
	}

	// Connect to the specific database
	dsnWithDB := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.Username, config.Password, config.Host, config.Port, config.DbName)
	dbWithSchema, err := sql.Open("mysql", dsnWithDB)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database '%s': %w", config.DbName, err)
	}

	// Test the connection
	if err := dbWithSchema.Ping(); err != nil {
		return nil, fmt.Errorf("database '%s' is unreachable: %w", config.DbName, err)
	}

	log.Printf("Connected to the database '%s' successfully.", config.DbName)
	return dbWithSchema, nil
}
