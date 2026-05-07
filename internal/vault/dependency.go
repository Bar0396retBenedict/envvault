package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// DependencyRecord maps each key to a list of keys it depends on.
type DependencyRecord struct {
	Deps map[string][]string `json:"deps"`
}

func dependencyFilePath(vaultPath string) string {
	dir := filepath.Dir(vaultPath)
	base := filepath.Base(vaultPath)
	return filepath.Join(dir, "."+base+".deps.json")
}

// LoadDependencyRecord loads the dependency record for the given vault file.
// If the file does not exist, an empty record is returned.
func LoadDependencyRecord(vaultPath string) (DependencyRecord, error) {
	rec := DependencyRecord{Deps: make(map[string][]string)}
	data, err := os.ReadFile(dependencyFilePath(vaultPath))
	if os.IsNotExist(err) {
		return rec, nil
	}
	if err != nil {
		return rec, fmt.Errorf("dependency: read: %w", err)
	}
	if err := json.Unmarshal(data, &rec); err != nil {
		return rec, fmt.Errorf("dependency: parse: %w", err)
	}
	return rec, nil
}

func saveDependencyRecord(vaultPath string, rec DependencyRecord) error {
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return fmt.Errorf("dependency: marshal: %w", err)
	}
	return os.WriteFile(dependencyFilePath(vaultPath), data, 0o600)
}

// AddDependency records that key depends on depKey.
func AddDependency(vaultPath, key, depKey string) error {
	rec, err := LoadDependencyRecord(vaultPath)
	if err != nil {
		return err
	}
	for _, existing := range rec.Deps[key] {
		if existing == depKey {
			return nil
		}
	}
	rec.Deps[key] = append(rec.Deps[key], depKey)
	sort.Strings(rec.Deps[key])
	return saveDependencyRecord(vaultPath, rec)
}

// RemoveDependency removes depKey from the dependency list of key.
func RemoveDependency(vaultPath, key, depKey string) error {
	rec, err := LoadDependencyRecord(vaultPath)
	if err != nil {
		return err
	}
	deps := rec.Deps[key]
	updated := deps[:0]
	for _, d := range deps {
		if d != depKey {
			updated = append(updated, d)
		}
	}
	if len(updated) == 0 {
		delete(rec.Deps, key)
	} else {
		rec.Deps[key] = updated
	}
	return saveDependencyRecord(vaultPath, rec)
}

// ListDependencies returns the keys that key depends on.
func ListDependencies(vaultPath, key string) ([]string, error) {
	rec, err := LoadDependencyRecord(vaultPath)
	if err != nil {
		return nil, err
	}
	return rec.Deps[key], nil
}

// Dependents returns all keys that depend on the given key.
func Dependents(vaultPath, key string) ([]string, error) {
	rec, err := LoadDependencyRecord(vaultPath)
	if err != nil {
		return nil, err
	}
	var result []string
	for k, deps := range rec.Deps {
		for _, d := range deps {
			if d == key {
				result = append(result, k)
				break
			}
		}
	}
	sort.Strings(result)
	return result, nil
}
