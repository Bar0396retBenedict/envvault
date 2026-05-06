package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// AliasRecord maps short alias names to vault key names.
type AliasRecord struct {
	Aliases map[string]string `json:"aliases"`
}

func aliasFilePath(vaultPath string) string {
	dir := filepath.Dir(vaultPath)
	base := filepath.Base(vaultPath)
	return filepath.Join(dir, "."+base+".aliases.json")
}

// LoadAliasRecord loads the alias record for the given vault file.
// If no record exists, an empty record is returned.
func LoadAliasRecord(vaultPath string) (AliasRecord, error) {
	path := aliasFilePath(vaultPath)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return AliasRecord{Aliases: make(map[string]string)}, nil
	}
	if err != nil {
		return AliasRecord{}, fmt.Errorf("alias: read: %w", err)
	}
	var rec AliasRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return AliasRecord{}, fmt.Errorf("alias: parse: %w", err)
	}
	if rec.Aliases == nil {
		rec.Aliases = make(map[string]string)
	}
	return rec, nil
}

func saveAliasRecord(vaultPath string, rec AliasRecord) error {
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return fmt.Errorf("alias: marshal: %w", err)
	}
	return os.WriteFile(aliasFilePath(vaultPath), data, 0600)
}

// SetAlias creates or updates an alias pointing to a vault key.
func SetAlias(vaultPath, alias, key string) error {
	if alias == "" {
		return fmt.Errorf("alias: alias name must not be empty")
	}
	if key == "" {
		return fmt.Errorf("alias: target key must not be empty")
	}
	rec, err := LoadAliasRecord(vaultPath)
	if err != nil {
		return err
	}
	rec.Aliases[alias] = key
	return saveAliasRecord(vaultPath, rec)
}

// RemoveAlias deletes an alias. Returns an error if the alias does not exist.
func RemoveAlias(vaultPath, alias string) error {
	rec, err := LoadAliasRecord(vaultPath)
	if err != nil {
		return err
	}
	if _, ok := rec.Aliases[alias]; !ok {
		return fmt.Errorf("alias: %q not found", alias)
	}
	delete(rec.Aliases, alias)
	return saveAliasRecord(vaultPath, rec)
}

// ResolveAlias returns the vault key that alias points to, or alias itself
// if no mapping exists (pass-through).
func ResolveAlias(vaultPath, alias string) (string, error) {
	rec, err := LoadAliasRecord(vaultPath)
	if err != nil {
		return "", err
	}
	if key, ok := rec.Aliases[alias]; ok {
		return key, nil
	}
	return alias, nil
}

// ListAliases returns all aliases sorted alphabetically.
func ListAliases(vaultPath string) ([][2]string, error) {
	rec, err := LoadAliasRecord(vaultPath)
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(rec.Aliases))
	for k := range rec.Aliases {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := make([][2]string, 0, len(keys))
	for _, k := range keys {
		out = append(out, [2]string{k, rec.Aliases[k]})
	}
	return out, nil
}
