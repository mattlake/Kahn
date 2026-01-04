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

// Default configuration constants
const (
	DefaultDatabasePath = "~/.kahn/kahn.db"
	DefaultBusyTimeout  = 5000 // milliseconds
	DefaultJournalMode  = "WAL"
	DefaultCacheSize    = 10000 // number of pages
	DefaultForeignKeys  = true
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
	viper.SetDefault("database.path", DefaultDatabasePath)
	viper.SetDefault("database.busy_timeout", DefaultBusyTimeout)
	viper.SetDefault("database.journal_mode", DefaultJournalMode)
	viper.SetDefault("database.cache_size", DefaultCacheSize)
	viper.SetDefault("database.foreign_keys", DefaultForeignKeys)

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
