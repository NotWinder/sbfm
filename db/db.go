// Package db handles the database
package db

import (
	"database/sql"
	"fmt"

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
            transport_id INTEGER,
            tls_id INTEGER,
            reality_id INTEGER,
            handshake_id INTEGER,
            FOREIGN KEY (transport_id) REFERENCES transports(id)
            FOREIGN KEY (tls_id) REFERENCES tls(id)
            FOREIGN KEY (reality_id) REFERENCES reality(id)
            FOREIGN KEY (handshake_id) REFERENCES handshake(id)
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
		    uuid TEXT NOT NULL UNIQUE,
		    sub TEXT NOT NULL UNIQUE
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating users table: %v", err)
	}

	// Create tls table with headers column
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS tls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		enabled BOOLEAN NOT NULL,
		server_name TEXT NOT NULL,
		min_version TEXT,
		max_version TEXT,
		certificate_path TEXT,
		key_path TEXT
	)
`)
	if err != nil {
		return fmt.Errorf("error creating tls table: %v", err)
	}

	// Create reality table with headers column
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS reality (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		enabled BOOLEAN NOT NULL,
		private_key TEXT NOT NULL,
		short_id INTEGER
	)
`)
	if err != nil {
		return fmt.Errorf("error creating reality table: %v", err)
	}

	// Create handshake table with headers column
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS handshake (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		server TEXT NOT NULL,
		server_port TEXT
	)
`)
	if err != nil {
		return fmt.Errorf("error creating handshake table: %v", err)
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
