// Package prompt is for printing the prompt
package prompt

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"

	"winder.website/sbfm/db"
)

// DisplayTransportList lists all available transports in the database
func DisplayTransportList(dbConnection *sql.DB) {
	if err := db.PrintTransports(dbConnection); err != nil {
		log.Println("Error displaying transports:", err)
	}
}

// DeleteTransportByID deletes a transport by its ID from the database
func DeleteTransportByID(dbConnection *sql.DB) {
	fmt.Print("Enter the ID of the transport you want to delete: ")
	var transportID int
	_, err := fmt.Scanf("%d\n", &transportID)
	if err != nil {
		log.Println("Invalid input:", err)
		return
	}

	err = db.DeleteTransport(dbConnection, transportID)
	if err != nil {
		log.Println("Error deleting transport:", err)
	} else {
		fmt.Println("Transport deleted successfully.")
	}
}

// AddTransportPrompt Function to handle transport input
func AddTransportPrompt(scanner *bufio.Scanner, dbConnection *sql.DB) {
	var transportType, transportPath string

	const defaultTransportType = "ws"
	const defaultTransportPath = ""

	// Helper function to scan input with default fallback
	readInput := func(prompt string, defaultValue string) string {
		fmt.Print(prompt)
		if scanner.Scan() {
			input := scanner.Text()
			if input == "" {
				return defaultValue
			}
			return input
		}
		return defaultValue
	}

	transportType = readInput(
		"Enter transport type (e.g., ws, http, etc.): ",
		defaultTransportType,
	)

	transportPath = readInput(
		"Enter transport path (e.g., /ws): ",
		defaultTransportPath,
	)

	// Insert into the transport table
	err := db.AddTransport(dbConnection, transportType, transportPath)
	if err != nil {
		fmt.Printf("Error inserting transport: %v\n", err)
	} else {
		fmt.Println("Transport configuration saved.")
	}
}
