// Package db handles the database
package db

import (
	"database/sql"
	"fmt"
	"log"

	//go-sqlite3 is the sql driver for sqlite in go
	_ "github.com/mattn/go-sqlite3"
)

// AddInbound Function to add an inbound
func AddInbound(
	db *sql.DB,
	inboundType, tag, listen, sniffTimeout string,
	listenPort int,
	transportID, tlsID, realityID, handshakeID *int,
	sniff, sniffOverrideDestination bool,
) {
	// Insert the inbound and associate it with the transport ID
	_, err := db.Exec(
		`
	INSERT INTO inbounds (type, tag, listen, listen_port, sniff, sniff_override_destination, sniff_timeout, transport_id, tls_id, reality_id, handshake_id)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		inboundType,
		tag,
		listen,
		listenPort,
		sniff,
		sniffOverrideDestination,
		sniffTimeout,
		transportID,
		tlsID,
		realityID,
		handshakeID,
	)
	if err != nil {
		log.Printf("Error adding inbound: %v", err)
		return
	}

	fmt.Println("Inbound added successfully.")
}

// AddTLS Function to add a tls
func AddTLS(
	db *sql.DB,
	enabled bool,
	serverName, minVersion, maxVersion, certificatePath, keyPath string,
) {
	// Insert the inbound and associate it with the transport ID
	_, err := db.Exec(
		`
	INSERT INTO tls (enabled, server_name, min_version, max_version, certificate_path, key_path)
	VALUES (?, ?, ?, ?, ?, ?)`,
		enabled,
		serverName,
		minVersion,
		maxVersion,
		certificatePath,
		keyPath,
	)
	if err != nil {
		log.Printf("Error adding tls: %v", err)
		return
	}

	fmt.Println("tls added successfully.")
}

// AddReality Function to add a reality
func AddReality(
	db *sql.DB,
	enabled bool,
	privateKey, shortID string,
) {
	// Insert the inbound and associate it with the transport ID
	_, err := db.Exec(
		`
	INSERT INTO reality (enabled, private_key, short_id)
	VALUES (?, ?, ?)`,
		enabled,
		privateKey,
		shortID,
	)
	if err != nil {
		log.Printf("Error adding reality: %v", err)
		return
	}

	fmt.Println("reality added successfully.")
}

// AddHandshake Function to add a handshake
func AddHandshake(
	db *sql.DB,
	server string,
	serverPort int,
) {
	// Insert the inbound and associate it with the transport ID
	_, err := db.Exec(
		`
	INSERT INTO handshake (server, server_port)
	VALUES (?, ?)`,
		server,
		serverPort,
	)
	if err != nil {
		log.Printf("Error adding handshake: %v", err)
		return
	}

	fmt.Println("handshake added successfully.")
}

// AddTransport inserts a new transport entry into the database.
func AddTransport(db *sql.DB, transportType, transportPath string) error {
	// Insert transport details into the transports table
	_, err := db.Exec(
		`INSERT INTO transports (type, path) VALUES (?, ?)`,
		transportType, transportPath)

	if err != nil {
		return fmt.Errorf("error inserting into transports table: %v", err)
	}

	return nil
}
