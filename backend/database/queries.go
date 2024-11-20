package database

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var tableRegistry sync.Map

func CreateTableOnce(db *sql.DB, tableName string, createQuery string) error {
	// Check the registry to see if the table has already been created
	if _, exists := tableRegistry.Load(tableName); exists {
		log.Printf("Table '%s' is already initialized.", tableName)
		return nil
	}

	// Create the table if not already created
	_, err := db.Exec(createQuery)
	if err != nil {
		return fmt.Errorf("failed to create table '%s': %w", tableName, err)
	}

	// Mark the table as created in the registry
	tableRegistry.Store(tableName, true)
	log.Printf("Table '%s' created successfully.", tableName)
	return nil
}

func GetUserByEmail(db *sql.DB, email string) *User {
	query := QFindUserByEmail()

	var user User

	err := db.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No user found with that email.")
			return nil
		}
		log.Fatalf("Failed to execute query: %v", err)
	}

	return &user
}

func CreateNewUser(db *sql.DB, user *User) (*User, error){
	query := QCreateNewUser();
	
	result, err := db.Exec(query, user.Name, user.Email, user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}

	id := int(lastInsertID)
	user.ID = &id

	// Return the updated user
	return user, nil
}


func QCreateUserTable() string {
    return `
        CREATE TABLE IF NOT EXISTS users (
            id INT AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(100) NOT NULL,
            email VARCHAR(100) UNIQUE NOT NULL,
            password VARCHAR(100) NOT NULL
        );
    `
}

func QFindUserByEmail() string{
	return `
		SELECT *
		FROM users
		WHERE email = ?
	`
}

func QCreateNewUser() string{
	return `
		INSERT INTO users (name, email, password)
		VALUES (?,?,?)
	`
}