// Package prompt is for printing the prompt
package prompt

import (
	"fmt"
	"strings"
)

// GetBoolInput prompts the user for a boolean input and returns the parsed boolean value.
func GetBoolInput(prompt string) (bool, error) {
	var input string
	fmt.Print(prompt)
	_, err := fmt.Scanln(&input)
	if err != nil {
		return false, fmt.Errorf("error reading input: %v", err)
	}

	// Convert input to lowercase for case-insensitive comparison
	input = strings.ToLower(input)

	switch input {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("invalid input: %s, please enter 'true' or 'false'", input)
	}
}

// DisplayMenu displays the main menu and returns the user's choice
func DisplayMenu() int {
	fmt.Println("\n1. Add user manually")
	fmt.Println("2. Add users from JSON file")
	fmt.Println("3. Print all users")
	fmt.Println("4. Delete user by ID")
	fmt.Println("5. Generate config.json")
	fmt.Println("6. Add log block to the database")
	fmt.Println("7. Add inbound")
	fmt.Println("8. Transport") // New option for Transport
	fmt.Println("9. Exit")
	fmt.Print("Choose an option: ")

	var choice int
	fmt.Scanln(&choice)
	return choice
}
