package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// GroupRecord maps group names to lists of vault keys.
type GroupRecord struct {
	Groups map[string][]string `json:"groups"`
}

func groupFilePath(vaultPath string) string {
	dir := filepath.Dir(vaultPath)
	base := filepath.Base(vaultPath)
	return filepath.Join(dir, "."+base+".groups.json")
}

// LoadGroupRecord reads the group record for the given vault file.
// If no record exists, an empty record is returned.
func LoadGroupRecord(vaultPath string) (GroupRecord, error) {
	path := groupFilePath(vaultPath)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return GroupRecord{Groups: make(map[string][]string)}, nil
	}
	if err != nil {
		return GroupRecord{}, fmt.Errorf("read group record: %w", err)
	}
	var rec GroupRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return GroupRecord{}, fmt.Errorf("parse group record: %w", err)
	}
	if rec.Groups == nil {
		rec.Groups = make(map[string][]string)
	}
	return rec, nil
}

func saveGroupRecord(vaultPath string, rec GroupRecord) error {
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal group record: %w", err)
	}
	return os.WriteFile(groupFilePath(vaultPath), data, 0600)
}

// AddToGroup adds a key to the named group, creating the group if needed.
// Duplicate entries are ignored.
func AddToGroup(vaultPath, group, key string) error {
	rec, err := LoadGroupRecord(vaultPath)
	if err != nil {
		return err
	}
	for _, k := range rec.Groups[group] {
		if k == key {
			return nil
		}
	}
	rec.Groups[group] = append(rec.Groups[group], key)
	sort.Strings(rec.Groups[group])
	return saveGroupRecord(vaultPath, rec)
}

// RemoveFromGroup removes a key from the named group.
// Removing the last key deletes the group entry.
func RemoveFromGroup(vaultPath, group, key string) error {
	rec, err := LoadGroupRecord(vaultPath)
	if err != nil {
		return err
	}
	keys := rec.Groups[group]
	newKeys := keys[:0]
	for _, k := range keys {
		if k != key {
			newKeys = append(newKeys, k)
		}
	}
	if len(newKeys) == 0 {
		delete(rec.Groups, group)
	} else {
		rec.Groups[group] = newKeys
	}
	return saveGroupRecord(vaultPath, rec)
}

// KeysForGroup returns all keys belonging to the named group.
func KeysForGroup(vaultPath, group string) ([]string, error) {
	rec, err := LoadGroupRecord(vaultPath)
	if err != nil {
		return nil, err
	}
	keys := rec.Groups[group]
	out := make([]string, len(keys))
	copy(out, keys)
	return out, nil
}

// ListGroups returns all group names in sorted order.
func ListGroups(vaultPath string) ([]string, error) {
	rec, err := LoadGroupRecord(vaultPath)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(rec.Groups))
	for g := range rec.Groups {
		names = append(names, g)
	}
	sort.Strings(names)
	return names, nil
}
