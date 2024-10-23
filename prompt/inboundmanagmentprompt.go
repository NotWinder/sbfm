// Package prompt is for printing the prompt
package prompt

import (
	"bufio"
	"database/sql"
	"fmt"
)

// DisplayInboundManagementMenu displays the menu for managing inbounds, transports, tls, etc.
func DisplayInboundManagementMenu() int {
	fmt.Println("\nManagement Menu:")
	fmt.Println("1. Add inbounds")
	fmt.Println("2. List all inbounds")
	fmt.Println("3. Delete inbound by ID")
	fmt.Println("4. Add transports")
	fmt.Println("5. List all transports")
	fmt.Println("6. Delete transport by ID")
	fmt.Println("7. Add TLS configurations")
	fmt.Println("8. List all TLS configurations")
	fmt.Println("9. Delete TLS by ID")
	fmt.Println("10. Add Reality configurations")
	fmt.Println("11. List all Reality configurations")
	fmt.Println("12. Delete Reality by ID")
	fmt.Println("13. Add Handshake configurations")
	fmt.Println("14. List all Handshake configurations")
	fmt.Println("15. Delete Handshake by ID")
	fmt.Println("0. Return to main menu")
	fmt.Print("Choose an option: ")

	var choice int
	fmt.Scanln(&choice)
	return choice
}

// HandleInboundManagementMenu handles user input for management options
func HandleInboundManagementMenu(scanner *bufio.Scanner, dbConnection *sql.DB) {
	for {
		choice := DisplayInboundManagementMenu()
		switch choice {
		case 1:
			AddInboundPrompt(scanner, dbConnection)
		case 2:
			DisplayInboundList(dbConnection)
		case 3:
			DeleteInboundByID(dbConnection)
		case 4:
			AddTransportPrompt(scanner, dbConnection)
		case 5:
			DisplayTransportList(dbConnection)
		case 6:
			DeleteTransportByID(dbConnection)
		case 7:
			AddTLSPrompt(scanner, dbConnection)
		case 8:
			DisplayTLSList(dbConnection)
		case 9:
			DeleteTLSByID(dbConnection)
		case 10:
			AddRealityPrompt(scanner, dbConnection)
		case 11:
			DisplayRealityList(dbConnection)
		case 12:
			DeleteRealityByID(dbConnection)
		case 13:
			AddHandshakePrompt(scanner, dbConnection)
		case 14:
			DisplayHandshakeList(dbConnection)
		case 15:
			DeleteHandshakeByID(dbConnection)
		case 0:
			return // Return to main menu
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}
