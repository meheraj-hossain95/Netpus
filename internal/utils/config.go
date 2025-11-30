package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

// Config represents application configuration
type Config struct {
	AutoStart        bool   `json:"autoStart"`
	Theme            string `json:"theme"`
	DataRetention    int    `json:"dataRetention"`    // Days to keep data, 0 = forever
	NetworkInterface string `json:"networkInterface"` // Reserved for future use
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		AutoStart:        false,
		Theme:            "auto",
		DataRetention:    30,
		NetworkInterface: "",
	}
}

// LoadConfig loads configuration from database
func LoadConfig(db interface{}) (*Config, error) {
	type SettingsDB interface {
		GetSetting(key string) (string, error)
	}

	sdb, ok := db.(SettingsDB)
	if !ok {
		return DefaultConfig(), fmt.Errorf("invalid database interface")
	}

	config := DefaultConfig()

	if val, err := sdb.GetSetting("autoStart"); err == nil && val != "" {
		config.AutoStart = val == "true"
	}

	if val, err := sdb.GetSetting("theme"); err == nil && val != "" {
		config.Theme = val
	}

	if val, err := sdb.GetSetting("dataRetention"); err == nil && val != "" {
		if days, err := strconv.Atoi(val); err == nil {
			config.DataRetention = days
		}
	}

	if val, err := sdb.GetSetting("networkInterface"); err == nil && val != "" {
		config.NetworkInterface = val
	}

	return config, nil
}

// Save saves configuration to database
func (c *Config) Save(db interface{}) error {
	type SettingsDB interface {
		SetSetting(key, value string) error
	}

	sdb, ok := db.(SettingsDB)
	if !ok {
		return fmt.Errorf("invalid database interface")
	}

	if err := sdb.SetSetting("autoStart", strconv.FormatBool(c.AutoStart)); err != nil {
		return err
	}

	if err := sdb.SetSetting("theme", c.Theme); err != nil {
		return err
	}

	if err := sdb.SetSetting("dataRetention", strconv.Itoa(c.DataRetention)); err != nil {
		return err
	}

	if err := sdb.SetSetting("networkInterface", c.NetworkInterface); err != nil {
		return err
	}

	return nil
}

// GetDatabasePath returns the platform-specific database path
func GetDatabasePath() string {
	var basePath string

	if runtime.GOOS == "windows" {
		basePath = os.Getenv("APPDATA")
		if basePath == "" {
			basePath = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming")
		}
	} else {
		basePath = os.Getenv("XDG_DATA_HOME")
		if basePath == "" {
			basePath = filepath.Join(os.Getenv("HOME"), ".local", "share")
		}
	}

	return filepath.Join(basePath, "netpus", "netpus.db")
}

// GetExecutablePath returns the current executable path
func GetExecutablePath() (string, error) {
	return os.Executable()
}
