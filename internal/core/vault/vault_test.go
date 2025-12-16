package vault

import (
	"os"
	"path/filepath"
	"testing"
)

// ==============================================================
// New Tests
// ==============================================================

func TestNew(t *testing.T) {
	v := New("/test/vault.enc")

	if v == nil {
		t.Fatal("New() returned nil")
	}
	if v.FilePath != "/test/vault.enc" {
		t.Errorf("FilePath = %q, want '/test/vault.enc'", v.FilePath)
	}
	if v.secrets == nil {
		t.Error("secrets should be initialized")
	}
	if v.unlocked {
		t.Error("New vault should be locked")
	}
}

// ==============================================================
// IsUnlocked Tests
// ==============================================================

func TestVault_IsUnlocked(t *testing.T) {
	v := New("/test/vault.enc")

	if v.IsUnlocked() {
		t.Error("New vault should not be unlocked")
	}
}

// ==============================================================
// Exists Tests
// ==============================================================

func TestVault_Exists(t *testing.T) {
	t.Run("file_does_not_exist", func(t *testing.T) {
		v := New("/nonexistent/vault.enc")
		if v.Exists() {
			t.Error("Exists() should return false for nonexistent file")
		}
	})

	t.Run("file_exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		vaultPath := filepath.Join(tmpDir, "vault.enc")

		// Create empty file
		if err := os.WriteFile(vaultPath, []byte("{}"), 0o600); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		v := New(vaultPath)
		if !v.Exists() {
			t.Error("Exists() should return true for existing file")
		}
	})
}

// ==============================================================
// Create Tests
// ==============================================================

func TestVault_Create(t *testing.T) {
	tmpDir := t.TempDir()
	vaultPath := filepath.Join(tmpDir, "vault.enc")

	v := New(vaultPath)
	err := v.Create("test-password")

	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Check vault is unlocked
	if !v.IsUnlocked() {
		t.Error("Vault should be unlocked after Create()")
	}

	// Check file was created
	if !v.Exists() {
		t.Error("Vault file should exist after Create()")
	}

	// Check key was derived
	if v.key == nil {
		t.Error("Key should be derived after Create()")
	}
}

// ==============================================================
// Unlock Tests
// ==============================================================

func TestVault_Unlock(t *testing.T) {
	tmpDir := t.TempDir()
	vaultPath := filepath.Join(tmpDir, "vault.enc")

	// Create vault first
	v := New(vaultPath)
	if err := v.Create("test-password"); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Lock and unlock
	v.Lock()
	if v.IsUnlocked() {
		t.Error("Vault should be locked after Lock()")
	}

	err := v.Unlock("test-password")
	if err != nil {
		t.Fatalf("Unlock() error = %v", err)
	}

	if !v.IsUnlocked() {
		t.Error("Vault should be unlocked after Unlock()")
	}
}

func TestVault_Unlock_WrongPassword(t *testing.T) {
	tmpDir := t.TempDir()
	vaultPath := filepath.Join(tmpDir, "vault.enc")

	v := New(vaultPath)
	if err := v.Create("correct-password"); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	v.Lock()

	err := v.Unlock("wrong-password")
	if err == nil {
		t.Error("Unlock() should fail with wrong password")
	}
}

func TestVault_Unlock_NonexistentFile(t *testing.T) {
	v := New("/nonexistent/vault.enc")
	err := v.Unlock("test-password")
	if err == nil {
		t.Error("Unlock() should fail for nonexistent file")
	}
}

// ==============================================================
// Lock Tests
// ==============================================================

func TestVault_Lock(t *testing.T) {
	tmpDir := t.TempDir()
	vaultPath := filepath.Join(tmpDir, "vault.enc")

	v := New(vaultPath)
	if err := v.Create("test-password"); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	v.Lock()

	if v.IsUnlocked() {
		t.Error("Vault should be locked after Lock()")
	}
	if v.key != nil {
		t.Error("Key should be cleared after Lock()")
	}
	if len(v.secrets) != 0 {
		t.Error("Secrets should be cleared after Lock()")
	}
}

// ==============================================================
// Set/Get/Delete Tests
// ==============================================================

func TestVault_SetGet(t *testing.T) {
	tmpDir := t.TempDir()
	vaultPath := filepath.Join(tmpDir, "vault.enc")

	v := New(vaultPath)
	if err := v.Create("test-password"); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Set a secret
	if err := v.Set("API_KEY", "secret-value-123"); err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	// Get the secret
	value, err := v.Get("API_KEY")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if value != "secret-value-123" {
		t.Errorf("Get() = %q, want 'secret-value-123'", value)
	}
}

func TestVault_Get_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	vaultPath := filepath.Join(tmpDir, "vault.enc")

	v := New(vaultPath)
	if err := v.Create("test-password"); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	_, err := v.Get("NONEXISTENT_KEY")
	if err == nil {
		t.Error("Get() should fail for nonexistent key")
	}
}

func TestVault_Set_WhenLocked(t *testing.T) {
	v := New("/test/vault.enc")

	err := v.Set("key", "value")
	if err == nil {
		t.Error("Set() should fail when vault is locked")
	}
}

func TestVault_Get_WhenLocked(t *testing.T) {
	v := New("/test/vault.enc")

	_, err := v.Get("key")
	if err == nil {
		t.Error("Get() should fail when vault is locked")
	}
}

func TestVault_Delete(t *testing.T) {
	tmpDir := t.TempDir()
	vaultPath := filepath.Join(tmpDir, "vault.enc")

	v := New(vaultPath)
	if err := v.Create("test-password"); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if err := v.Set("KEY_TO_DELETE", "value"); err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	if err := v.Delete("KEY_TO_DELETE"); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err := v.Get("KEY_TO_DELETE")
	if err == nil {
		t.Error("Get() should fail after Delete()")
	}
}

func TestVault_Delete_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	vaultPath := filepath.Join(tmpDir, "vault.enc")

	v := New(vaultPath)
	if err := v.Create("test-password"); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	err := v.Delete("NONEXISTENT")
	if err == nil {
		t.Error("Delete() should fail for nonexistent key")
	}
}

func TestVault_Delete_WhenLocked(t *testing.T) {
	v := New("/test/vault.enc")

	err := v.Delete("key")
	if err == nil {
		t.Error("Delete() should fail when vault is locked")
	}
}

// ==============================================================
// List/Count Tests
// ==============================================================

func TestVault_List(t *testing.T) {
	tmpDir := t.TempDir()
	vaultPath := filepath.Join(tmpDir, "vault.enc")

	v := New(vaultPath)
	if err := v.Create("test-password"); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	v.Set("KEY1", "value1")
	v.Set("KEY2", "value2")
	v.Set("KEY3", "value3")

	keys, err := v.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(keys) != 3 {
		t.Errorf("List() returned %d keys, want 3", len(keys))
	}
}

func TestVault_List_WhenLocked(t *testing.T) {
	v := New("/test/vault.enc")

	_, err := v.List()
	if err == nil {
		t.Error("List() should fail when vault is locked")
	}
}

func TestVault_Count(t *testing.T) {
	tmpDir := t.TempDir()
	vaultPath := filepath.Join(tmpDir, "vault.enc")

	v := New(vaultPath)
	if err := v.Create("test-password"); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if v.Count() != 0 {
		t.Error("Count() should be 0 for empty vault")
	}

	v.Set("KEY1", "value1")
	v.Set("KEY2", "value2")

	if v.Count() != 2 {
		t.Errorf("Count() = %d, want 2", v.Count())
	}
}

// ==============================================================
// Export/Import Tests
// ==============================================================

func TestVault_Export(t *testing.T) {
	tmpDir := t.TempDir()
	vaultPath := filepath.Join(tmpDir, "vault.enc")

	v := New(vaultPath)
	if err := v.Create("test-password"); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	v.Set("KEY1", "value1")
	v.Set("KEY2", "value2")

	exported, err := v.Export()
	if err != nil {
		t.Fatalf("Export() error = %v", err)
	}

	if len(exported) != 2 {
		t.Errorf("Export() returned %d secrets, want 2", len(exported))
	}

	if exported["KEY1"] != "value1" {
		t.Errorf("Export()['KEY1'] = %q, want 'value1'", exported["KEY1"])
	}
}

func TestVault_Export_WhenLocked(t *testing.T) {
	v := New("/test/vault.enc")

	_, err := v.Export()
	if err == nil {
		t.Error("Export() should fail when vault is locked")
	}
}

func TestVault_Import(t *testing.T) {
	tmpDir := t.TempDir()
	vaultPath := filepath.Join(tmpDir, "vault.enc")

	v := New(vaultPath)
	if err := v.Create("test-password"); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	secrets := map[string]string{
		"IMPORTED_KEY1": "imported_value1",
		"IMPORTED_KEY2": "imported_value2",
	}

	if err := v.Import(secrets); err != nil {
		t.Fatalf("Import() error = %v", err)
	}

	if v.Count() != 2 {
		t.Errorf("Count() = %d after import, want 2", v.Count())
	}

	val, _ := v.Get("IMPORTED_KEY1")
	if val != "imported_value1" {
		t.Errorf("Get('IMPORTED_KEY1') = %q, want 'imported_value1'", val)
	}
}

func TestVault_Import_WhenLocked(t *testing.T) {
	v := New("/test/vault.enc")

	err := v.Import(map[string]string{"key": "value"})
	if err == nil {
		t.Error("Import() should fail when vault is locked")
	}
}

// ==============================================================
// Persistence Tests
// ==============================================================

func TestVault_Persistence(t *testing.T) {
	tmpDir := t.TempDir()
	vaultPath := filepath.Join(tmpDir, "vault.enc")

	// Create vault and add secrets
	v1 := New(vaultPath)
	if err := v1.Create("test-password"); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	v1.Set("PERSISTENT_KEY", "persistent_value")
	v1.Lock()

	// Open new vault instance and verify
	v2 := New(vaultPath)
	if err := v2.Unlock("test-password"); err != nil {
		t.Fatalf("Unlock() error = %v", err)
	}

	val, err := v2.Get("PERSISTENT_KEY")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if val != "persistent_value" {
		t.Errorf("Persisted value = %q, want 'persistent_value'", val)
	}
}
