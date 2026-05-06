package vault

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func makeNoteVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "test.vault")
}

func TestLoadNoteRecordMissing(t *testing.T) {
	vp := makeNoteVault(t)
	rec, err := LoadNoteRecord(vp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rec.Notes) != 0 {
		t.Errorf("expected empty notes, got %d", len(rec.Notes))
	}
}

func TestSetAndGetNote(t *testing.T) {
	vp := makeNoteVault(t)
	if err := SetNote(vp, "DB_URL", "primary database connection"); err != nil {
		t.Fatalf("SetNote: %v", err)
	}
	e, ok := GetNote(vp, "DB_URL")
	if !ok {
		t.Fatal("expected note to exist")
	}
	if e.Note != "primary database connection" {
		t.Errorf("unexpected note text: %q", e.Note)
	}
	if e.Key != "DB_URL" {
		t.Errorf("unexpected key: %q", e.Key)
	}
}

func TestSetNoteUpdatesExisting(t *testing.T) {
	vp := makeNoteVault(t)
	_ = SetNote(vp, "API_KEY", "old note")
	before, _ := GetNote(vp, "API_KEY")
	time.Sleep(2 * time.Millisecond)
	_ = SetNote(vp, "API_KEY", "new note")
	after, _ := GetNote(vp, "API_KEY")
	if after.Note != "new note" {
		t.Errorf("expected updated note, got %q", after.Note)
	}
	if !after.UpdatedAt.After(before.UpdatedAt) {
		t.Error("UpdatedAt should advance on update")
	}
}

func TestRemoveNote(t *testing.T) {
	vp := makeNoteVault(t)
	_ = SetNote(vp, "SECRET", "some note")
	if err := RemoveNote(vp, "SECRET"); err != nil {
		t.Fatalf("RemoveNote: %v", err)
	}
	_, ok := GetNote(vp, "SECRET")
	if ok {
		t.Error("expected note to be removed")
	}
}

func TestRemoveNoteNotExisting(t *testing.T) {
	vp := makeNoteVault(t)
	if err := RemoveNote(vp, "NONEXISTENT"); err != nil {
		t.Errorf("RemoveNote on missing key should not error: %v", err)
	}
}

func TestNoteFilePermissions(t *testing.T) {
	vp := makeNoteVault(t)
	_ = SetNote(vp, "X", "y")
	path := noteFilePath(vp)
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("expected 0600, got %o", info.Mode().Perm())
	}
}

func TestNoteTimestampNonZero(t *testing.T) {
	vp := makeNoteVault(t)
	_ = SetNote(vp, "K", "v")
	e, _ := GetNote(vp, "K")
	if e.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should not be zero")
	}
}
