// Package prompt is for printing the prompt
package prompt

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"winder.website/sbfm/db"
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
	fmt.Println("\n1. User Management") // New user management option
	fmt.Println("2. Generate config.json")
	fmt.Println("3. Add log block to the database")
	fmt.Println("4. Add inbound")
	fmt.Println("5. TLS")
	fmt.Println("6. Reality")
	fmt.Println("7. Handshake")
	fmt.Println("8. Transport")
	fmt.Println("0. Exit")
	fmt.Print("Choose an option: ")

	var choice int
	fmt.Scanln(&choice)
	return choice
}

// DisplayUserManagementMenu displays the user management menu and returns the user's choice
func DisplayUserManagementMenu() int {
	fmt.Println("\nUser Management:")
	fmt.Println("1. Add user manually")
	fmt.Println("2. Add users from JSON file")
	fmt.Println("3. Print all users")
	fmt.Println("4. Delete user by ID")
	fmt.Println("5. Return to main menu")
	fmt.Print("Choose an option: ")

	var choice int
	fmt.Scanln(&choice)
	return choice
}

// HandleInboundInput Function to handle inbound input
func HandleInboundInput(scanner *bufio.Scanner, dbConnection *sql.DB) {
	var inboundType, tag, listen, sniffTimeout string
	var sniff, sniffOverrideDestination bool
	var listenPort, transportID, tlsID, realityID, handshakeID int

	// Set default values
	const defaultSniffTimeout = "300ms"
	const defaultInboundType = "vless"
	const defaultTag = "vless-ws"
	const defaultListen = "::"
	const defaultListenPort = 8080
	const defaultSniff = true
	const defaultSniffOverrideDestination = false

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

	// Read inputs with default fallback
	inboundType = readInput(
		"Enter inbounds type (e.g., vless, vmess, etc.) [default: vless]: ",
		defaultInboundType,
	)
	tag = readInput("Enter inbounds tag (e.g., vless-ws) [default: vless-ws]: ", defaultTag)
	listen = readInput("Enter inbounds listenIP (e.g., ::) [default: ::]: ", defaultListen)

	fmt.Print("Enter inbounds listenPort (e.g., 8080) [default: 8080]: ")
	_, err := fmt.Scanf("%d\n", &listenPort)
	if err != nil || listenPort == 0 {
		listenPort = defaultListenPort
	}

	sniffTimeout = readInput(
		"Enter inbounds sniffTimeout (e.g., 300ms) [default: 300ms]: ",
		defaultSniffTimeout,
	)

	sniff, err = GetBoolInput("Enter inbounds sniff (true/false) [default: true]: ")
	if err != nil {
		log.Println(err)
		sniff = defaultSniff
	}

	sniffOverrideDestination, err = GetBoolInput(
		"Enter inbounds sniffOverrideDestination (true/false) [default: false]: ",
	)
	if err != nil {
		log.Println(err)
		sniffOverrideDestination = defaultSniffOverrideDestination
	}

	// Print available transports
	if err := db.PrintTransports(dbConnection); err != nil {
		log.Println(err)
		return
	}

	// Prompt user for transport ID
	fmt.Print("Enter the transport ID you want to use: ")
	_, err = fmt.Scanf("%d\n", &transportID)
	if err != nil {
		log.Println("Invalid input for transport ID:", err)
		return
	}

	// Print available tlss
	if err := db.PrintTLS(dbConnection); err != nil {
		log.Println(err)
		return
	}

	// Prompt user for tls ID
	fmt.Print("Enter the tls ID you want to use: ")
	_, err = fmt.Scanf("%d\n", &tlsID)
	if err != nil {
		log.Println("Invalid input for tls ID:", err)
		return
	}

	// Print available realitys
	if err := db.PrintReality(dbConnection); err != nil {
		log.Println(err)
		return
	}

	// Prompt user for reality ID
	fmt.Print("Enter the reality ID you want to use: ")
	_, err = fmt.Scanf("%d\n", &realityID)
	if err != nil {
		log.Println("Invalid input for reality ID:", err)
		return
	}

	// Print available handshakes
	if err := db.PrintHandshake(dbConnection); err != nil {
		log.Println(err)
		return
	}

	// Prompt user for handshake ID
	fmt.Print("Enter the handshake ID you want to use: ")
	_, err = fmt.Scanf("%d\n", &handshakeID)
	if err != nil {
		log.Println("Invalid input for handshake ID:", err)
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
		tlsID,
		realityID,
		handshakeID,
		sniff,
		sniffOverrideDestination,
	)
}

// HandleTLSInput Function to handle TLS input
func HandleTLSInput(scanner *bufio.Scanner, dbConnection *sql.DB) {
	var enabled bool
	var serverName, minVersion, maxVersion, certificatePath, keyPath string

	// Set default values
	const defaultenabled = false
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
	if !enabled { // if false was entered, use the default value
		enabled = defaultenabled
	} else {
		enabled = enabledValue
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
		enabled,
		serverName,
		minVersion,
		maxVersion,
		certificatePath,
		keyPath,
	)
}

// HandleTransportInput Function to handle transport input
func HandleTransportInput(scanner *bufio.Scanner, dbConnection *sql.DB) {
	var transportType, transportPath string

	const defaultTransportType = "ws"
	const defaultTransportPath = ""

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

	transportType = readInput(
		"Enter transport type (e.g., ws, http, etc.): ",
		defaultTransportType,
	)

	transportPath = readInput(
		"Enter transport path (e.g., /ws): ",
		defaultTransportPath,
	)

	// Insert into the transport table
	err := db.InsertTransport(dbConnection, transportType, transportPath)
	if err != nil {
		fmt.Printf("Error inserting transport: %v\n", err)
	} else {
		fmt.Println("Transport configuration saved.")
	}
}

// HandleRealityInput to handle transport input
func HandleRealityInput(scanner *bufio.Scanner, dbConnection *sql.DB) {
	var enabled bool
	var privateKey, shortID string

	const defaultenabled = false
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
	if !enabled { // if false was entered, use the default value
		enabled = defaultenabled
	} else {
		enabled = enabledValue
	}

	privateKey = readInput(
		"Enter transport path (e.g., jdasflkjdsfj) [default= jadasjlfkjsflkj]: ",
		defaultprivetkey,
	)

	shortID = readInput(
		"Enter transport path (e.g., jdsfka) [ default= hakdsj]: ",
		defaultshortid,
	)

	db.AddReality(
		dbConnection,
		enabled,
		privateKey,
		shortID,
	)
}

// HandleHandshakeInput Function to handle transport input
func HandleHandshakeInput(scanner *bufio.Scanner, dbConnection *sql.DB) {
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
