// Package main is the main package of the program
package main

import (
	"database/sql"
	"fmt"
	"log"

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

	// Main application loop
	for {
		choice := prompt.DisplayMenu() // Use the prompt package to display menu

		switch choice {
		case 1:
			db.AddUserManually(dbConnection)
		case 2:
			db.AddUsersFromJSON(dbConnection)
		case 3:
			db.PrintAllUsers(dbConnection)
		case 4:
			db.DeleteUserByID(dbConnection)
		case 5:
			jsonhandler.GenerateConfigFile(dbConnection)
		case 6:
			// Add log data
			err = db.AddLogData(dbConnection, false, "info", "/var/log/app.log", true)
			if err != nil {
				log.Fatalf("Failed to add log data: %v", err)
			}
		case 7:
			var inboundType, tag, listen, sniffTimeout string
			var sniff, sniffOverrideDestination bool
			var listenPort, transportID int

			// Set default values
			const defaultSniffTimeout = "300ms"
			const defaultInboundType = "vless"
			const defaultTag = "vless-ws"
			const defaultListen = "::"
			const defaultListenPort = 8080
			const defaultSniff = true
			const defaultSniffOverrideDestination = false

			fmt.Print("Enter inbounds type (e.g., vless, vmess, etc.) [default: vless]: ")
			fmt.Scan(&inboundType)
			if inboundType == "" {
				inboundType = defaultInboundType
			}

			fmt.Print("Enter inbounds tag (e.g., vless-ws) [default: vless-ws]: ")
			fmt.Scan(&tag)
			if tag == "" {
				tag = defaultTag
			}

			fmt.Print("Enter inbounds listenIP (e.g., ::) [default: ::]: ")
			fmt.Scan(&listen)
			if listen == "" {
				listen = defaultListen
			}

			fmt.Print("Enter inbounds listenPort (e.g., 8080) [default: 8080]: ")
			_, err := fmt.Scan(&listenPort)
			if err != nil || listenPort == 0 {
				listenPort = defaultListenPort
			}

			fmt.Print("Enter inbounds sniffTimeout (e.g., 300ms) [default: 300ms]: ")
			fmt.Scan(&sniffTimeout)
			if sniffTimeout == "" {
				sniffTimeout = defaultSniffTimeout
			}

			sniffValue, err := prompt.GetBoolInput(
				"Enter inbounds sniff (true/false) [default: true]: ",
			)
			if err != nil {
				log.Println(err)
				return
			}
			if !sniffValue { // if false was entered, use the default value
				sniff = defaultSniff
			} else {
				sniff = sniffValue
			}

			sniffODValue, err := prompt.GetBoolInput(
				"Enter inbounds sniffOverrideDestination (true/false) [default: false]: ",
			)
			if err != nil {
				log.Println(err)
				return
			}
			if !sniffODValue { // if false was entered, use the default value
				sniffOverrideDestination = defaultSniffOverrideDestination
			} else {
				sniffOverrideDestination = sniffODValue
			}

			// Print available transports
			if err := db.PrintTransports(dbConnection); err != nil {
				log.Println(err)
				return
			}

			// Prompt user for transport ID
			fmt.Print("Enter the transport ID you want to use: ")
			_, errr := fmt.Scan(&transportID)
			if errr != nil {
				log.Println("Invalid input for transport ID:", errr)
				return
			}

			db.AddInbound(
				dbConnection,
				inboundType,
				tag,
				listen,
				sniffTimeout,
				listenPort,
				transportID,
				sniff,
				sniffOverrideDestination,
			)
		case 8: // Transport configuration
			var transportType, transportPath string

			fmt.Print("Enter transport type (e.g., ws, http, etc.): ")
			fmt.Scan(&transportType)

			fmt.Print("Enter transport path (e.g., /ws): ")
			fmt.Scan(&transportPath)

			// Insert into the transport table
			err := db.InsertTransport(dbConnection, transportType, transportPath)
			if err != nil {
				fmt.Printf("Error inserting transport: %v\n", err)
			} else {
				fmt.Println("Transport configuration saved.")
			}
		case 9:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}
