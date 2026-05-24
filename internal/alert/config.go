package alert

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFile = "alert_config.json"

var configDir = func() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".confsnap")
}()

// Config holds alert dispatch configuration.
type Config struct {
	Enabled      bool   `json:"enabled"`
	Level        string `json:"level"`
	LogFile      string `json:"log_file,omitempty"`
	Stdout       bool   `json:"stdout"`
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		Enabled: true,
		Level:   "WARNING",
		Stdout:  true,
	}
}

// SaveConfig persists the alert config to disk.
func SaveConfig(cfg Config) error {
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("alert config: mkdir: %w", err)
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("alert config: marshal: %w", err)
	}
	return os.WriteFile(filepath.Join(configDir, configFile), data, 0644)
}

// LoadConfig reads the alert config from disk, returning defaults if not found.
func LoadConfig() (Config, error) {
	path := filepath.Join(configDir, configFile)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return DefaultConfig(), nil
	}
	if err != nil {
		return Config{}, fmt.Errorf("alert config: read: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("alert config: unmarshal: %w", err)
	}
	return cfg, nil
}

// BuildDispatcher constructs a Dispatcher from the given Config.
func BuildDispatcher(cfg Config) *Dispatcher {
	var handlers []Handler
	if cfg.Stdout {
		handlers = append(handlers, &StdoutHandler{})
	}
	if cfg.LogFile != "" {
		handlers = append(handlers, &FileHandler{Path: cfg.LogFile})
	}
	return NewDispatcher(handlers...)
}
