// Package prompt is for printing the prompt
package prompt

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"winder.website/sbfm/db"
)

// DisplayInboundList lists all available inbounds in the database
func DisplayInboundList(dbConnection *sql.DB) {
	if err := db.PrintInbounds(dbConnection); err != nil {
		log.Println("Error displaying inbounds:", err)
	}
}

// DeleteInboundByID deletes an inbound by its ID from the database
func DeleteInboundByID(dbConnection *sql.DB) {
	fmt.Print("Enter the ID of the inbound you want to delete: ")
	var inboundID int
	_, err := fmt.Scanf("%d\n", &inboundID)
	if err != nil {
		log.Println("Invalid input:", err)
		return
	}

	err = db.DeleteInbound(dbConnection, inboundID)
	if err != nil {
		log.Println("Error deleting inbound:", err)
	} else {
		fmt.Println("Inbound deleted successfully.")
	}
}

// AddInboundPrompt Function to handle inbound input
func AddInboundPrompt(scanner *bufio.Scanner, dbConnection *sql.DB) {
	var inboundType, tag, listen, sniffTimeout string
	var sniff, sniffOverrideDestination bool
	var listenPort int
	var transportID, tlsID, realityID, handshakeID *int

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

	// Optional Transport configuration
	useTransport, err := GetBoolInput(
		"Do you want to configure Transport (true/false)?: ",
	)
	if err == nil && useTransport {

		// Print available transports
		if err := db.PrintTransports(dbConnection); err != nil {
			log.Println(err)
			return
		}
		// If Transport is configured, assign transport ID
		fmt.Print("Enter the transport ID you want to use: ")
		transportID = new(int)
		if scanner.Scan() {
			*transportID, err = strconv.Atoi(scanner.Text())
			if err != nil {
				log.Println("Invalid transport ID. Defaulting to nil.")
				transportID = nil
			}
		}
	}

	// Optional TLS configuration
	useTLS, err := GetBoolInput("Do you want to configure TLS (true/false)? [default: false]: ")
	if err == nil && useTLS {

		// Print available TLS
		if err := db.PrintTLS(dbConnection); err != nil {
			log.Println(err)
			return
		}
		// If TLS is configured, assign tlsID
		fmt.Print("Enter the tls ID you want to use: ")
		tlsID = new(int)
		if scanner.Scan() {
			*tlsID, err = strconv.Atoi(scanner.Text())
			if err != nil {
				log.Println("Invalid TLS ID. Defaulting to nil.")
				tlsID = nil
			}
		}
	}

	// Optional Reality configuration
	useReality, err := GetBoolInput(
		"Do you want to configure Reality (true/false)? [default: false]: ",
	)
	if err == nil && useReality {

		// Print available realities
		if err := db.PrintReality(dbConnection); err != nil {
			log.Println(err)
			return
		}

		realityID = new(int)
		fmt.Print("Enter the reality ID you want to use: ")
		if scanner.Scan() {
			*realityID, err = strconv.Atoi(scanner.Text())
			if err != nil {
				log.Println("Invalid reality ID. Defaulting to nil.")
				realityID = nil
			}
		}
	}

	// Optional Handshake configuration
	useHandshake, err := GetBoolInput(
		"Do you want to configure Handshake (true/false)? [default: false]: ",
	)
	if err == nil && useHandshake {

		// Print available handshakes
		if err := db.PrintHandshake(dbConnection); err != nil {
			log.Println(err)
			return
		}

		handshakeID = new(int)
		fmt.Print("Enter the handshake ID you want to use: ")
		if scanner.Scan() {
			*handshakeID, err = strconv.Atoi(scanner.Text())
			if err != nil {
				log.Println("Invalid handshake ID. Defaulting to nil.")
				handshakeID = nil
			}
		}
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
