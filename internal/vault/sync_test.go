package vault

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadSyncRecordMissing(t *testing.T) {
	rec, err := LoadSyncRecord("/nonexistent/path/.envvault.sync")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(rec) != 0 {
		t.Fatalf("expected empty record, got %d entries", len(rec))
	}
}

func TestSaveAndLoadSyncRecord(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".envvault.sync")

	rec := SyncRecord{
		"staging": {
			Environment: "staging",
			SyncedAt:    time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Checksum:    "abc123",
		},
	}

	if err := SaveSyncRecord(path, rec); err != nil {
		t.Fatalf("SaveSyncRecord: %v", err)
	}

	loaded, err := LoadSyncRecord(path)
	if err != nil {
		t.Fatalf("LoadSyncRecord: %v", err)
	}

	meta, ok := loaded["staging"]
	if !ok {
		t.Fatal("expected 'staging' key in loaded record")
	}
	if meta.Checksum != "abc123" {
		t.Errorf("checksum: got %q, want %q", meta.Checksum, "abc123")
	}
}

func TestSaveAndLoadSyncRecordPermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".envvault.sync")

	if err := SaveSyncRecord(path, make(SyncRecord)); err != nil {
		t.Fatalf("SaveSyncRecord: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("file permissions: got %o, want 0600", perm)
	}
}

func TestMergeVaults(t *testing.T) {
	dst := &Vault{data: map[string]string{"A": "1", "B": "2"}}
	src := &Vault{data: map[string]string{"B": "updated", "C": "3"}}

	count := MergeVaults(dst, src)
	if count != 2 {
		t.Errorf("MergeVaults count: got %d, want 2", count)
	}
	if dst.data["A"] != "1" {
		t.Errorf("A should be unchanged, got %q", dst.data["A"])
	}
	if dst.data["B"] != "updated" {
		t.Errorf("B should be updated, got %q", dst.data["B"])
	}
	if dst.data["C"] != "3" {
		t.Errorf("C should be added, got %q", dst.data["C"])
	}
}

func TestMergeVaultsNoChanges(t *testing.T) {
	dst := &Vault{data: map[string]string{"A": "1"}}
	src := &Vault{data: map[string]string{"A": "1"}}

	count := MergeVaults(dst, src)
	if count != 0 {
		t.Errorf("expected 0 changes, got %d", count)
	}
}
