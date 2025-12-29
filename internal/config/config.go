package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
		Path        string `mapstructure:"path"`
		BusyTimeout int    `mapstructure:"busy_timeout"`
		JournalMode string `mapstructure:"journal_mode"`
		CacheSize   int    `mapstructure:"cache_size"`
		ForeignKeys bool   `mapstructure:"foreign_keys"`
	} `mapstructure:"database"`
}

/*
* LoadConfig loads configuration from multiple sources with the following priority:
* 1. Command-line flags (highest priority)
* 2. Environment variables
* 3. Config file
* 4. Default values (lowest priority)
 */
func LoadConfig() (*Config, error) {
	config := &Config{}

	// Set default values
	viper.SetDefault("database.path", "~/.kahn/kahn.db")
	viper.SetDefault("database.busy_timeout", 5000)
	viper.SetDefault("database.journal_mode", "WAL")
	viper.SetDefault("database.cache_size", 10000)
	viper.SetDefault("database.foreign_keys", true)

	// Set up command-line flags
	pflag.String("config", "", "Path to config file")
	pflag.String("db-path", "", "Path to database file")
	pflag.Parse()

	// Bind command-line flags to viper
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return nil, fmt.Errorf("failed to bind flags: %w", err)
	}

	// Bind db-path flag to database.path
	err = viper.BindPFlag("database.path", pflag.Lookup("db-path"))
	if err != nil {
		return nil, fmt.Errorf("failed to bind db-path flag: %w", err)
	}

	// Set up environment variable prefix
	viper.SetEnvPrefix("KAHN")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Set up config file search paths
	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	// Add search paths in order of preference
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.kahn")
	viper.AddConfigPath("/etc/kahn")

	// If a config file is specified via flag, use it
	if configFile := viper.GetString("config"); configFile != "" {
		viper.SetConfigFile(configFile)
	}

	// Try to read config file (ignore if not found)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	config.Database.Path = expandPath(config.Database.Path)

	return config, nil
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Printf("Warning: could not determine home directory: %v", err)
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}

func EnsureConfigDir(configPath string) error {
	configDir := filepath.Dir(configPath)
	if configDir == "." {
		return nil
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory %s: %w", configDir, err)
	}

	return nil
}

func WriteDefaultConfig(configPath string) error {
	defaultConfig := `# Kahn Configuration File
# This file controls the behavior of Kahn task manager

[database]
# Path to the SQLite database file
# Supports ~ expansion for home directory
path = "~/.kahn/kahn.db"

# Database connection timeout in milliseconds
busy_timeout = 5000

# Journal mode for SQLite (WAL recommended for better concurrency)
# Options: DELETE, TRUNCATE, PERSIST, MEMORY, WAL, OFF
journal_mode = "WAL"

# SQLite cache size (number of pages)
cache_size = 10000

# Enable foreign key constraints
foreign_keys = true
`

	if _, err := os.Stat(configPath); err == nil {
		return fmt.Errorf("config file already exists: %s", configPath)
	}

	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(configPath, []byte(defaultConfig), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
