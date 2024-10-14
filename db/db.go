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

// CreateTables function is responsible for creating the tables in the database.
func CreateTables(db *sql.DB) error {
	// Create log table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS log (
			disabled BOOLEAN,
			level TEXT,
			output TEXT,
			timestamp BOOLEAN
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating log table: %v", err)
	}

	// Create inbounds table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS inbounds (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			type TEXT NOT NULL,
			tag TEXT UNIQUE NOT NULL,
			listen TEXT NOT NULL,
			listen_port INTEGER NOT NULL,
	        tcp_fast_open BOOLEAN,
	        tcp_multi_path BOOLEAN,
	        udp_fragment BOOLEAN,
	        udp_timeout TEXT,
	        detour TEXT,
			sniff BOOLEAN NOT NULL,
			sniff_override_destination BOOLEAN NOT NULL,
			sniff_timeout TEXT NOT NULL,
            domain_strategy TEXT,
            udp_disable_domain_unmapping BOOLEAN,
            transport_id INTEGER, -- New column for transport ID
            FOREIGN KEY (transport_id) REFERENCES transports(id) -- Foreign key to transports table
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating inbounds table: %v", err)
	}

	// Create users table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
		    name TEXT NOT NULL,
		    uuid TEXT NOT NULL UNIQUE
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating users table: %v", err)
	}

	// Create transports table with headers column
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS transports (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT NOT NULL,
		path TEXT NOT NULL
	)
`)
	if err != nil {
		return fmt.Errorf("error creating transports table: %v", err)
	}

	return nil
}

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

// AddInbound Function to add an inbound
func AddInbound(
	db *sql.DB,
	inboundType, tag, listen, sniffTimeout string,
	listenPort, transportID int,
	sniff, sniffOverrideDestination bool,
) {
	// Insert the inbound and associate it with the transport ID
	_, err := db.Exec(
		`
	INSERT INTO inbounds (type, tag, listen, listen_port, sniff, sniff_override_destination, sniff_timeout, transport_id)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		inboundType,
		tag,
		listen,
		listenPort,
		sniff,
		sniffOverrideDestination,
		sniffTimeout,
		transportID, // Use the transport ID
	)
	if err != nil {
		log.Printf("Error adding inbound: %v", err)
		return
	}

	fmt.Println("Inbound and transport added successfully.")
}

// AddLogData inserts a new log entry into the log table.
func AddLogData(db *sql.DB, disabled bool, level, output string, timestamp bool) error {
	query := `
        INSERT INTO log (disabled, level, output, timestamp) 
        VALUES (?, ?, ?, ?)
    `
	_, err := db.Exec(query, disabled, level, output, timestamp)
	if err != nil {
		return fmt.Errorf("error inserting log data: %v", err)
	}
	fmt.Println("Log data added successfully.")
	return nil
}

// InsertTransport inserts a new transport entry into the database.
func InsertTransport(db *sql.DB, transportType, transportPath string) error {
	// Insert transport details into the transports table
	_, err := db.Exec(
		`INSERT INTO transports (type, path) VALUES (?, ?)`,
		transportType, transportPath)

	if err != nil {
		return fmt.Errorf("error inserting into transports table: %v", err)
	}

	return nil
}

// PrintTransports prints all the data in the trasport table
func PrintTransports(dbConnection *sql.DB) error {
	rows, err := dbConnection.Query(`SELECT id, type, path FROM transports`)
	if err != nil {
		return fmt.Errorf("error querying transports table: %v", err)
	}
	defer rows.Close()

	fmt.Println("Available Transports:")
	fmt.Println("ID\tType\tPath")
	for rows.Next() {
		var id int
		var transportType, path string
		if err := rows.Scan(&id, &transportType, &path); err != nil {
			return fmt.Errorf("error scanning transport row: %v", err)
		}
		fmt.Printf("%d\t%s\t%s\n", id, transportType, path)
	}
	return nil
}
