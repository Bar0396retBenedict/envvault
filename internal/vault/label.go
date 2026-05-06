package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// LabelRecord maps vault keys to human-readable labels (display names).
type LabelRecord struct {
	Labels map[string]string `json:"labels"`
}

func labelFilePath(vaultPath string) string {
	dir := filepath.Dir(vaultPath)
	base := filepath.Base(vaultPath)
	return filepath.Join(dir, "."+base+".labels.json")
}

// LoadLabelRecord reads the label record for the given vault file.
// If the file does not exist, an empty record is returned.
func LoadLabelRecord(vaultPath string) (LabelRecord, error) {
	path := labelFilePath(vaultPath)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return LabelRecord{Labels: make(map[string]string)}, nil
	}
	if err != nil {
		return LabelRecord{}, fmt.Errorf("read label record: %w", err)
	}
	var rec LabelRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return LabelRecord{}, fmt.Errorf("parse label record: %w", err)
	}
	if rec.Labels == nil {
		rec.Labels = make(map[string]string)
	}
	return rec, nil
}

func saveLabelRecord(vaultPath string, rec LabelRecord) error {
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal label record: %w", err)
	}
	return os.WriteFile(labelFilePath(vaultPath), data, 0600)
}

// SetLabel assigns a human-readable label to a vault key.
func SetLabel(vaultPath, key, label string) error {
	if key == "" {
		return fmt.Errorf("key must not be empty")
	}
	rec, err := LoadLabelRecord(vaultPath)
	if err != nil {
		return err
	}
	rec.Labels[key] = label
	return saveLabelRecord(vaultPath, rec)
}

// RemoveLabel removes the label for the given key.
func RemoveLabel(vaultPath, key string) error {
	rec, err := LoadLabelRecord(vaultPath)
	if err != nil {
		return err
	}
	if _, ok := rec.Labels[key]; !ok {
		return fmt.Errorf("no label set for key %q", key)
	}
	delete(rec.Labels, key)
	return saveLabelRecord(vaultPath, rec)
}

// GetLabel returns the label for the given key, or an empty string if none is set.
func GetLabel(vaultPath, key string) (string, error) {
	rec, err := LoadLabelRecord(vaultPath)
	if err != nil {
		return "", err
	}
	return rec.Labels[key], nil
}

// ListLabels returns all key→label pairs sorted by key.
func ListLabels(vaultPath string) ([]string, []string, error) {
	rec, err := LoadLabelRecord(vaultPath)
	if err != nil {
		return nil, nil, err
	}
	keys := make([]string, 0, len(rec.Labels))
	for k := range rec.Labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	labels := make([]string, len(keys))
	for i, k := range keys {
		labels[i] = rec.Labels[k]
	}
	return keys, labels, nil
}
