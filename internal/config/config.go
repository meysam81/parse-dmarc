package config

import (
	"encoding/json"
	"os"
)

// Config holds the application configuration
type Config struct {
	IMAP     IMAPConfig     `json:"imap"`
	Database DatabaseConfig `json:"database"`
	Server   ServerConfig   `json:"server"`
}

// IMAPConfig holds IMAP server configuration
type IMAPConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Mailbox  string `json:"mailbox"`
	UseTLS   bool   `json:"use_tls"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Path string `json:"path"`
}

// ServerConfig holds web server configuration
type ServerConfig struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}

// Load loads configuration from a JSON file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Set defaults
	if cfg.IMAP.Port == 0 {
		cfg.IMAP.Port = 993
	}
	if cfg.IMAP.Mailbox == "" {
		cfg.IMAP.Mailbox = "INBOX"
	}
	if !cfg.IMAP.UseTLS {
		cfg.IMAP.UseTLS = true
	}
	if cfg.Database.Path == "" {
		cfg.Database.Path = "./dmarc.db"
	}
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if cfg.Server.Host == "" {
		cfg.Server.Host = "0.0.0.0"
	}

	return &cfg, nil
}

// GenerateSample creates a sample configuration file
func GenerateSample(path string) error {
	sample := Config{
		IMAP: IMAPConfig{
			Host:     "imap.example.com",
			Port:     993,
			Username: "your-email@example.com",
			Password: "your-password",
			Mailbox:  "INBOX",
			UseTLS:   true,
		},
		Database: DatabaseConfig{
			Path: "./dmarc.db",
		},
		Server: ServerConfig{
			Port: 8080,
			Host: "0.0.0.0",
		},
	}

	data, err := json.MarshalIndent(sample, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
