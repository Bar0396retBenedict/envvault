package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func makeLabelVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "test.vault")
}

func TestLoadLabelRecordMissing(t *testing.T) {
	vp := makeLabelVault(t)
	rec, err := LoadLabelRecord(vp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rec.Labels) != 0 {
		t.Fatalf("expected empty labels, got %v", rec.Labels)
	}
}

func TestSetAndGetLabel(t *testing.T) {
	vp := makeLabelVault(t)
	if err := SetLabel(vp, "DB_PASSWORD", "Database password"); err != nil {
		t.Fatalf("SetLabel: %v", err)
	}
	lbl, err := GetLabel(vp, "DB_PASSWORD")
	if err != nil {
		t.Fatalf("GetLabel: %v", err)
	}
	if lbl != "Database password" {
		t.Fatalf("expected 'Database password', got %q", lbl)
	}
}

func TestSetLabelUpdatesExisting(t *testing.T) {
	vp := makeLabelVault(t)
	_ = SetLabel(vp, "API_KEY", "Old label")
	_ = SetLabel(vp, "API_KEY", "New label")
	lbl, _ := GetLabel(vp, "API_KEY")
	if lbl != "New label" {
		t.Fatalf("expected 'New label', got %q", lbl)
	}
}

func TestRemoveLabel(t *testing.T) {
	vp := makeLabelVault(t)
	_ = SetLabel(vp, "TOKEN", "Auth token")
	if err := RemoveLabel(vp, "TOKEN"); err != nil {
		t.Fatalf("RemoveLabel: %v", err)
	}
	lbl, _ := GetLabel(vp, "TOKEN")
	if lbl != "" {
		t.Fatalf("expected empty label after removal, got %q", lbl)
	}
}

func TestRemoveLabelNotSet(t *testing.T) {
	vp := makeLabelVault(t)
	err := RemoveLabel(vp, "MISSING_KEY")
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestListLabels(t *testing.T) {
	vp := makeLabelVault(t)
	_ = SetLabel(vp, "Z_KEY", "Zeta")
	_ = SetLabel(vp, "A_KEY", "Alpha")
	_ = SetLabel(vp, "M_KEY", "Mu")
	keys, labels, err := ListLabels(vp)
	if err != nil {
		t.Fatalf("ListLabels: %v", err)
	}
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
	if keys[0] != "A_KEY" || keys[1] != "M_KEY" || keys[2] != "Z_KEY" {
		t.Fatalf("keys not sorted: %v", keys)
	}
	if labels[0] != "Alpha" {
		t.Fatalf("expected Alpha at index 0, got %q", labels[0])
	}
}

func TestSetLabelEmptyKey(t *testing.T) {
	vp := makeLabelVault(t)
	err := SetLabel(vp, "", "some label")
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestLabelFilePermissions(t *testing.T) {
	vp := makeLabelVault(t)
	_ = SetLabel(vp, "SECRET", "My secret")
	path := labelFilePath(vp)
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat label file: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("expected 0600, got %o", info.Mode().Perm())
	}
}
