// Package prompt is for printing the prompt
package prompt

import (
	"database/sql"
	"fmt"

	"winder.website/sbfm/db"
)

// DisplayUserManagementMenu displays the user management menu and returns the user's choice
func DisplayUserManagementMenu() int {
	fmt.Println("\nUser Management:")
	fmt.Println("1. Add user manually")
	fmt.Println("2. Add users from JSON file")
	fmt.Println("3. Print all users")
	fmt.Println("4. Delete user by ID")
	fmt.Println("5. Activate/Deactivate user by ID")
	fmt.Println("0. Return to main menu")
	fmt.Print("Choose an option: ")

	var choice int
	fmt.Scanln(&choice)
	return choice
}

// HandleUserManagementMenu handles user input for management options
func HandleUserManagementMenu(dbConnection *sql.DB) {
	for {
		choice := DisplayUserManagementMenu()
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
			db.ToggleUserActiveStatus(dbConnection)
		case 0:
			// Return to the main menu
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}
