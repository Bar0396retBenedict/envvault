package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot represents a point-in-time capture of a vault's contents.
type Snapshot struct {
	CreatedAt time.Time         `json:"created_at"`
	Label     string            `json:"label"`
	Entries   map[string]string `json:"entries"`
}

// snapshotDir returns the directory where snapshots for a vault file are stored.
func snapshotDir(vaultPath string) string {
	base := filepath.Base(vaultPath)
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]
	return filepath.Join(filepath.Dir(vaultPath), ".snapshots", name)
}

// TakeSnapshot decrypts the vault and saves a labeled snapshot to disk.
func TakeSnapshot(vaultPath, passphrase, label string) (*Snapshot, error) {
	v, err := Load(vaultPath, passphrase)
	if err != nil {
		return nil, fmt.Errorf("snapshot: load vault: %w", err)
	}

	snap := &Snapshot{
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Entries:   make(map[string]string),
	}
	for k, val := range v.data {
		snap.Entries[k] = val
	}

	dir := snapshotDir(vaultPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("snapshot: mkdir: %w", err)
	}

	filename := fmt.Sprintf("%d_%s.json", snap.CreatedAt.UnixNano(), sanitizeLabel(label))
	path := filepath.Join(dir, filename)

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("snapshot: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return nil, fmt.Errorf("snapshot: write: %w", err)
	}
	return snap, nil
}

// ListSnapshots returns all snapshots for the given vault file, ordered by filename.
func ListSnapshots(vaultPath string) ([]Snapshot, error) {
	dir := snapshotDir(vaultPath)
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("snapshot: read dir: %w", err)
	}

	var snaps []Snapshot
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, fmt.Errorf("snapshot: read file %s: %w", e.Name(), err)
		}
		var s Snapshot
		if err := json.Unmarshal(data, &s); err != nil {
			return nil, fmt.Errorf("snapshot: parse %s: %w", e.Name(), err)
		}
		snaps = append(snaps, s)
	}
	return snaps, nil
}

// sanitizeLabel replaces characters unsafe for filenames with underscores.
func sanitizeLabel(label string) string {
	out := make([]byte, len(label))
	for i := 0; i < len(label); i++ {
		c := label[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' {
			out[i] = c
		} else {
			out[i] = '_'
		}
	}
	return string(out)
}
