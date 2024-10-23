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
type DNS struct {
	Servers          []Servers `json:"servers,omitempty"`
	Rules            []Rules   `json:"rules,omitempty"`
	Final            string    `json:"final,omitempty"`
	Strategy         string    `json:"strategy,omitempty"`
	DisableCache     bool      `json:"disable_cache,omitempty"`
	DisableExpire    bool      `json:"disable_expire,omitempty"`
	IndependentCache bool      `json:"independent_cache,omitempty"`
	ReverseMapping   bool      `json:"reverse_mapping,omitempty"`
	ClientSubnet     string    `json:"client_subnet,omitempty"`
	Fakeip           FakeIP    `json:"fakeip,omitempty"`
}

// Servers is the structure of the DNS servers block.
// Add Servers fields as needed.
type Servers struct {
	Tag             string `json:"tag,omitempty"`
	Address         string `json:"address,omitempty"`
	AddressResolver string `json:"address_resolver,omitempty"`
	AddressStrategy string `json:"address_strategy,omitempty"`
	Strategy        string `json:"strategy,omitempty"`
	Detour          string `json:"detour,omitempty"`
	ClientSubnet    string `json:"client_subnet,omitempty"`
}

// Rules is the structure of the DNS rules block.
// Add Rules fields as needed.
type Rules struct{}

// FakeIP is the structure of the DNS fakeip block.
// Add Rules fields as needed.
type FakeIP struct {
	Enabled    bool   `json:"enabled,omitempty"`
	Inet4Range string `json:"inet4_range,omitempty"`
	Inet6Range string `json:"inet6_range,omitempty"`
}

// NTP is the structure of the NTP block.
// Add NTP fields as needed.
type NTP struct {
	Enabled    bool   `json:"enabled,omitempty"`
	Server     string `json:"server,omitempty"`
	ServerPort int    `json:"server_port,omitempty"`
	Interval   string `json:"interval,omitempty"`
	// Dial Fields
}

// Inbound is the structure of the Inbound block.
type Inbound struct {
	Type                      string    `json:"type,omitempty"`
	Tag                       string    `json:"tag,omitempty"`
	Listen                    string    `json:"listen,omitempty"`
	ListenPort                int       `json:"listen_port,omitempty"`
	TCPFastOpen               bool      `json:"tcp_fast_open,omitempty"`
	TCPMultiPath              bool      `json:"tcp_multi_path,omitempty"`
	UDPFragment               bool      `json:"udp_fragment,omitempty"`
	UDPTimeout                string    `json:"udp_timeout,omitempty"`
	Detour                    string    `json:"detour,omitempty"`
	Sniff                     bool      `json:"sniff,omitempty"`
	SniffOverrideDestination  bool      `json:"sniff_override_destination"`
	SniffTimeout              string    `json:"sniff_timeout,omitempty"`
	DomainStrategy            string    `json:"domain_strategy,omitempty"`
	UDPDisableDomainUnmapping bool      `json:"udp_disable_domain_unmapping,omitempty"`
	Users                     []User    `json:"users,omitempty"`
	TLS                       TLS       `json:"tls,omitempty"`
	Transport                 Transport `json:"transport,omitempty"`
}

// User is the structure of the user block in the inbound block.
type User struct {
	Name string `json:"name,omitempty"`
	UUID string `json:"uuid,omitempty"`
	SUB  string `json:"sub,omitempty"`
}

// TLS is the structure of the TLS block in the inbound block.
// Add TLS fields as needed.
type TLS struct {
	Enabled         bool    `json:"enabled,omitempty"`
	ServerName      string  `json:"server_name,omitempty"`
	MinVersion      string  `json:"min_version,omitempty"`
	MaxVersion      string  `json:"max_version,omitempty"`
	CertificatePath string  `json:"certificate_path,omitempty"`
	KeyPath         string  `json:"key_path,omitempty"`
	Reality         Reality `json:"reality,omitempty"`
}

// Reality is the structure of the Reality block in the inbound block.
// Add Reality fields as needed.
type Reality struct {
	Enabled    bool      `json:"enabled,omitempty"`
	Handshake  Handshake `json:"handshake,omitempty"`
	PrivateKey string    `json:"private_key,omitempty"`
	ShortID    string    `json:"short_id,omitempty"`
}

// Handshake is the structure of the Handshake block in the inbound block.
// Add Handshake fields as needed.
type Handshake struct {
	Server     string `json:"server,omitempty"`
	ServerPort int    `json:"server_port,omitempty"`
}

// Transport is the structure of the transport block in the inbound block.
type Transport struct {
	Type string `json:"type,omitempty"`
	Path string `json:"path,omitempty"`
}

// Outbound is the structure of the Outbound block.
// Add Outbound fields as needed.
type Outbound struct{}

// Route is the structure of the Route block.
// Add Route fields as needed.
type Route struct {
	Rules               []Rules   `json:"rules,omitempty"`
	RuleSet             []RuleSet `json:"rule_set,omitempty"`
	Final               string    `json:"final,omitempty"`
	AutoDetectInterface bool      `json:"auto_detect_interface,omitempty"`
	OverrideAndroidVPN  bool      `json:"override_android_vpn,omitempty"`
	DefaultInterface    string    `json:"default_interface,omitempty"`
	DefaultMark         int       `json:"default_mark,omitempty"`
}

// RuleSet is the structure of the RuleSet block.
// Add RuleSet fields as needed.
type RuleSet struct{}

// Experimental is the structure of the Experimental block.
// Add Experimental fields as needed.
type Experimental struct {
	CacheFile CacheFile `json:"cache_file,omitempty"`
	ClashAPI  ClashAPI  `json:"clash_api,omitempty"`
	V2rayAPI  V2rayAPI  `json:"v2ray_api,omitempty"`
}

// CacheFile is the structure of the CacheFile block.
// Add CacheFile fields as needed.
type CacheFile struct{}

// ClashAPI is the structure of the ClashAPI block.
// Add ClashAPI fields as needed.
type ClashAPI struct{}

// V2rayAPI is the structure of the V2rayAPI block.
// Add V2rayAPI fields as needed.
type V2rayAPI struct{}

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

	// Fetch all users once
	var allUsers []User
	userRows, err := db.Query("SELECT name, uuid FROM users")
	if err != nil {
		return fmt.Errorf("error querying users table: %v", err)
	}
	defer userRows.Close()

	for userRows.Next() {
		var user User
		err := userRows.Scan(&user.Name, &user.UUID)
		if err != nil {
			return fmt.Errorf("error scanning user row: %v", err)
		}
		allUsers = append(allUsers, user)
	}

	// Query to fetch inbounds, transports, tls, reality, and handshake data
	rows, err := db.Query(`
    SELECT 
        i.type, i.tag, i.listen, i.listen_port, i.tcp_fast_open, i.tcp_multi_path, 
        i.udp_fragment, i.udp_timeout, i.detour, i.sniff, i.sniff_override_destination, 
        i.sniff_timeout, i.domain_strategy, i.udp_disable_domain_unmapping, 
        t.type AS transport_type, t.path,
        tls.enabled, tls.server_name, tls.min_version, tls.max_version, 
        tls.certificate_path, tls.key_path,
        r.enabled AS reality_enabled, r.private_key, r.short_id,
        h.server, h.server_port
    FROM inbounds i
    LEFT JOIN transports t ON i.transport_id = t.id
    LEFT JOIN tls ON i.tls_id = tls.id
    LEFT JOIN reality r ON reality_id = r.id
    LEFT JOIN handshake h ON handshake_id = h.id
`)

	// Check if the query resulted in an error
	if err != nil {
		fmt.Printf("Error querying inbounds table: %v\n", err)
		return fmt.Errorf("error querying inbounds table: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var inbound Inbound

		// Using sql.Null* types for optional fields
		var udpDisableDomainUnmapping, tcpFastOpen, tcpMultiPath, udpFragment, tlsEnabled, realityEnabled sql.NullBool
		var serverName, minVersion, maxVersion, certPath, keyPath, privateKey, shortID, handshakeServer, udpTimeout, detour, domainStrategy, transportType, transportPath sql.NullString
		var handshakeServerPort sql.NullInt64

		err := rows.Scan(
			&inbound.Type, &inbound.Tag, &inbound.Listen, &inbound.ListenPort,
			&tcpFastOpen, &tcpMultiPath, &udpFragment, &udpTimeout, &detour,
			&inbound.Sniff, &inbound.SniffOverrideDestination, &inbound.SniffTimeout,
			&domainStrategy, &udpDisableDomainUnmapping,
			&transportType, &transportPath,
			&tlsEnabled, &serverName, &minVersion, &maxVersion, &certPath, &keyPath,
			&realityEnabled, &privateKey, &shortID,
			&handshakeServer, &handshakeServerPort,
		)
		if err != nil {
			return fmt.Errorf("error scanning inbound row: %v", err)
		}

		// Populate TLS block
		inbound.TLS.Enabled = tlsEnabled.Valid && tlsEnabled.Bool
		if serverName.Valid {
			inbound.TLS.ServerName = serverName.String
		}
		if minVersion.Valid {
			inbound.TLS.MinVersion = minVersion.String
		}
		if maxVersion.Valid {
			inbound.TLS.MaxVersion = maxVersion.String
		}
		if certPath.Valid {
			inbound.TLS.CertificatePath = certPath.String
		}
		if keyPath.Valid {
			inbound.TLS.KeyPath = keyPath.String
		}

		// Populate Reality block
		inbound.TLS.Reality.Enabled = realityEnabled.Valid && realityEnabled.Bool
		if privateKey.Valid {
			inbound.TLS.Reality.PrivateKey = privateKey.String
		}
		if shortID.Valid {
			inbound.TLS.Reality.ShortID = shortID.String
		}

		// Populate Handshake block in Reality
		if handshakeServer.Valid {
			inbound.TLS.Reality.Handshake.Server = handshakeServer.String
		}
		if handshakeServerPort.Valid {
			inbound.TLS.Reality.Handshake.ServerPort = int(handshakeServerPort.Int64)
		}

		// Populate Transport block
		if transportType.Valid {
			inbound.Transport.Type = transportType.String
		}
		if transportPath.Valid {
			inbound.Transport.Path = transportPath.String
		}

		// Assign all users to the inbound entry.
		inbound.Users = allUsers

		config.Inbounds = append(config.Inbounds, inbound)
	}

	return nil
}
