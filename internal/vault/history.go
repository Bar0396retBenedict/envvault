package vault

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// HistoryEntry records a single change to a vault key.
type HistoryEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Key       string    `json:"key"`
	OldValue  string    `json:"old_value"`
	NewValue  string    `json:"new_value"`
	Action    string    `json:"action"` // "set", "delete", "rename"
}

// HistoryRecord holds all history entries for a vault.
type HistoryRecord struct {
	Entries []HistoryEntry `json:"entries"`
}

func historyFilePath(vaultPath string) string {
	dir := filepath.Dir(vaultPath)
	base := filepath.Base(vaultPath)
	return filepath.Join(dir, "."+base+".history.json")
}

// LoadHistory loads the history record for the given vault file.
// Returns an empty record if the file does not exist.
func LoadHistory(vaultPath string) (HistoryRecord, error) {
	p := historyFilePath(vaultPath)
	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return HistoryRecord{}, nil
	}
	if err != nil {
		return HistoryRecord{}, err
	}
	var rec HistoryRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return HistoryRecord{}, err
	}
	return rec, nil
}

// AppendHistory appends a new entry to the history record and saves it.
func AppendHistory(vaultPath string, entry HistoryEntry) error {
	rec, err := LoadHistory(vaultPath)
	if err != nil {
		return err
	}
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}
	rec.Entries = append(rec.Entries, entry)
	return saveHistory(vaultPath, rec)
}

func saveHistory(vaultPath string, rec HistoryRecord) error {
	p := historyFilePath(vaultPath)
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0600)
}
