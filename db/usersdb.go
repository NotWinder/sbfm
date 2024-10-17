// Package db handles the database
package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"winder.website/sbfm/jsonhandler"
	//go-sqlite3 is the sql driver for sqlite in go
	_ "github.com/mattn/go-sqlite3"
)

// AddUserManually is responsible for adding a user to the database duh
func AddUserManually(db *sql.DB) {
	var name string
	fmt.Print("Enter user name: ")
	fmt.Scanln(&name)

	uuid := uuid.New().String()

	_, err := db.Exec("INSERT INTO users (name, uuid) VALUES (?, ?)", name, uuid)
	if err != nil {
		log.Printf("Error adding user: %v", err)
		return
	}

	fmt.Printf("User added successfully. UUID: %s\n", uuid)
}

// AddUsersFromJSON is responsible for adding multiple users from a json file
func AddUsersFromJSON(db *sql.DB) {
	fmt.Print("Enter JSON file name: ")
	var filename string
	fmt.Scanln(&filename)

	data, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("Error reading file: %v", err)
		return
	}

	var users []jsonhandler.User
	err = json.Unmarshal(data, &users)
	if err != nil {
		log.Printf("Error parsing JSON: %v", err)
		return
	}

	for _, user := range users {
		_, err := db.Exec("INSERT INTO users (name, uuid) VALUES (?, ?)", user.Name, user.UUID)
		if err != nil {
			log.Printf("Error adding user %s: %v", user.Name, err)
		} else {
			fmt.Printf("User %s added successfully\n", user.Name)
		}
	}
}

// PrintAllUsers is responsible for fetching and printing all the users in the users table
func PrintAllUsers(db *sql.DB) {
	rows, err := db.Query("SELECT id, name, uuid FROM users")
	if err != nil {
		log.Printf("Error querying users: %v", err)
		return
	}
	defer rows.Close()

	fmt.Println("\nAll Users:")
	fmt.Println("ID.Name.UUID")
	fmt.Println("------------")

	for rows.Next() {
		var id int
		var name, uuid string
		err := rows.Scan(&id, &name, &uuid)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		fmt.Printf("%d.%s.%s\n", id, name, uuid)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
	}
}

// DeleteUserByID is responsible for what ever the name says idiot
func DeleteUserByID(db *sql.DB) {
	var id int
	fmt.Print("Enter the ID of the user to delete: ")
	_, err := fmt.Scanln(&id)
	if err != nil {
		log.Printf("Error reading input: %v", err)
		return
	}

	result, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		return
	}

	if rowsAffected == 0 {
		fmt.Printf("No user found with ID %d\n", id)
	} else {
		fmt.Printf("User with ID %d deleted successfully\n", id)
	}
}
