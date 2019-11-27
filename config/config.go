package config

import (
	"encoding/json"
	"io/ioutil"
)

// ApiServer represents the API server configuration details
type ApiServer struct {
	CorsAllowedOrigins []string `json:"cors_allowed_origins"`
	CorsDebug          bool     `json:"cors_debug"`
}

// AdminServer represents the Admin server configuration details
type AdminServer struct {
	ListenURL string    `json:"listen_url"`
	UseTLS    bool      `json:"use_tls"`
	CertPath  string    `json:"cert_path"`
	KeyPath   string    `json:"key_path"`
	ApiConf   ApiServer `json:"api"`
}

// PhishServer represents the Phish server configuration details
type PhishServer struct {
	ListenURL string `json:"listen_url"`
	UseTLS    bool   `json:"use_tls"`
	CertPath  string `json:"cert_path"`
	KeyPath   string `json:"key_path"`
}

// LoggingConfig represents configuration details for Gophish logging.
type LoggingConfig struct {
	Filename string `json:"filename"`
}

// Config represents the configuration information.
type Config struct {
	AdminConf      AdminServer   `json:"admin_server"`
	PhishConf      PhishServer   `json:"phish_server"`
	DBName         string        `json:"db_name"`
	DBPath         string        `json:"db_path"`
	DBSSLCaPath    string        `json:"db_sslca_path"`
	MigrationsPath string        `json:"migrations_prefix"`
	TestFlag       bool          `json:"test_flag"`
	ContactAddress string        `json:"contact_address"`
	Logging        LoggingConfig `json:"logging"`
}

// Version contains the current gophish version
var Version = ""

// ServerName is the server type that is returned in the transparency response.
const ServerName = "gophish"

// LoadConfig loads the configuration from the specified filepath
func LoadConfig(filepath string) (*Config, error) {
	// Get the config file
	configFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = json.Unmarshal(configFile, config)
	if err != nil {
		return nil, err
	}
	// Choosing the migrations directory based on the database used.
	config.MigrationsPath = config.MigrationsPath + config.DBName
	// Explicitly set the TestFlag to false to prevent config.json overrides
	config.TestFlag = false
	return config, nil
}
