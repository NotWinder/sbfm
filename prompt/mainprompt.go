// Package prompt is for printing the prompt
package prompt

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"

	"winder.website/sbfm/db"
	"winder.website/sbfm/jsonhandler"
)

// DisplayMenu displays the main menu and returns the user's choice
func DisplayMenu() int {
	fmt.Println("\n1. User Management")
	fmt.Println("2. Generate config.json")
	fmt.Println("3. Add log block to the database")
	fmt.Println("4. Manage Inbounds, Transports, TLS, Reality, Handshake")
	fmt.Println("5. make users client files")
	fmt.Println("6. make users sub files")
	fmt.Println("0. Exit")
	fmt.Print("Choose an option: ")

	var choice int
	fmt.Scanln(&choice)
	return choice
}

// HandleMenu handles the main menu input
func HandleMenu(scanner *bufio.Scanner, dbConnection *sql.DB) {
	for {
		choice := DisplayMenu()
		switch choice {
		case 1:
			HandleUserManagementMenu(scanner, dbConnection)
		case 2:
			jsonhandler.GenerateConfigFile(dbConnection)
		case 3:
			// Add log data
			err := db.AddLogData(dbConnection, false, "info", "/var/log/app.log", true)
			if err != nil {
				log.Fatalf("Failed to add log data: %v", err)
			}
		case 4:
			HandleInboundManagementMenu(scanner, dbConnection)
		case 5:
			templateFilePath := "./template.json"
			if err := db.GenerateUserJSONFiles(dbConnection, templateFilePath); err != nil {
				log.Fatalf("failed to generate user JSON files: %v", err)
			}
		case 6:
			db.GenerateUserConfigFiles(dbConnection)
		case 0:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}
