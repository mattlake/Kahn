package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test helper functions - moved from main config package as they're only used in tests

// EnsureConfigDir creates the config directory if it doesn't exist
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

// WriteDefaultConfig writes a default configuration file
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

func TestLoadConfig_DefaultValues(t *testing.T) {
	// Clear any existing environment variables that might interfere
	os.Unsetenv("KAHN_DATABASE_PATH")
	os.Unsetenv("KAHN_DATABASE_BUSY_TIMEOUT")
	os.Unsetenv("KAHN_DATABASE_JOURNAL_MODE")
	os.Unsetenv("KAHN_DATABASE_CACHE_SIZE")
	os.Unsetenv("KAHN_DATABASE_FOREIGN_KEYS")

	config, err := LoadConfig()
	require.NoError(t, err, "LoadConfig should not return error with defaults")
	require.NotNil(t, config, "Config should not be nil")

	// Test default values (path gets expanded)
	expectedPath := filepath.Join(os.Getenv("HOME"), ".kahn", "kahn.db")
	assert.Equal(t, expectedPath, config.Database.Path, "Default database path should match expanded version")
	assert.Equal(t, DefaultBusyTimeout, config.Database.BusyTimeout, "Default busy timeout should match")
	assert.Equal(t, DefaultJournalMode, config.Database.JournalMode, "Default journal mode should match")
	assert.Equal(t, DefaultCacheSize, config.Database.CacheSize, "Default cache size should match")
	assert.Equal(t, DefaultForeignKeys, config.Database.ForeignKeys, "Default foreign keys should be true")
}

func TestExpandPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Path with tilde",
			input:    "~/test.db",
			expected: os.Getenv("HOME") + "/test.db",
		},
		{
			name:     "Path without tilde",
			input:    "/tmp/test.db",
			expected: "/tmp/test.db",
		},
		{
			name:     "Empty path",
			input:    "",
			expected: "",
		},
		{
			name:     "Path with multiple tildes",
			input:    "~/~/test.db",
			expected: os.Getenv("HOME") + "/~/test.db",
		},
		{
			name:     "Absolute path",
			input:    "/absolute/path/test.db",
			expected: "/absolute/path/test.db",
		},
		{
			name:     "Relative path",
			input:    "relative/path/test.db",
			expected: "relative/path/test.db",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandPath(tt.input)
			assert.Equal(t, tt.expected, result, "Path expansion should match expected")
		})
	}
}

func TestEnsureConfigDir(t *testing.T) {
	tests := []struct {
		name        string
		configPath  string
		expectError bool
	}{
		{
			name:        "Create new directory",
			configPath:  filepath.Join(os.TempDir(), "kahn_test_config", "test.db"),
			expectError: false,
		},
		{
			name:        "Use existing directory",
			configPath:  filepath.Join(os.TempDir(), "kahn_test_config", "test2.db"),
			expectError: false,
		},
		{
			name:        "Root directory (should not fail)",
			configPath:  "/tmp/test.db",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := EnsureConfigDir(tt.configPath)

			if tt.expectError {
				assert.Error(t, err, "Should return error")
			} else {
				assert.NoError(t, err, "Should not return error")

				// Verify directory exists
				dir := filepath.Dir(tt.configPath)
				if dir != "/" && dir != "." {
					_, err := os.Stat(dir)
					assert.NoError(t, err, "Directory should exist")
				}
			}
		})
	}
}

func TestWriteDefaultConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "default_config.toml")

	err := WriteDefaultConfig(configPath)
	require.NoError(t, err, "WriteDefaultConfig should not return error")

	// Verify file was created
	_, err = os.Stat(configPath)
	assert.NoError(t, err, "Config file should exist")

	// Verify file content
	content, err := os.ReadFile(configPath)
	require.NoError(t, err, "Should be able to read config file")

	contentStr := string(content)
	assert.Contains(t, contentStr, "[database]", "Config should contain database section")
	assert.Contains(t, contentStr, "path", "Config should contain path setting")
	assert.Contains(t, contentStr, "busy_timeout", "Config should contain busy_timeout setting")
	assert.Contains(t, contentStr, "journal_mode", "Config should contain journal_mode setting")
	assert.Contains(t, contentStr, "cache_size", "Config should contain cache_size setting")
	assert.Contains(t, contentStr, "foreign_keys", "Config should contain foreign_keys setting")
}

func TestWriteDefaultConfig_ExistingFile(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "existing_config.toml")

	// Create existing file
	existingContent := "# Existing config"
	err := os.WriteFile(configPath, []byte(existingContent), 0644)
	require.NoError(t, err, "Should be able to create existing file")

	// Try to write default config (should not overwrite)
	err = WriteDefaultConfig(configPath)
	assert.Error(t, err, "WriteDefaultConfig should return error for existing file")
	assert.Contains(t, err.Error(), "already exists", "Error should mention file exists")

	// Verify original content is unchanged
	content, err := os.ReadFile(configPath)
	require.NoError(t, err, "Should be able to read config file")
	assert.Equal(t, existingContent, string(content), "Original content should be unchanged")
}
