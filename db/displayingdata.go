// Package db handles the database
package db

import (
	"database/sql"
	"fmt"

	//go-sqlite3 is the sql driver for sqlite in go
	_ "github.com/mattn/go-sqlite3"
)

// PrintInbounds prints all the data in the inbounds table
func PrintInbounds(dbConnection *sql.DB) error {
	rows, err := dbConnection.Query(
		`SELECT id, type, tag, listen, listen_port, sniff, sniff_override_destination, sniff_timeout, transport_id, tls_id, reality_id, handshake_id FROM inbounds`,
	)
	if err != nil {
		return fmt.Errorf("error querying inbounds table: %v", err)
	}
	defer rows.Close()

	fmt.Println("Available Inbounds:")
	fmt.Println(
		"ID,\tType,\tTag,\tListen,\tListenPort,\tSniff,\tSniffOverrideDestination,\tSniffTimeout,\tTransportID,\tTLSID,\tRealityID,\tHandshakeID",
	)
	for rows.Next() {
		var id int
		var inboundType, port, settings string
		if err := rows.Scan(&id, &inboundType, &port, &settings); err != nil {
			return fmt.Errorf("error scanning inbound row: %v", err)
		}
		fmt.Printf("%d\t%s\t%s\t%s\n", id, inboundType, port, settings)
	}
	return nil
}

// PrintTransports prints all the data in the trasport table
func PrintTransports(dbConnection *sql.DB) error {
	rows, err := dbConnection.Query(`SELECT id, type, path FROM transports`)
	if err != nil {
		return fmt.Errorf("error querying transports table: %v", err)
	}
	defer rows.Close()

	fmt.Println("Available Transports:")
	fmt.Println("ID\tType\tPath")
	for rows.Next() {
		var id int
		var transportType, path string
		if err := rows.Scan(&id, &transportType, &path); err != nil {
			return fmt.Errorf("error scanning transport row: %v", err)
		}
		fmt.Printf("%d\t%s\t%s\n", id, transportType, path)
	}
	return nil
}

// PrintTLS prints all the data in the tls table
func PrintTLS(dbConnection *sql.DB) error {
	rows, err := dbConnection.Query(
		`SELECT id, enabled, server_name, min_version, max_version, certificate_path, key_path FROM tls`,
	)
	if err != nil {
		return fmt.Errorf("error querying tls table: %v", err)
	}
	defer rows.Close()

	fmt.Println("Available TLS:")
	fmt.Println("ID\tEnabled\tServer_Name\tMinVersion\tMaxVersion\tCertificate_Path\tKeyPath")
	for rows.Next() {
		var id int
		var enabled bool
		var serverName, minVersion, maxVersion, certificatePath, keyPath string
		if err := rows.Scan(&id, &enabled, &serverName, &minVersion, &maxVersion, &certificatePath, &keyPath); err != nil {
			return fmt.Errorf("error scanning tls row: %v", err)
		}
		fmt.Printf(
			"%d\t%t\t%s\n%s\n%s\n%s\n%s\n",
			id,
			enabled,
			serverName,
			minVersion,
			maxVersion,
			certificatePath,
			keyPath,
		)
	}
	return nil
}

// PrintReality prints all the data in the reality table
func PrintReality(dbConnection *sql.DB) error {
	rows, err := dbConnection.Query(
		`SELECT id, enabled, private_key, short_id FROM reality`,
	)
	if err != nil {
		return fmt.Errorf("error querying reality table: %v", err)
	}
	defer rows.Close()

	fmt.Println("Available Reality:")
	fmt.Println("ID\tEnabled\tPrivetKey\tShortID")
	for rows.Next() {
		var id int
		var enabled bool
		var privateKey, shortID string
		if err := rows.Scan(&id, &enabled, &privateKey, &shortID); err != nil {
			return fmt.Errorf("error scanning reality row: %v", err)
		}
		fmt.Printf(
			"%d\t%t\t%s\n%s\n",
			id,
			enabled,
			privateKey,
			shortID,
		)
	}
	return nil
}

// PrintHandshake prints all the data in the handshake table
func PrintHandshake(dbConnection *sql.DB) error {
	rows, err := dbConnection.Query(
		`SELECT id, server, server_port FROM handshake`,
	)
	if err != nil {
		return fmt.Errorf("error querying handshake table: %v", err)
	}
	defer rows.Close()

	fmt.Println("Available Handshake:")
	fmt.Println("ID\tServer\tServerPort")
	for rows.Next() {
		var id int
		var server, serverPort string
		if err := rows.Scan(&id, &server, &serverPort); err != nil {
			return fmt.Errorf("error scanning handshake row: %v", err)
		}
		fmt.Printf(
			"%d\t%s\n%s\n",
			id,
			server,
			serverPort,
		)
	}
	return nil
}
