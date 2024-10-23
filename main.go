// Package main is the main package of the program
package main

import (
	"bufio"
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"winder.website/sbfm/db"
	"winder.website/sbfm/prompt"
)

func main() {
	dbConnection, err := sql.Open("sqlite3", "./config.db")
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer dbConnection.Close()

	// Create tables if they don't exist
	err = db.CreateTables(dbConnection)
	if err != nil {
		log.Fatal("Error creating tables:", err)
	}

	// Create a scanner for reading input
	scanner := bufio.NewScanner(os.Stdin)

	// Main application loop
	prompt.HandleMenu(scanner, dbConnection)
}
