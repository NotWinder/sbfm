// Package jsonhandler is responsible for generating the config.json file.
package jsonhandler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	// go-sqlite3 is the SQL driver for SQLite in Go
	_ "github.com/mattn/go-sqlite3"
)

// Config is the structure of the config.json file.
type Config struct {
	Log          Log          `json:"log,omitempty"`
	DNS          DNS          `json:"dns,omitempty"`
	NTP          NTP          `json:"ntp,omitempty"`
	Inbounds     []Inbound    `json:"inbounds,omitempty"`
	Outbounds    []Outbound   `json:"outbounds,omitempty"`
	Route        Route        `json:"route,omitempty"`
	Experimental Experimental `json:"experimental,omitempty"`
}

// Log is the structure of the Log block.
type Log struct {
	Disabled  bool   `json:"disabled,omitempty"`
	Level     string `json:"level,omitempty"`
	Output    string `json:"output,omitempty"`
	Timestamp bool   `json:"timestamp,omitempty"`
}

// DNS is the structure of the DNS block.
// Add DNS fields as needed.
type DNS struct{}

// NTP is the structure of the NTP block.
// Add NTP fields as needed.
type NTP struct{}

// Inbound is the structure of the Inbound block.
type Inbound struct {
	Type                      string    `json:"type"`
	Tag                       string    `json:"tag"`
	Listen                    string    `json:"listen"`
	ListenPort                int       `json:"listen_port"`
	TCPFastOpen               bool      `json:"tcp_fast_open,omitempty"`
	TCPMultiPath              bool      `json:"tcp_multi_path,omitempty"`
	UDPFragment               bool      `json:"udp_fragment,omitempty"`
	UDPTimeout                string    `json:"udp_timeout,omitempty"`
	Detour                    string    `json:"detour,omitempty"`
	Sniff                     bool      `json:"sniff"`
	SniffOverrideDestination  bool      `json:"sniff_override_destination"`
	SniffTimeout              string    `json:"sniff_timeout"`
	DomainStrategy            string    `json:"domain_strategy,omitempty"`
	UDPDisableDomainUnmapping bool      `json:"udp_disable_domain_unmapping,omitempty"`
	Users                     []User    `json:"users"`
	TLS                       TLS       `json:"tls"`
	Transport                 Transport `json:"transport"`
}

// User is the structure of the user block in the inbound block.
type User struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

// TLS is the structure of the TLS block in the inbound block.
// Add TLS fields as needed.
type TLS struct{}

// Transport is the structure of the transport block in the inbound block.
type Transport struct {
	Type string `json:"type"`
	Path string `json:"path"`
}

// Outbound is the structure of the Outbound block.
// Add Outbound fields as needed.
type Outbound struct{}

// Route is the structure of the Route block.
// Add Route fields as needed.
type Route struct{}

// Experimental is the structure of the Experimental block.
// Add Experimental fields as needed.
type Experimental struct{}

// GenerateConfigFile generates the config.json file from the data in the database.
func GenerateConfigFile(db *sql.DB) error {
	// Create a Config instance.
	config := Config{}

	// Populate the Config instance from the database.
	err := PopulateConfig(db, &config)
	if err != nil {
		return fmt.Errorf("error populating config: %v", err)
	}

	// Generate JSON.
	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	// Write JSON to file.
	err = os.WriteFile("config.json", jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing JSON to file: %v", err)
	}

	fmt.Println("Config file generated successfully!")
	return nil
}

// PopulateConfig populates the config.json file from the data in the database.
func PopulateConfig(db *sql.DB, config *Config) error {
	// Populate Log
	err := db.QueryRow("SELECT disabled, level, output, timestamp FROM log").Scan(
		&config.Log.Disabled, &config.Log.Level, &config.Log.Output, &config.Log.Timestamp)
	if err != nil {
		return fmt.Errorf("error querying log table: %v", err)
	}

	fmt.Println("inbound1")
	// Fetch all users once
	var allUsers []User
	userRows, err := db.Query("SELECT name, uuid FROM users")
	if err != nil {
		return fmt.Errorf("error querying users table: %v", err)
	}
	defer userRows.Close()

	fmt.Println("inbound2")

	for userRows.Next() {
		var user User
		err := userRows.Scan(&user.Name, &user.UUID)
		if err != nil {
			return fmt.Errorf("error scanning user row: %v", err)
		}
		allUsers = append(allUsers, user)
	}

	fmt.Println("inbound3")

	// Populate Inbounds with Transport details
	rows, err := db.Query(`
	SELECT 
		i.type, i.tag, i.listen, i.listen_port, i.tcp_fast_open, i.tcp_multi_path, 
		i.udp_fragment, i.udp_timeout, i.detour, i.sniff, i.sniff_override_destination, 
		i.sniff_timeout, i.domain_strategy, i.udp_disable_domain_unmapping, 
		t.type AS transport_type, t.path
	FROM inbounds i
	LEFT JOIN transports t ON i.transport_id = t.id
`)
	if err != nil {
		return fmt.Errorf("error querying inbounds table: %v", err)
	}
	defer rows.Close()

	fmt.Println("inbound4")
	for rows.Next() {
		var inbound Inbound

		// Using sql.Null* types for optional fields
		var tcpFastOpen, tcpMultiPath, udpFragment, udpDisableDomainUnmapping sql.NullBool
		var udpTimeout, detour, domainStrategy sql.NullString
		var transportType, transportPath sql.NullString

		err := rows.Scan(
			&inbound.Type,
			&inbound.Tag,
			&inbound.Listen,
			&inbound.ListenPort,
			&tcpFastOpen,
			&tcpMultiPath,
			&udpFragment,
			&udpTimeout,
			&detour,
			&inbound.Sniff,
			&inbound.SniffOverrideDestination,
			&inbound.SniffTimeout,
			&domainStrategy,
			&udpDisableDomainUnmapping,
			&transportType,
			&transportPath,
		)
		if err != nil {
			return fmt.Errorf("error scanning inbound row: %v", err)
		}

		// Set optional fields if they are valid (i.e., not NULL)
		inbound.TCPFastOpen = tcpFastOpen.Valid && tcpFastOpen.Bool
		inbound.TCPMultiPath = tcpMultiPath.Valid && tcpMultiPath.Bool
		inbound.UDPFragment = udpFragment.Valid && udpFragment.Bool
		if udpTimeout.Valid {
			inbound.UDPTimeout = udpTimeout.String
		}
		if detour.Valid {
			inbound.Detour = detour.String
		}
		if domainStrategy.Valid {
			inbound.DomainStrategy = domainStrategy.String
		}
		inbound.UDPDisableDomainUnmapping = udpDisableDomainUnmapping.Valid &&
			udpDisableDomainUnmapping.Bool

		// Assign all users to the inbound entry.
		inbound.Users = allUsers

		// Populate Transport block
		if transportType.Valid {
			inbound.Transport.Type = transportType.String
		}
		if transportPath.Valid {
			inbound.Transport.Path = transportPath.String
		}

		// Debugging print to check if transport values are coming through
		fmt.Printf("Inbound Transport: Type=%s, Path=%s\n",
			inbound.Transport.Type, inbound.Transport.Path)

		// Add the inbound to the config
		config.Inbounds = append(config.Inbounds, inbound)
	}
	fmt.Println("inbound5")

	return nil
}
