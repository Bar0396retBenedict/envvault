package vault

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// NoteEntry holds a human-readable annotation attached to a vault key.
type NoteEntry struct {
	Key       string    `json:"key"`
	Note      string    `json:"note"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NoteRecord is the top-level structure persisted to disk.
type NoteRecord struct {
	Notes map[string]NoteEntry `json:"notes"`
}

func noteFilePath(vaultPath string) string {
	dir := filepath.Dir(vaultPath)
	base := filepath.Base(vaultPath)
	return filepath.Join(dir, "."+base+".notes.json")
}

// LoadNoteRecord reads the note record for the given vault file.
// If the file does not exist, an empty record is returned.
func LoadNoteRecord(vaultPath string) (NoteRecord, error) {
	path := noteFilePath(vaultPath)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return NoteRecord{Notes: make(map[string]NoteEntry)}, nil
	}
	if err != nil {
		return NoteRecord{}, err
	}
	var rec NoteRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return NoteRecord{}, err
	}
	if rec.Notes == nil {
		rec.Notes = make(map[string]NoteEntry)
	}
	return rec, nil
}

func saveNoteRecord(vaultPath string, rec NoteRecord) error {
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(noteFilePath(vaultPath), data, 0o600)
}

// SetNote attaches or replaces a note for the given key.
func SetNote(vaultPath, key, note string) error {
	rec, err := LoadNoteRecord(vaultPath)
	if err != nil {
		return err
	}
	rec.Notes[key] = NoteEntry{Key: key, Note: note, UpdatedAt: time.Now().UTC()}
	return saveNoteRecord(vaultPath, rec)
}

// RemoveNote deletes the note for the given key.
// Returns nil even if no note existed.
func RemoveNote(vaultPath, key string) error {
	rec, err := LoadNoteRecord(vaultPath)
	if err != nil {
		return err
	}
	delete(rec.Notes, key)
	return saveNoteRecord(vaultPath, rec)
}

// GetNote returns the NoteEntry for a key and a boolean indicating existence.
func GetNote(vaultPath, key string) (NoteEntry, bool) {
	rec, err := LoadNoteRecord(vaultPath)
	if err != nil {
		return NoteEntry{}, false
	}
	e, ok := rec.Notes[key]
	return e, ok
}
