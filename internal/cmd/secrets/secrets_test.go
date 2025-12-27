package secretscmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/badie/bdev/internal/core/vault"
)

func TestVaultIntegration(t *testing.T) {
	// Setup temporary directory
	tmpDir, err := os.MkdirTemp("", "vault_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	vaultPath := filepath.Join(tmpDir, "vault.json")
	v := vault.New(vaultPath)

	password := "strongpassword"

	// Test Create
	if err := v.Create(password); err != nil {
		t.Fatalf("Failed to create vault: %v", err)
	}

	if !v.Exists() {
		t.Error("Vault should exist after creation")
	}

	// Test Unlock
	if err := v.Unlock(password); err != nil {
		t.Fatalf("Failed to unlock vault: %v", err)
	}

	if !v.IsUnlocked() {
		t.Error("Vault should be unlocked")
	}

	// Test Set/Get
	key := "api_key"
	value := "secret_value"

	if err := v.Set(key, value); err != nil {
		t.Fatalf("Failed to set secret: %v", err)
	}

	got, err := v.Get(key)
	if err != nil {
		t.Fatalf("Failed to get secret: %v", err)
	}

	if got != value {
		t.Errorf("Expected %s, got %s", value, got)
	}

	// Test List
	keys, err := v.List()
	if err != nil {
		t.Fatalf("Failed to list secrets: %v", err)
	}

	if len(keys) != 1 || keys[0] != key {
		t.Errorf("Expected [api_key], got %v", keys)
	}

	// Test Delete
	if err := v.Delete(key); err != nil {
		t.Fatalf("Failed to delete secret: %v", err)
	}

	_, err = v.Get(key)
	if err == nil {
		t.Error("Expected error when getting deleted secret")
	}

	// Test Lock
	v.Lock()
	if v.IsUnlocked() {
		t.Error("Vault should be locked")
	}

	// Test persistence
	v2 := vault.New(vaultPath)
	if !v2.Exists() {
		t.Error("Existing vault should be detected")
	}

	if err := v2.Unlock(password); err != nil {
		t.Fatalf("Failed to unlock existing vault: %v", err)
	}

	// Verify secret is gone (from Delete above)
	keys, err = v2.List()
	if err != nil {
		t.Fatalf("Failed to list secrets from reloaded vault: %v", err)
	}
	if len(keys) != 0 {
		t.Errorf("Expected empty vault, got %v", keys)
	}
}
