package config

import (
	"os"
	"path/filepath"

	"github.com/caarlos0/env/v11"
	"github.com/goccy/go-json"
)

// Config holds the application configuration
type Config struct {
	IMAP        IMAPConfig     `json:"imap"`
	IMAPConfigs []IMAPConfig   `json:"imap_configs"`
	Database    DatabaseConfig `json:"database"`
	Server      ServerConfig   `json:"server"`
}

// IMAPConfig holds IMAP server configuration
type IMAPConfig struct {
	Host     string `json:"host" env:"IMAP_HOST"`
	Port     int    `json:"port" env:"IMAP_PORT" envDefault:"993"`
	Username string `json:"username" env:"IMAP_USERNAME"`
	Password string `json:"password" env:"IMAP_PASSWORD"`
	Mailbox  string `json:"mailbox" env:"IMAP_MAILBOX" envDefault:"INBOX"`
	UseTLS   bool   `json:"use_tls" env:"IMAP_USE_TLS" envDefault:"true"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Path string `json:"path" env:"DATABASE_PATH"`
}

// ServerConfig holds web server configuration
type ServerConfig struct {
	Port int    `json:"port" env:"SERVER_PORT" envDefault:"8080"`
	Host string `json:"host" env:"SERVER_HOST" envDefault:"0.0.0.0"`
}

// GetIMAPConfigs returns all IMAP configurations, normalizing single and multiple configs
func (c *Config) GetIMAPConfigs() []IMAPConfig {
	// If IMAPConfigs array is populated, use it
	if len(c.IMAPConfigs) > 0 {
		return c.IMAPConfigs
	}

	// Otherwise, use the single IMAP config for backward compatibility
	if c.IMAP.Host != "" {
		return []IMAPConfig{c.IMAP}
	}

	return []IMAPConfig{}
}

func defaultDBPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".parse-dmarc/db.sqlite"), nil
}

func ensureDBPathExists(dbPath string) error {
	parent := filepath.Dir(dbPath)
	err := os.MkdirAll(parent, 0755)
	if err != nil {
		return err
	}
	return nil
}

// Load loads configuration from a JSON file
func Load(path string) (*Config, error) {
	var cfg Config
	var err error

	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, err
		}
	}

	// Support IMAP_CONFIGS environment variable for multiple configs
	imapConfigsJSON := os.Getenv("IMAP_CONFIGS")
	if imapConfigsJSON != "" {
		var imapConfigs []IMAPConfig
		if err := json.Unmarshal([]byte(imapConfigsJSON), &imapConfigs); err != nil {
			return nil, err
		}
		cfg.IMAPConfigs = imapConfigs
	}

	// Only parse environment variables if we're using single IMAP config
	// (env.Parse incorrectly applies envDefault to array elements)
	if len(cfg.IMAPConfigs) == 0 {
		if err := env.Parse(&cfg); err != nil {
			return nil, err
		}
	} else {
		// Parse only non-IMAP fields to avoid envDefault overwriting array values
		if err := env.Parse(&cfg.Database); err != nil {
			return nil, err
		}
		if err := env.Parse(&cfg.Server); err != nil {
			return nil, err
		}
	}

	// Apply defaults to all IMAP configs
	allConfigs := cfg.GetIMAPConfigs()
	for i := range allConfigs {
		if allConfigs[i].Port == 0 {
			allConfigs[i].Port = 993
		}
		if allConfigs[i].Mailbox == "" {
			allConfigs[i].Mailbox = "INBOX"
		}
	}

	// Update the config with normalized values
	if len(cfg.IMAPConfigs) > 0 {
		cfg.IMAPConfigs = allConfigs
	} else if cfg.IMAP.Host != "" {
		cfg.IMAP = allConfigs[0]
	}
	if cfg.Database.Path == "" {
		cfg.Database.Path, err = defaultDBPath()
		if err != nil {
			return nil, err
		}
	}
	err = ensureDBPathExists(cfg.Database.Path)
	if err != nil {
		return nil, err
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
	dbPath, err := defaultDBPath()
	if err != nil {
		return err
	}
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
			Path: dbPath,
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
