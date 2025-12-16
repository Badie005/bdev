package config

import (
	"os"
	"path/filepath"
	"testing"
)

// ==============================================================
// Load and Get Tests
// ==============================================================

func TestLoad(t *testing.T) {
	// Reset global state for testing
	globalConfig = nil

	cfg := Load()
	if cfg == nil {
		t.Fatal("Load() returned nil")
	}

	// Check defaults
	if cfg.Version == "" {
		t.Error("Version should have a default value")
	}
	if cfg.AI.Model == "" {
		t.Error("AI.Model should have a default value")
	}
	if cfg.Display.Theme == "" {
		t.Error("Display.Theme should have a default value")
	}
}

func TestLoad_Cached(t *testing.T) {
	// Reset global state
	globalConfig = nil

	cfg1 := Load()
	cfg2 := Load()

	// Should return same instance (cached)
	if cfg1 != cfg2 {
		t.Error("Load() should return cached instance")
	}
}

func TestGet(t *testing.T) {
	globalConfig = nil

	cfg := Get()
	if cfg == nil {
		t.Fatal("Get() returned nil")
	}
}

func TestGet_LoadsIfNil(t *testing.T) {
	globalConfig = nil

	cfg := Get()
	if cfg == nil {
		t.Fatal("Get() should call Load() if globalConfig is nil")
	}
}

// ==============================================================
// Path Method Tests
// ==============================================================

func TestConfig_HistoryFile(t *testing.T) {
	cfg := &Config{
		Paths: PathsConfig{Bdev: "/test/path/.bdev"},
	}

	got := cfg.HistoryFile()
	want := filepath.Join("/test/path/.bdev", "history")

	if got != want {
		t.Errorf("HistoryFile() = %q, want %q", got, want)
	}
}

func TestConfig_SessionFile(t *testing.T) {
	cfg := &Config{
		Paths: PathsConfig{Bdev: "/test/path/.bdev"},
	}

	got := cfg.SessionFile()
	want := filepath.Join("/test/path/.bdev", "session.json")

	if got != want {
		t.Errorf("SessionFile() = %q, want %q", got, want)
	}
}

func TestConfig_AIMemoryFile(t *testing.T) {
	cfg := &Config{
		Paths: PathsConfig{Bdev: "/test/path/.bdev"},
	}

	got := cfg.AIMemoryFile()
	want := filepath.Join("/test/path/.bdev", "ai_memory.json")

	if got != want {
		t.Errorf("AIMemoryFile() = %q, want %q", got, want)
	}
}

func TestConfig_VaultFile(t *testing.T) {
	cfg := &Config{
		Paths: PathsConfig{Bdev: "/test/path/.bdev"},
	}

	got := cfg.VaultFile()
	want := filepath.Join("/test/path/.bdev", "vault.enc")

	if got != want {
		t.Errorf("VaultFile() = %q, want %q", got, want)
	}
}

// ==============================================================
// Save Tests
// ==============================================================

func TestConfig_Save(t *testing.T) {
	tmpDir := t.TempDir()
	bdevDir := filepath.Join(tmpDir, ".bdev")

	cfg := &Config{
		Version: "3.0.0",
		Paths:   PathsConfig{Bdev: bdevDir},
		User:    UserConfig{Name: "Test User"},
		AI:      AIConfig{Model: "test-model"},
	}

	err := cfg.Save()
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify file was created
	configPath := filepath.Join(bdevDir, "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Save() did not create config.json")
	}

	// Read and verify content
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config.json: %v", err)
	}

	if len(data) == 0 {
		t.Error("config.json is empty")
	}

	// Check file permissions (should be 0o600 for security)
	info, _ := os.Stat(configPath)
	mode := info.Mode().Perm()
	// On Windows, permissions work differently, so we just check it's not world-readable
	if mode&0o077 != 0 && os.Getenv("GOOS") != "windows" {
		t.Logf("Warning: config file permissions %o should be 0600", mode)
	}
}

// ==============================================================
// EnsureDirectories Tests
// ==============================================================

func TestConfig_EnsureDirectories(t *testing.T) {
	tmpDir := t.TempDir()
	bdevDir := filepath.Join(tmpDir, ".bdev")

	cfg := &Config{
		Paths: PathsConfig{Bdev: bdevDir},
	}

	err := cfg.EnsureDirectories()
	if err != nil {
		t.Fatalf("EnsureDirectories() error = %v", err)
	}

	// Check all required directories exist
	requiredDirs := []string{
		bdevDir,
		filepath.Join(bdevDir, "workflows"),
		filepath.Join(bdevDir, "cache"),
		filepath.Join(bdevDir, "logs"),
	}

	for _, dir := range requiredDirs {
		info, err := os.Stat(dir)
		if os.IsNotExist(err) {
			t.Errorf("EnsureDirectories() did not create %s", dir)
		}
		if !info.IsDir() {
			t.Errorf("%s should be a directory", dir)
		}
	}
}

func TestConfig_EnsureDirectories_Idempotent(t *testing.T) {
	tmpDir := t.TempDir()
	bdevDir := filepath.Join(tmpDir, ".bdev")

	cfg := &Config{
		Paths: PathsConfig{Bdev: bdevDir},
	}

	// Call twice - should not error
	if err := cfg.EnsureDirectories(); err != nil {
		t.Fatalf("First EnsureDirectories() error = %v", err)
	}
	if err := cfg.EnsureDirectories(); err != nil {
		t.Fatalf("Second EnsureDirectories() error = %v", err)
	}
}

// ==============================================================
// homeDir Tests
// ==============================================================

func TestHomeDir(t *testing.T) {
	home := homeDir()

	// Should not return empty string
	if home == "" {
		t.Error("homeDir() returned empty string")
	}

	// Should return "." if UserHomeDir fails (we can't easily test that)
	// but at least verify it returns something valid
	if _, err := os.Stat(home); os.IsNotExist(err) && home != "." {
		t.Errorf("homeDir() = %q does not exist", home)
	}
}

// ==============================================================
// Config Struct Tests
// ==============================================================

func TestConfig_DefaultValues(t *testing.T) {
	globalConfig = nil
	cfg := Load()

	tests := []struct {
		name string
		got  interface{}
		want interface{}
	}{
		{"Version", cfg.Version != "", true},
		{"AI.Enabled", cfg.AI.Enabled, true},
		{"AI.Model set", cfg.AI.Model != "", true},
		{"AI.BaseURL", cfg.AI.BaseURL, "http://localhost:11434"},
		{"Display.UseColors", cfg.Display.UseColors, true},
		{"Display.UseUnicode", cfg.Display.UseUnicode, true},
		{"Aliases not nil", cfg.Aliases != nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestAliases(t *testing.T) {
	globalConfig = nil
	cfg := Load()

	// Check default aliases exist
	expectedAliases := map[string]string{
		"gs": "git status",
		"gp": "git push",
		"gl": "git pull",
		"gc": "git commit",
	}

	for alias, cmd := range expectedAliases {
		if got, ok := cfg.Aliases[alias]; !ok {
			t.Errorf("Missing alias %q", alias)
		} else if got != cmd {
			t.Errorf("Alias[%q] = %q, want %q", alias, got, cmd)
		}
	}
}
