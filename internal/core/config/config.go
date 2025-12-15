package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the B.DEV CLI configuration
type Config struct {
	Version string            `json:"version" mapstructure:"version"`
	User    UserConfig        `json:"user" mapstructure:"user"`
	Paths   PathsConfig       `json:"paths" mapstructure:"paths"`
	AI      AIConfig          `json:"ai" mapstructure:"ai"`
	Display DisplayConfig     `json:"display" mapstructure:"display"`
	Aliases map[string]string `json:"aliases" mapstructure:"aliases"`
}

// UserConfig contains user preferences
type UserConfig struct {
	Name   string `json:"name" mapstructure:"name"`
	Editor string `json:"editor" mapstructure:"editor"`
	Shell  string `json:"shell" mapstructure:"shell"`
}

// PathsConfig contains path configurations
type PathsConfig struct {
	Projects string `json:"projects" mapstructure:"projects"`
	Bdev     string `json:"bdev" mapstructure:"bdev"`
	Backups  string `json:"backups" mapstructure:"backups"`
}

// AIConfig contains AI engine settings
type AIConfig struct {
	Enabled       bool   `json:"enabled" mapstructure:"enabled"`
	Model         string `json:"model" mapstructure:"model"`
	FallbackModel string `json:"fallback_model" mapstructure:"fallback_model"`
	Timeout       int    `json:"timeout" mapstructure:"timeout"`
	BaseURL       string `json:"base_url" mapstructure:"base_url"`
}

// DisplayConfig contains display settings
type DisplayConfig struct {
	UseColors  bool   `json:"use_colors" mapstructure:"use_colors"`
	UseUnicode bool   `json:"use_unicode" mapstructure:"use_unicode"`
	Theme      string `json:"theme" mapstructure:"theme"`
	DateFormat string `json:"date_format" mapstructure:"date_format"`
}

// Global config instance
var globalConfig *Config

// Load loads the configuration from disk
func Load() *Config {
	if globalConfig != nil {
		return globalConfig
	}

	cfg := &Config{
		Version: "3.0.0",
		User: UserConfig{
			Name:   "B.DEV",
			Editor: "code",
			Shell:  "powershell",
		},
		Paths: PathsConfig{
			Projects: filepath.Join(homeDir(), "Dev", "Projects"),
			Bdev:     filepath.Join(homeDir(), "Dev", ".bdev"),
			Backups:  filepath.Join("D:", "Backups"),
		},
		AI: AIConfig{
			Enabled:       true,
			Model:         "llama3.2",
			FallbackModel: "phi3:mini",
			Timeout:       120,
			BaseURL:       "http://localhost:11434",
		},
		Display: DisplayConfig{
			UseColors:  true,
			UseUnicode: true,
			Theme:      "claude",
			DateFormat: "relative",
		},
		Aliases: map[string]string{
			"gs": "git status",
			"gp": "git push",
			"gl": "git pull",
			"gc": "git commit",
		},
	}

	// Try to load from viper
	if err := viper.Unmarshal(cfg); err == nil {
		globalConfig = cfg
	}

	globalConfig = cfg
	return cfg
}

// Get returns the global config instance
func Get() *Config {
	if globalConfig == nil {
		return Load()
	}
	return globalConfig
}

// HistoryFile returns the path to command history file
func (c *Config) HistoryFile() string {
	return filepath.Join(c.Paths.Bdev, "history")
}

// SessionFile returns the path to session file
func (c *Config) SessionFile() string {
	return filepath.Join(c.Paths.Bdev, "session.json")
}

// AIMemoryFile returns the path to AI memory file
func (c *Config) AIMemoryFile() string {
	return filepath.Join(c.Paths.Bdev, "ai_memory.json")
}

// VaultFile returns the path to secrets vault
func (c *Config) VaultFile() string {
	return filepath.Join(c.Paths.Bdev, "vault.enc")
}

// Save persists the configuration to disk
func (c *Config) Save() error {
	configPath := filepath.Join(c.Paths.Bdev, "config.json")

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0o600)
}

// homeDir returns the user's home directory
func homeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return home
}

// EnsureDirectories creates all required B.DEV directories
func (c *Config) EnsureDirectories() error {
	dirs := []string{
		c.Paths.Bdev,
		filepath.Join(c.Paths.Bdev, "workflows"),
		filepath.Join(c.Paths.Bdev, "cache"),
		filepath.Join(c.Paths.Bdev, "logs"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	return nil
}
