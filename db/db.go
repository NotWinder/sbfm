// Package db handles the database
package db

import (
	"database/sql"
	"fmt"
	"log"

	//go-sqlite3 is the sql driver for sqlite in go
	_ "github.com/mattn/go-sqlite3"
)

// CreateTables function is responsible for creating the tables in the database.
func CreateTables(db *sql.DB) error {
	// Create log table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS log (
			disabled BOOLEAN,
			level TEXT,
			output TEXT,
			timestamp BOOLEAN
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating log table: %v", err)
	}

	// Create inbounds table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS inbounds (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			type TEXT NOT NULL,
			tag TEXT UNIQUE NOT NULL,
			listen TEXT NOT NULL,
			listen_port INTEGER NOT NULL,
	        tcp_fast_open BOOLEAN,
	        tcp_multi_path BOOLEAN,
	        udp_fragment BOOLEAN,
	        udp_timeout TEXT,
	        detour TEXT,
			sniff BOOLEAN NOT NULL,
			sniff_override_destination BOOLEAN NOT NULL,
			sniff_timeout TEXT NOT NULL,
            domain_strategy TEXT,
            udp_disable_domain_unmapping BOOLEAN,
            transport_id INTEGER,
            tls_id INTEGER,
            reality_id INTEGER,
            handshake_id INTEGER,
            FOREIGN KEY (transport_id) REFERENCES transports(id)
            FOREIGN KEY (tls_id) REFERENCES tls(id)
            FOREIGN KEY (reality_id) REFERENCES reality(id)
            FOREIGN KEY (handshake_id) REFERENCES handshake(id)
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating inbounds table: %v", err)
	}

	// Create users table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
		    name TEXT NOT NULL,
		    uuid TEXT NOT NULL UNIQUE
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating users table: %v", err)
	}

	// Create tls table with headers column
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS tls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		enabled BOOLEAN NOT NULL,
		server_name TEXT NOT NULL,
		min_version TEXT,
		max_version TEXT,
		certificate_path TEXT,
		key_path TEXT
	)
`)
	if err != nil {
		return fmt.Errorf("error creating tls table: %v", err)
	}

	// Create reality table with headers column
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS reality (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		enabled BOOLEAN NOT NULL,
		private_key TEXT NOT NULL,
		short_id INTEGER
	)
`)
	if err != nil {
		return fmt.Errorf("error creating reality table: %v", err)
	}

	// Create handshake table with headers column
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS handshake (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		server TEXT NOT NULL,
		server_port TEXT
	)
`)
	if err != nil {
		return fmt.Errorf("error creating handshake table: %v", err)
	}

	// Create transports table with headers column
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS transports (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT NOT NULL,
		path TEXT NOT NULL
	)
`)
	if err != nil {
		return fmt.Errorf("error creating transports table: %v", err)
	}

	return nil
}

// AddLogData inserts a new log entry into the log table.
func AddLogData(db *sql.DB, disabled bool, level, output string, timestamp bool) error {
	query := `
        INSERT INTO log (disabled, level, output, timestamp) 
        VALUES (?, ?, ?, ?)
    `
	_, err := db.Exec(query, disabled, level, output, timestamp)
	if err != nil {
		return fmt.Errorf("error inserting log data: %v", err)
	}
	fmt.Println("Log data added successfully.")
	return nil
}

// AddInbound Function to add an inbound
func AddInbound(
	db *sql.DB,
	inboundType, tag, listen, sniffTimeout string,
	listenPort, transportID, tlsID, realityID, handshakeID int,
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

// InsertTransport inserts a new transport entry into the database.
func InsertTransport(db *sql.DB, transportType, transportPath string) error {
	// Insert transport details into the transports table
	_, err := db.Exec(
		`INSERT INTO transports (type, path) VALUES (?, ?)`,
		transportType, transportPath)

	if err != nil {
		return fmt.Errorf("error inserting into transports table: %v", err)
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
