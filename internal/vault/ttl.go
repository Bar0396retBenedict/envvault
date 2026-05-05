package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// TTLEntry records an expiry time for a single vault key.
type TTLEntry struct {
	Key       string    `json:"key"`
	ExpiresAt time.Time `json:"expires_at"`
}

// TTLRecord holds all TTL entries for a vault file.
type TTLRecord struct {
	Entries []TTLEntry `json:"entries"`
}

func ttlFilePath(vaultPath string) string {
	dir := filepath.Dir(vaultPath)
	base := filepath.Base(vaultPath)
	return filepath.Join(dir, "."+base+".ttl.json")
}

// LoadTTLRecord loads the TTL record associated with vaultPath.
// Returns an empty record if the file does not exist.
func LoadTTLRecord(vaultPath string) (TTLRecord, error) {
	path := ttlFilePath(vaultPath)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return TTLRecord{}, nil
	}
	if err != nil {
		return TTLRecord{}, fmt.Errorf("ttl: read %s: %w", path, err)
	}
	var rec TTLRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return TTLRecord{}, fmt.Errorf("ttl: unmarshal: %w", err)
	}
	return rec, nil
}

// SaveTTLRecord persists the TTL record next to the vault file.
func SaveTTLRecord(vaultPath string, rec TTLRecord) error {
	path := ttlFilePath(vaultPath)
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return fmt.Errorf("ttl: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("ttl: write %s: %w", path, err)
	}
	return nil
}

// SetTTL sets or replaces the expiry duration for key in the vault at vaultPath.
func SetTTL(vaultPath, key string, ttl time.Duration) error {
	rec, err := LoadTTLRecord(vaultPath)
	if err != nil {
		return err
	}
	expiresAt := time.Now().UTC().Add(ttl)
	for i, e := range rec.Entries {
		if e.Key == key {
			rec.Entries[i].ExpiresAt = expiresAt
			return SaveTTLRecord(vaultPath, rec)
		}
	}
	rec.Entries = append(rec.Entries, TTLEntry{Key: key, ExpiresAt: expiresAt})
	return SaveTTLRecord(vaultPath, rec)
}

// ExpiredKeys returns all keys whose TTL has elapsed as of now.
func ExpiredKeys(vaultPath string) ([]string, error) {
	rec, err := LoadTTLRecord(vaultPath)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	var expired []string
	for _, e := range rec.Entries {
		if now.After(e.ExpiresAt) {
			expired = append(expired, e.Key)
		}
	}
	return expired, nil
}

// PurgeTTLEntry removes the TTL record for key, if present.
func PurgeTTLEntry(vaultPath, key string) error {
	rec, err := LoadTTLRecord(vaultPath)
	if err != nil {
		return err
	}
	filtered := rec.Entries[:0]
	for _, e := range rec.Entries {
		if e.Key != key {
			filtered = append(filtered, e)
		}
	}
	rec.Entries = filtered
	return SaveTTLRecord(vaultPath, rec)
}
