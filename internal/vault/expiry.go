package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ExpiryRecord maps key names to their expiry timestamps.
type ExpiryRecord map[string]time.Time

func expiryFilePath(vaultPath string) string {
	dir := filepath.Dir(vaultPath)
	base := filepath.Base(vaultPath)
	return filepath.Join(dir, "."+base+".expiry.json")
}

// LoadExpiryRecord reads the expiry metadata for a vault file.
// If no expiry file exists, an empty record is returned.
func LoadExpiryRecord(vaultPath string) (ExpiryRecord, error) {
	path := expiryFilePath(vaultPath)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return ExpiryRecord{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read expiry file: %w", err)
	}
	var rec ExpiryRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return nil, fmt.Errorf("parse expiry file: %w", err)
	}
	return rec, nil
}

// SaveExpiryRecord persists the expiry record to disk.
func SaveExpiryRecord(vaultPath string, rec ExpiryRecord) error {
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal expiry record: %w", err)
	}
	path := expiryFilePath(vaultPath)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("write expiry file: %w", err)
	}
	return nil
}

// SetExpiry assigns an expiry time to a key in the given vault.
// The key must exist in the vault.
func SetExpiry(vaultPath, passphrase, key string, expiresAt time.Time) error {
	v, err := Load(vaultPath, passphrase)
	if err != nil {
		return fmt.Errorf("load vault: %w", err)
	}
	if _, ok := v.Get(key); !ok {
		return fmt.Errorf("key %q not found in vault", key)
	}
	rec, err := LoadExpiryRecord(vaultPath)
	if err != nil {
		return err
	}
	rec[key] = expiresAt
	return SaveExpiryRecord(vaultPath, rec)
}

// ExpiredKeysList returns all keys whose expiry time is before now.
func ExpiredKeysList(vaultPath string) ([]string, error) {
	rec, err := LoadExpiryRecord(vaultPath)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	var expired []string
	for key, exp := range rec {
		if now.After(exp) {
			expired = append(expired, key)
		}
	}
	return expired, nil
}
