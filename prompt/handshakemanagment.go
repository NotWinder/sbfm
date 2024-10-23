// Package prompt is for printing the prompt
package prompt

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"

	"winder.website/sbfm/db"
)

// DisplayHandshakeList lists all available Handshake configurations in the database
func DisplayHandshakeList(dbConnection *sql.DB) {
	if err := db.PrintHandshake(dbConnection); err != nil {
		log.Println("Error displaying Handshake configurations:", err)
	}
}

// DeleteHandshakeByID deletes a Handshake configuration by its ID from the database
func DeleteHandshakeByID(dbConnection *sql.DB) {
	fmt.Print("Enter the ID of the Handshake configuration you want to delete: ")
	var handshakeID int
	_, err := fmt.Scanf("%d\n", &handshakeID)
	if err != nil {
		log.Println("Invalid input:", err)
		return
	}

	err = db.DeleteHandshake(dbConnection, handshakeID)
	if err != nil {
		log.Println("Error deleting Handshake configuration:", err)
	} else {
		fmt.Println("Handshake configuration deleted successfully.")
	}
}

// AddHandshakePrompt Function to handle transport input
func AddHandshakePrompt(scanner *bufio.Scanner, dbConnection *sql.DB) {
	var server string
	var serverPort int

	const defaultserver = "www.yahoo.com"
	const defaultserverPort = 443

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

	server = readInput(
		"Enter handshakes server address (e.g., www.yahoo.com) [default= www.yahoo.com]: ",
		defaultserver,
	)

	fmt.Print("Enter handshakes server Port (e.g., 443) [default: 443]: ")
	_, err := fmt.Scanf("%d\n", &serverPort)
	if err != nil || serverPort == 0 {
		serverPort = defaultserverPort
	}

	db.AddHandshake(
		dbConnection,
		server,
		serverPort,
	)
}
