package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/crypto/argon2"
)

// Vault stores encrypted secrets
type Vault struct {
	FilePath string
	key      []byte
	secrets  map[string]string
	unlocked bool
	mu       sync.RWMutex
}

// VaultData is the encrypted file format
type VaultData struct {
	Salt  []byte `json:"salt"`
	Nonce []byte `json:"nonce"`
	Data  []byte `json:"data"`
}

// Argon2 parameters
const (
	argonTime    = 1
	argonMemory  = 64 * 1024
	argonThreads = 4
	argonKeyLen  = 32
	saltLen      = 32
)

// New creates a new vault instance
func New(filePath string) *Vault {
	return &Vault{
		FilePath: filePath,
		secrets:  make(map[string]string),
	}
}

// IsUnlocked returns whether the vault is unlocked
func (v *Vault) IsUnlocked() bool {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.unlocked
}

// Exists checks if the vault file exists
func (v *Vault) Exists() bool {
	_, err := os.Stat(v.FilePath)
	return err == nil
}

// Create creates a new vault with the given password
func (v *Vault) Create(password string) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	// Generate salt
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive key from password
	v.key = argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)
	v.secrets = make(map[string]string)
	v.unlocked = true

	return v.saveWithSalt(salt)
}

// Unlock opens an existing vault with the given password
func (v *Vault) Unlock(password string) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	// Read vault file
	data, err := os.ReadFile(v.FilePath)
	if err != nil {
		return fmt.Errorf("failed to read vault: %w", err)
	}

	var vaultData VaultData
	if err := json.Unmarshal(data, &vaultData); err != nil {
		return fmt.Errorf("invalid vault format: %w", err)
	}

	// Derive key from password
	v.key = argon2.IDKey([]byte(password), vaultData.Salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	// Decrypt
	block, err := aes.NewCipher(v.key)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %w", err)
	}

	plaintext, err := gcm.Open(nil, vaultData.Nonce, vaultData.Data, nil)
	if err != nil {
		return fmt.Errorf("wrong password or corrupted vault")
	}

	// Parse secrets
	if err := json.Unmarshal(plaintext, &v.secrets); err != nil {
		return fmt.Errorf("failed to parse secrets: %w", err)
	}

	v.unlocked = true
	return nil
}

// Lock locks the vault
func (v *Vault) Lock() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.key = nil
	v.secrets = make(map[string]string)
	v.unlocked = false
}

// Set stores a secret
func (v *Vault) Set(key, value string) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	if !v.unlocked {
		return fmt.Errorf("vault is locked")
	}

	v.secrets[key] = value
	return v.save()
}

// Get retrieves a secret
func (v *Vault) Get(key string) (string, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if !v.unlocked {
		return "", fmt.Errorf("vault is locked")
	}

	value, ok := v.secrets[key]
	if !ok {
		return "", fmt.Errorf("secret not found: %s", key)
	}

	return value, nil
}

// Delete removes a secret
func (v *Vault) Delete(key string) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	if !v.unlocked {
		return fmt.Errorf("vault is locked")
	}

	if _, ok := v.secrets[key]; !ok {
		return fmt.Errorf("secret not found: %s", key)
	}

	delete(v.secrets, key)
	return v.save()
}

// List returns all secret keys
func (v *Vault) List() ([]string, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if !v.unlocked {
		return nil, fmt.Errorf("vault is locked")
	}

	keys := make([]string, 0, len(v.secrets))
	for k := range v.secrets {
		keys = append(keys, k)
	}
	return keys, nil
}

// Count returns the number of secrets
func (v *Vault) Count() int {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return len(v.secrets)
}

// save encrypts and writes the vault to disk
func (v *Vault) save() error {
	// Read existing salt
	data, err := os.ReadFile(v.FilePath)
	if err != nil {
		return err
	}

	var vaultData VaultData
	if err := json.Unmarshal(data, &vaultData); err != nil {
		return err
	}

	return v.saveWithSalt(vaultData.Salt)
}

// saveWithSalt encrypts and writes the vault with a specific salt
func (v *Vault) saveWithSalt(salt []byte) error {
	// Serialize secrets
	plaintext, err := json.Marshal(v.secrets)
	if err != nil {
		return fmt.Errorf("failed to serialize secrets: %w", err)
	}

	// Create cipher
	block, err := aes.NewCipher(v.key)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	// Create vault data
	vaultData := VaultData{
		Salt:  salt,
		Nonce: nonce,
		Data:  ciphertext,
	}

	// Serialize
	data, err := json.Marshal(vaultData)
	if err != nil {
		return fmt.Errorf("failed to serialize vault: %w", err)
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(v.FilePath), 0700); err != nil {
		return err
	}

	// Write with restricted permissions
	return os.WriteFile(v.FilePath, data, 0600)
}

// Export returns all secrets (for backup)
func (v *Vault) Export() (map[string]string, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if !v.unlocked {
		return nil, fmt.Errorf("vault is locked")
	}

	// Return a copy
	export := make(map[string]string, len(v.secrets))
	for k, val := range v.secrets {
		export[k] = val
	}
	return export, nil
}

// Import adds secrets from a map
func (v *Vault) Import(secrets map[string]string) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	if !v.unlocked {
		return fmt.Errorf("vault is locked")
	}

	for k, val := range secrets {
		v.secrets[k] = val
	}

	return v.save()
}
