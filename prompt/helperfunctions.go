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
