package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// SyncMeta holds metadata about the last sync operation.
type SyncMeta struct {
	Environment string    `json:"environment"`
	SyncedAt    time.Time `json:"synced_at"`
	Checksum    string    `json:"checksum"`
}

// SyncRecord maps environment names to their sync metadata.
type SyncRecord map[string]SyncMeta

// LoadSyncRecord reads the sync metadata file from the given path.
// Returns an empty record if the file does not exist.
func LoadSyncRecord(path string) (SyncRecord, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return make(SyncRecord), nil
	}
	if err != nil {
		return nil, fmt.Errorf("read sync record: %w", err)
	}
	var rec SyncRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return nil, fmt.Errorf("parse sync record: %w", err)
	}
	return rec, nil
}

// SaveSyncRecord writes the sync metadata to the given path.
func SaveSyncRecord(path string, rec SyncRecord) error {
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal sync record: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write sync record: %w", err)
	}
	return nil
}

// MergeVaults merges keys from src into dst, returning the number of keys added
// or updated. Keys present in dst but absent in src are left untouched.
func MergeVaults(dst, src *Vault) int {
	count := 0
	for k, v := range src.data {
		existing, ok := dst.data[k]
		if !ok || existing != v {
			dst.data[k] = v
			count++
		}
	}
	return count
}
