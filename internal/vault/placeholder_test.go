package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func makePlaceholderVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "test.vault")
}

func TestLoadPlaceholderRecordMissing(t *testing.T) {
	vp := makePlaceholderVault(t)
	rec, err := LoadPlaceholderRecord(vp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rec.Placeholders) != 0 {
		t.Errorf("expected empty record, got %v", rec.Placeholders)
	}
}

func TestSetAndListPlaceholder(t *testing.T) {
	vp := makePlaceholderVault(t)
	if err := SetPlaceholder(vp, "AWS_ACCESS_KEY_ID", "AWS access key (20 chars)"); err != nil {
		t.Fatalf("SetPlaceholder: %v", err)
	}
	if err := SetPlaceholder(vp, "DB_PASSWORD", "PostgreSQL password"); err != nil {
		t.Fatalf("SetPlaceholder: %v", err)
	}
	entries, err := ListPlaceholders(vp)
	if err != nil {
		t.Fatalf("ListPlaceholders: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Key != "AWS_ACCESS_KEY_ID" {
		t.Errorf("expected sorted first key AWS_ACCESS_KEY_ID, got %s", entries[0].Key)
	}
}

func TestSetPlaceholderUpdatesExisting(t *testing.T) {
	vp := makePlaceholderVault(t)
	_ = SetPlaceholder(vp, "API_KEY", "old description")
	_ = SetPlaceholder(vp, "API_KEY", "new description")
	entries, _ := ListPlaceholders(vp)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Description != "new description" {
		t.Errorf("expected updated description, got %q", entries[0].Description)
	}
}

func TestRemovePlaceholder(t *testing.T) {
	vp := makePlaceholderVault(t)
	_ = SetPlaceholder(vp, "TOKEN", "auth token")
	if err := RemovePlaceholder(vp, "TOKEN"); err != nil {
		t.Fatalf("RemovePlaceholder: %v", err)
	}
	entries, _ := ListPlaceholders(vp)
	if len(entries) != 0 {
		t.Errorf("expected empty list after removal")
	}
}

func TestRemovePlaceholderNotFound(t *testing.T) {
	vp := makePlaceholderVault(t)
	if err := RemovePlaceholder(vp, "MISSING_KEY"); err == nil {
		t.Error("expected error for missing key, got nil")
	}
}

func TestSetPlaceholderEmptyKey(t *testing.T) {
	vp := makePlaceholderVault(t)
	if err := SetPlaceholder(vp, "", "some description"); err == nil {
		t.Error("expected error for empty key")
	}
}

func TestSetPlaceholderEmptyDescription(t *testing.T) {
	vp := makePlaceholderVault(t)
	if err := SetPlaceholder(vp, "KEY", ""); err == nil {
		t.Error("expected error for empty description")
	}
}

func TestPlaceholderFilePermissions(t *testing.T) {
	vp := makePlaceholderVault(t)
	_ = SetPlaceholder(vp, "SECRET", "some secret value")
	pfp := placeholderFilePath(vp)
	info, err := os.Stat(pfp)
	if err != nil {
		t.Fatalf("stat placeholder file: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("expected permissions 0600, got %v", info.Mode().Perm())
	}
}
