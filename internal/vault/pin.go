package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// PinRecord holds the set of pinned keys for a vault file.
type PinRecord struct {
	Pins map[string]PinEntry `json:"pins"`
}

// PinEntry records metadata about a single pinned key.
type PinEntry struct {
	Note      string    `json:"note,omitempty"`
	PinnedAt  time.Time `json:"pinned_at"`
}

func pinFilePath(vaultPath string) string {
	dir := filepath.Dir(vaultPath)
	base := filepath.Base(vaultPath)
	return filepath.Join(dir, "."+base+".pins.json")
}

// LoadPinRecord loads the pin record for the given vault file.
// Returns an empty record if the file does not exist.
func LoadPinRecord(vaultPath string) (PinRecord, error) {
	p := pinFilePath(vaultPath)
	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return PinRecord{Pins: make(map[string]PinEntry)}, nil
	}
	if err != nil {
		return PinRecord{}, fmt.Errorf("read pin record: %w", err)
	}
	var rec PinRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return PinRecord{}, fmt.Errorf("parse pin record: %w", err)
	}
	if rec.Pins == nil {
		rec.Pins = make(map[string]PinEntry)
	}
	return rec, nil
}

func savePinRecord(vaultPath string, rec PinRecord) error {
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal pin record: %w", err)
	}
	p := pinFilePath(vaultPath)
	return os.WriteFile(p, data, 0600)
}

// PinKey marks a key as pinned with an optional note.
func PinKey(vaultPath, key, note string) error {
	rec, err := LoadPinRecord(vaultPath)
	if err != nil {
		return err
	}
	rec.Pins[key] = PinEntry{Note: note, PinnedAt: time.Now().UTC()}
	return savePinRecord(vaultPath, rec)
}

// UnpinKey removes a key from the pin record.
func UnpinKey(vaultPath, key string) error {
	rec, err := LoadPinRecord(vaultPath)
	if err != nil {
		return err
	}
	if _, ok := rec.Pins[key]; !ok {
		return fmt.Errorf("key %q is not pinned", key)
	}
	delete(rec.Pins, key)
	return savePinRecord(vaultPath, rec)
}

// ListPinnedKeys returns the pinned keys in sorted order.
func ListPinnedKeys(vaultPath string) ([]string, PinRecord, error) {
	rec, err := LoadPinRecord(vaultPath)
	if err != nil {
		return nil, rec, err
	}
	keys := make([]string, 0, len(rec.Pins))
	for k := range rec.Pins {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys, rec, nil
}
