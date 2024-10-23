// Package db handles the database
package db

import (
	"database/sql"

	//go-sqlite3 is the sql driver for sqlite in go
	_ "github.com/mattn/go-sqlite3"
)

// DeleteInbound deletes an inbound by ID
func DeleteInbound(dbConnection *sql.DB, inboundID int) error {
	_, err := dbConnection.Exec("DELETE FROM inbounds WHERE id = ?", inboundID)
	return err
}

// DeleteTransport deletes a transport by ID
func DeleteTransport(dbConnection *sql.DB, transportID int) error {
	_, err := dbConnection.Exec("DELETE FROM transports WHERE id = ?", transportID)
	return err
}

// DeleteTLS deletes a TLS configuration by ID
func DeleteTLS(dbConnection *sql.DB, tlsID int) error {
	_, err := dbConnection.Exec("DELETE FROM tls WHERE id = ?", tlsID)
	return err
}

// DeleteReality deletes a Reality configuration by ID
func DeleteReality(dbConnection *sql.DB, realityID int) error {
	_, err := dbConnection.Exec("DELETE FROM reality WHERE id = ?", realityID)
	return err
}

// DeleteHandshake deletes a Handshake configuration by ID
func DeleteHandshake(dbConnection *sql.DB, handshakeID int) error {
	_, err := dbConnection.Exec("DELETE FROM handshake WHERE id = ?", handshakeID)
	return err
}
