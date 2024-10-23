// Package prompt is for printing the prompt
package prompt

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"

	"winder.website/sbfm/db"
)

// DisplayTLSList lists all available TLS configurations in the database
func DisplayTLSList(dbConnection *sql.DB) {
	if err := db.PrintTLS(dbConnection); err != nil {
		log.Println("Error displaying TLS configurations:", err)
	}
}

// DeleteTLSByID deletes a TLS configuration by its ID from the database
func DeleteTLSByID(dbConnection *sql.DB) {
	fmt.Print("Enter the ID of the TLS configuration you want to delete: ")
	var tlsID int
	_, err := fmt.Scanf("%d\n", &tlsID)
	if err != nil {
		log.Println("Invalid input:", err)
		return
	}

	err = db.DeleteTLS(dbConnection, tlsID)
	if err != nil {
		log.Println("Error deleting TLS configuration:", err)
	} else {
		fmt.Println("TLS configuration deleted successfully.")
	}
}

// AddTLSPrompt Function to handle TLS input
func AddTLSPrompt(scanner *bufio.Scanner, dbConnection *sql.DB) {
	var serverName, minVersion, maxVersion, certificatePath, keyPath string

	// Set default values
	const defaultserverName = "www.yahoo.com"
	const defaultminVersion = ""
	const defaultmaxVersion = ""
	const defaultcertificatePath = ""
	const defaultkeyPath = ""

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

	enabledValue, err := GetBoolInput(
		"Enter if you want tls to be enabled (true/false) [default: false]: ",
	)
	if err != nil {
		log.Println(err)
		return
	}

	serverName = readInput(
		"Enter server-name type (e.g., www.yahoo.com, etc.) [default: www.yahoo.com]: ",
		defaultserverName,
	)

	minVersion = readInput(
		"Enter tls minVersion (e.g., 1.2) [default: ]: ",
		defaultminVersion,
	)

	maxVersion = readInput(
		"Enter tls maxVersion (e.g., 1.3) [default: ]: ",
		defaultmaxVersion,
	)

	certificatePath = readInput(
		"Enter tls certificatePath (e.g., /path/to/cert) [default: ]: ",
		defaultcertificatePath,
	)

	keyPath = readInput(
		"Enter tls keyPath (e.g., /path/to/key) [default: ]: ",
		defaultkeyPath,
	)

	db.AddTLS(
		dbConnection,
		enabledValue,
		serverName,
		minVersion,
		maxVersion,
		certificatePath,
		keyPath,
	)
}
