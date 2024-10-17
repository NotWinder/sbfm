// Package main is the main package of the program
package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"winder.website/sbfm/db"
	"winder.website/sbfm/jsonhandler"
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
	for {
		choice := prompt.DisplayMenu() // Use the prompt package to display menu

		switch choice {
		case 1: // User Management
			for {
				userChoice := prompt.DisplayUserManagementMenu()

				switch userChoice {
				case 1:
					db.AddUserManually(dbConnection)
				case 2:
					db.AddUsersFromJSON(dbConnection)
				case 3:
					db.PrintAllUsers(dbConnection)
				case 4:
					db.DeleteUserByID(dbConnection)
				case 5:
					// Return to the main menu
					break
				default:
					fmt.Println("Invalid option. Please try again.")
				}

				// Break out of the loop if the user chose to return to the main menu
				if userChoice == 5 {
					break
				}
			}
		case 2:
			jsonhandler.GenerateConfigFile(dbConnection)
		case 3:
			// Add log data
			err = db.AddLogData(dbConnection, false, "info", "/var/log/app.log", true)
			if err != nil {
				log.Fatalf("Failed to add log data: %v", err)
			}
		case 4:
			prompt.HandleInboundInput(scanner, dbConnection)
		case 5:
			prompt.HandleTLSInput(scanner, dbConnection)
		case 6:
			prompt.HandleRealityInput(scanner, dbConnection)
		case 7:
			prompt.HandleHandshakeInput(scanner, dbConnection)
		case 8:
			prompt.HandleTransportInput(scanner, dbConnection)
		case 0:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}
