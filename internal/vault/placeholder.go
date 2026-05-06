package vault

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// placeholderFilePath returns the path to the placeholder record file
// adjacent to the given vault file.
func placeholderFilePath(vaultPath string) string {
	dir := filepath.Dir(vaultPath)
	base := strings.TrimSuffix(filepath.Base(vaultPath), filepath.Ext(vaultPath))
	return filepath.Join(dir, "."+base+".placeholders.json")
}

// PlaceholderRecord maps vault keys to their placeholder descriptions.
type PlaceholderRecord struct {
	Placeholders map[string]string `json:"placeholders"`
}

// LoadPlaceholderRecord reads the placeholder record for the given vault.
// If the file does not exist, an empty record is returned.
func LoadPlaceholderRecord(vaultPath string) (PlaceholderRecord, error) {
	p := placeholderFilePath(vaultPath)
	var rec PlaceholderRecord
	if err := loadJSON(p, &rec); err != nil {
		if os.IsNotExist(err) {
			return PlaceholderRecord{Placeholders: make(map[string]string)}, nil
		}
		return rec, err
	}
	if rec.Placeholders == nil {
		rec.Placeholders = make(map[string]string)
	}
	return rec, nil
}

// savePlaceholderRecord persists the placeholder record to disk.
func savePlaceholderRecord(vaultPath string, rec PlaceholderRecord) error {
	p := placeholderFilePath(vaultPath)
	return saveJSON(p, rec, 0o600)
}

// SetPlaceholder associates a human-readable placeholder description with a
// vault key. The description is intended to document what value is expected,
// e.g. "AWS access key ID (20 uppercase chars)".
func SetPlaceholder(vaultPath, key, description string) error {
	if key == "" {
		return fmt.Errorf("key must not be empty")
	}
	if description == "" {
		return fmt.Errorf("description must not be empty")
	}
	rec, err := LoadPlaceholderRecord(vaultPath)
	if err != nil {
		return err
	}
	rec.Placeholders[key] = description
	return savePlaceholderRecord(vaultPath, rec)
}

// RemovePlaceholder deletes the placeholder entry for the given key.
func RemovePlaceholder(vaultPath, key string) error {
	rec, err := LoadPlaceholderRecord(vaultPath)
	if err != nil {
		return err
	}
	if _, ok := rec.Placeholders[key]; !ok {
		return fmt.Errorf("no placeholder found for key %q", key)
	}
	delete(rec.Placeholders, key)
	return savePlaceholderRecord(vaultPath, rec)
}

// ListPlaceholders returns all placeholder entries sorted by key.
func ListPlaceholders(vaultPath string) ([]PlaceholderEntry, error) {
	rec, err := LoadPlaceholderRecord(vaultPath)
	if err != nil {
		return nil, err
	}
	keys := sortedKeys(rec.Placeholders)
	out := make([]PlaceholderEntry, 0, len(keys))
	for _, k := range keys {
		out = append(out, PlaceholderEntry{Key: k, Description: rec.Placeholders[k]})
	}
	return out, nil
}

// PlaceholderEntry is a single key/description pair returned by ListPlaceholders.
type PlaceholderEntry struct {
	Key         string
	Description string
}
