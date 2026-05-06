package vault

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func makeHistoryVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "test.vault")
}

func TestLoadHistoryMissing(t *testing.T) {
	vaultPath := makeHistoryVault(t)
	rec, err := LoadHistory(vaultPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rec.Entries) != 0 {
		t.Fatalf("expected empty entries, got %d", len(rec.Entries))
	}
}

func TestAppendHistory(t *testing.T) {
	vaultPath := makeHistoryVault(t)
	entry := HistoryEntry{
		Key:      "API_KEY",
		OldValue: "",
		NewValue: "abc123",
		Action:   "set",
	}
	if err := AppendHistory(vaultPath, entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rec, err := LoadHistory(vaultPath)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if len(rec.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(rec.Entries))
	}
	if rec.Entries[0].Key != "API_KEY" {
		t.Errorf("expected key API_KEY, got %s", rec.Entries[0].Key)
	}
}

func TestAppendHistoryMultiple(t *testing.T) {
	vaultPath := makeHistoryVault(t)
	for i := 0; i < 3; i++ {
		_ = AppendHistory(vaultPath, HistoryEntry{Key: "K", Action: "set"})
	}
	rec, _ := LoadHistory(vaultPath)
	if len(rec.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(rec.Entries))
	}
}

func TestHistoryEntryTimestampSet(t *testing.T) {
	vaultPath := makeHistoryVault(t)
	_ = AppendHistory(vaultPath, HistoryEntry{Key: "X", Action: "delete"})
	rec, _ := LoadHistory(vaultPath)
	if rec.Entries[0].Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestHistoryFilePermissions(t *testing.T) {
	vaultPath := makeHistoryVault(t)
	_ = AppendHistory(vaultPath, HistoryEntry{Key: "Y", Action: "set"})
	p := historyFilePath(vaultPath)
	info, err := os.Stat(p)
	if err != nil {
		t.Fatalf("stat error: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}

func TestHistoryPreservesOldAndNewValue(t *testing.T) {
	vaultPath := makeHistoryVault(t)
	entry := HistoryEntry{
		Key:       "DB_PASS",
		OldValue:  "old",
		NewValue:  "new",
		Action:    "set",
		Timestamp: time.Now().UTC(),
	}
	_ = AppendHistory(vaultPath, entry)
	rec, _ := LoadHistory(vaultPath)
	if rec.Entries[0].OldValue != "old" || rec.Entries[0].NewValue != "new" {
		t.Errorf("values not preserved: old=%s new=%s", rec.Entries[0].OldValue, rec.Entries[0].NewValue)
	}
}
