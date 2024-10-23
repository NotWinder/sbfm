package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"winder.website/sbfm/jsonhandler"
)

// GenerateUserConfigFiles generates configuration files for each user based on their name and sub value
func GenerateUserConfigFiles(dbConnection *sql.DB) error {
	// Step 1: Query all users from the database, including the sub field
	rows, err := dbConnection.Query(`SELECT uuid, name, sub FROM users`)
	if err != nil {
		return fmt.Errorf("error querying users table: %v", err)
	}
	defer rows.Close()

	var users []jsonhandler.User
	for rows.Next() {
		var user jsonhandler.User
		if err := rows.Scan(&user.UUID, &user.Name, &user.SUB); err != nil {
			return fmt.Errorf("error scanning user row: %v", err)
		}
		users = append(users, user)
	}

	// Step 2: Create the users directory for configuration files
	configsDir := "./sing-box/sub" // Directory to store config files
	if err := os.MkdirAll(configsDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating configs directory: %v", err)
	}

	// Step 3: Generate config files for each user
	for _, user := range users {
		// Define the content for the configuration file
		content := fmt.Sprintf(`location /sub/%s {
    alias /etc/sing-box/users/%s.json;
}`, user.SUB, user.Name)

		// Step 4: Write the content to a new file in the configs directory
		fileName := filepath.Join(configsDir, fmt.Sprintf("%s", user.Name))
		if err := os.WriteFile(fileName, []byte(content), 0644); err != nil {
			log.Printf("error writing config file for user %s: %v", user.Name, err)
			continue // Continue processing other users even if one fails
		}
		log.Printf("Generated config file: %s", fileName)
	}

	return nil
}
