// Package prompt is for printing the prompt
package prompt

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"

	"winder.website/sbfm/db"
)

// DisplayRealityList lists all available Reality configurations in the database
func DisplayRealityList(dbConnection *sql.DB) {
	if err := db.PrintReality(dbConnection); err != nil {
		log.Println("Error displaying Reality configurations:", err)
	}
}

// DeleteRealityByID deletes a Reality configuration by its ID from the database
func DeleteRealityByID(dbConnection *sql.DB) {
	fmt.Print("Enter the ID of the Reality configuration you want to delete: ")
	var realityID int
	_, err := fmt.Scanf("%d\n", &realityID)
	if err != nil {
		log.Println("Invalid input:", err)
		return
	}

	err = db.DeleteReality(dbConnection, realityID)
	if err != nil {
		log.Println("Error deleting Reality configuration:", err)
	} else {
		fmt.Println("Reality configuration deleted successfully.")
	}
}

// AddRealityPrompt to handle transport input
func AddRealityPrompt(scanner *bufio.Scanner, dbConnection *sql.DB) {
	var privateKey, shortID string

	const defaultprivetkey = "wKKkpH2-ccPqK3JUfrGiCcd62uZSsLBOScNRBd_BMUk"
	const defaultshortid = "3a630a0a"

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
		"Enter if you want reality to be enabled (true/false) [default: false]: ",
	)
	if err != nil {
		log.Println(err)
		return
	}

	privateKey = readInput(
		"Enter reality's privetkey (e.g., jdasflkjdsfj) [default= wKKkpH2-ccPqK3JUfrGiCcd62uZSsLBOScNRBd_BMUk]: ",
		defaultprivetkey,
	)

	shortID = readInput(
		"Enter reality's shortID (e.g., 3a630a0a) [ default= 3a630a0a]: ",
		defaultshortid,
	)

	db.AddReality(
		dbConnection,
		enabledValue,
		privateKey,
		shortID,
	)
}
