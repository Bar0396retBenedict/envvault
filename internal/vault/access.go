package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// AccessRecord tracks per-key read/write access counts and timestamps.
type AccessRecord struct {
	Entries map[string]*AccessEntry `json:"entries"`
}

// AccessEntry holds access statistics for a single key.
type AccessEntry struct {
	Key       string    `json:"key"`
	ReadCount int       `json:"read_count"`
	WriteCount int      `json:"write_count"`
	LastRead  time.Time `json:"last_read,omitempty"`
	LastWrite time.Time `json:"last_write,omitempty"`
}

func accessFilePath(vaultPath string) string {
	dir := filepath.Dir(vaultPath)
	base := filepath.Base(vaultPath)
	return filepath.Join(dir, "."+base+".access.json")
}

// LoadAccessRecord loads the access record for the given vault file.
// Returns an empty record if the file does not exist.
func LoadAccessRecord(vaultPath string) (*AccessRecord, error) {
	p := accessFilePath(vaultPath)
	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return &AccessRecord{Entries: make(map[string]*AccessEntry)}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("access: read file: %w", err)
	}
	var rec AccessRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return nil, fmt.Errorf("access: unmarshal: %w", err)
	}
	if rec.Entries == nil {
		rec.Entries = make(map[string]*AccessEntry)
	}
	return &rec, nil
}

// saveAccessRecord persists the access record to disk.
func saveAccessRecord(vaultPath string, rec *AccessRecord) error {
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return fmt.Errorf("access: marshal: %w", err)
	}
	p := accessFilePath(vaultPath)
	if err := os.WriteFile(p, data, 0600); err != nil {
		return fmt.Errorf("access: write file: %w", err)
	}
	return nil
}

// RecordRead increments the read counter for the given key.
func RecordRead(vaultPath, key string) error {
	rec, err := LoadAccessRecord(vaultPath)
	if err != nil {
		return err
	}
	e := ensureEntry(rec, key)
	e.ReadCount++
	e.LastRead = time.Now().UTC()
	return saveAccessRecord(vaultPath, rec)
}

// RecordWrite increments the write counter for the given key.
func RecordWrite(vaultPath, key string) error {
	rec, err := LoadAccessRecord(vaultPath)
	if err != nil {
		return err
	}
	e := ensureEntry(rec, key)
	e.WriteCount++
	e.LastWrite = time.Now().UTC()
	return saveAccessRecord(vaultPath, rec)
}

// ListAccessEntries returns all entries sorted by key name.
func ListAccessEntries(vaultPath string) ([]*AccessEntry, error) {
	rec, err := LoadAccessRecord(vaultPath)
	if err != nil {
		return nil, err
	}
	entries := make([]*AccessEntry, 0, len(rec.Entries))
	for _, e := range rec.Entries {
		entries = append(entries, e)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
	return entries, nil
}

func ensureEntry(rec *AccessRecord, key string) *AccessEntry {
	if e, ok := rec.Entries[key]; ok {
		return e
	}
	e := &AccessEntry{Key: key}
	rec.Entries[key] = e
	return e
}
