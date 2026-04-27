// Package vault provides functionality for reading, writing, and managing
// encrypted environment variable files (vaults).
package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/envvault/internal/crypto"
)

// Vault represents an encrypted collection of environment variables.
type Vault struct {
	Version   int               `json:"version"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Env       map[string]string `json:"env"`
}

// New creates a new empty Vault.
func New() *Vault {
	now := time.Now().UTC()
	return &Vault{
		Version:   1,
		CreatedAt: now,
		UpdatedAt: now,
		Env:       make(map[string]string),
	}
}

// Set adds or updates an environment variable in the vault.
func (v *Vault) Set(key, value string) {
	v.Env[key] = value
	v.UpdatedAt = time.Now().UTC()
}

// Get retrieves an environment variable from the vault.
func (v *Vault) Get(key string) (string, bool) {
	val, ok := v.Env[key]
	return val, ok
}

// Delete removes an environment variable from the vault.
func (v *Vault) Delete(key string) {
	delete(v.Env, key)
	v.UpdatedAt = time.Now().UTC()
}

// Save encrypts the vault and writes it to the given file path.
func (v *Vault) Save(path, passphrase string) error {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("vault: marshal: %w", err)
	}

	key, err := crypto.DeriveKey(passphrase, nil)
	if err != nil {
		return fmt.Errorf("vault: derive key: %w", err)
	}

	ciphertext, err := crypto.Encrypt(key, data)
	if err != nil {
		return fmt.Errorf("vault: encrypt: %w", err)
	}

	if err := os.WriteFile(path, ciphertext, 0600); err != nil {
		return fmt.Errorf("vault: write file: %w", err)
	}
	return nil
}

// Load reads an encrypted vault file and decrypts it.
func Load(path, passphrase string) (*Vault, error) {
	ciphertext, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("vault: read file: %w", err)
	}

	key, err := crypto.DeriveKey(passphrase, nil)
	if err != nil {
		return nil, fmt.Errorf("vault: derive key: %w", err)
	}

	data, err := crypto.Decrypt(key, ciphertext)
	if err != nil {
		return nil, fmt.Errorf("vault: decrypt: %w", err)
	}

	var v Vault
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("vault: unmarshal: %w", err)
	}
	return &v, nil
}
