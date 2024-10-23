package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"winder.website/sbfm/jsonhandler"
)

// GenerateUserJSONFiles generates JSON files for each user based on a template
func GenerateUserJSONFiles(dbConnection *sql.DB, templateFilePath string) error {
	// Step 1: Query all users from the database
	rows, err := dbConnection.Query(`SELECT uuid, name FROM users`)
	if err != nil {
		return fmt.Errorf("error querying users table: %v", err)
	}
	defer rows.Close()

	var users []jsonhandler.User
	for rows.Next() {
		var user jsonhandler.User
		if err := rows.Scan(&user.UUID, &user.Name); err != nil {
			return fmt.Errorf("error scanning user row: %v", err)
		}
		users = append(users, user)
	}

	// Step 2: Read the JSON template file
	templateData, err := os.ReadFile(templateFilePath)
	if err != nil {
		return fmt.Errorf("error reading template file: %v", err)
	}

	// Step 3: Create the users directory
	usersDir := "./sing-box/users"
	if err := os.MkdirAll(usersDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating users directory: %v", err)
	}

	// Step 4: Generate JSON files for each user
	for _, user := range users {
		// Unmarshal the template JSON into a generic structure
		var jsonData interface{}
		if err := json.Unmarshal(templateData, &jsonData); err != nil {
			log.Printf("error unmarshalling template JSON: %v", err)
			continue
		}

		// Step 5: Replace UUIDs in the JSON structure
		replaceUUID(jsonData, user.UUID)

		// Step 6: Marshal the modified structure back to JSON
		modifiedJSON, err := json.MarshalIndent(jsonData, "", "  ")
		if err != nil {
			log.Printf("error marshalling modified JSON for user %s: %v", user.Name, err)
			continue
		}

		// Step 7: Write the modified JSON to a new file in the users directory
		fileName := filepath.Join(usersDir, fmt.Sprintf("%s.json", user.Name))
		if err := os.WriteFile(fileName, modifiedJSON, 0644); err != nil {
			log.Printf("error writing JSON file for user %s: %v", user.Name, err)
			continue // Continue processing other users even if one fails
		}
		log.Printf("Generated JSON file: %s", fileName)
	}

	return nil
}

// replaceUUID recursively replaces UUID placeholders in the JSON structure with the actual UUID
func replaceUUID(data interface{}, newUUID string) {
	switch v := data.(type) {
	case map[string]interface{}: // If the data is a map
		for key, value := range v {
			if key == "uuid" {
				v[key] = newUUID // Replace UUID
			} else {
				replaceUUID(value, newUUID) // Recurse into the value
			}
		}
	case []interface{}: // If the data is a slice
		for _, item := range v {
			replaceUUID(item, newUUID) // Recurse into each item
		}
	}
}
