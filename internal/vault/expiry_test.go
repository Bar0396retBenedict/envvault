package vault

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func makeExpiryVault(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.vault")
	passphrase := "expiry-secret"
	v := New()
	v.Set("API_KEY", "abc123")
	v.Set("DB_PASS", "hunter2")
	if err := v.Save(path, passphrase); err != nil {
		t.Fatalf("save vault: %v", err)
	}
	return path, passphrase
}

func TestLoadExpiryRecordMissing(t *testing.T) {
	dir := t.TempDir()
	rec, err := LoadExpiryRecord(filepath.Join(dir, "none.vault"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rec) != 0 {
		t.Errorf("expected empty record, got %v", rec)
	}
}

func TestSaveAndLoadExpiryRecord(t *testing.T) {
	path, _ := makeExpiryVault(t)
	exp := time.Now().Add(24 * time.Hour).Truncate(time.Second)
	rec := ExpiryRecord{"API_KEY": exp}
	if err := SaveExpiryRecord(path, rec); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadExpiryRecord(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if !loaded["API_KEY"].Equal(exp) {
		t.Errorf("expected %v, got %v", exp, loaded["API_KEY"])
	}
}

func TestExpiryFilePermissions(t *testing.T) {
	path, _ := makeExpiryVault(t)
	if err := SaveExpiryRecord(path, ExpiryRecord{}); err != nil {
		t.Fatalf("save: %v", err)
	}
	info, err := os.Stat(expiryFilePath(path))
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("expected 0600, got %04o", perm)
	}
}

func TestSetExpiry(t *testing.T) {
	path, pass := makeExpiryVault(t)
	exp := time.Now().Add(48 * time.Hour)
	if err := SetExpiry(path, pass, "API_KEY", exp); err != nil {
		t.Fatalf("SetExpiry: %v", err)
	}
	rec, err := LoadExpiryRecord(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if _, ok := rec["API_KEY"]; !ok {
		t.Error("expected API_KEY in expiry record")
	}
}

func TestSetExpiryMissingKey(t *testing.T) {
	path, pass := makeExpiryVault(t)
	err := SetExpiry(path, pass, "NONEXISTENT", time.Now().Add(time.Hour))
	if err == nil {
		t.Error("expected error for missing key")
	}
}

func TestExpiredKeysList(t *testing.T) {
	path, pass := makeExpiryVault(t)
	past := time.Now().Add(-1 * time.Hour)
	future := time.Now().Add(1 * time.Hour)
	if err := SetExpiry(path, pass, "API_KEY", past); err != nil {
		t.Fatalf("set past expiry: %v", err)
	}
	if err := SetExpiry(path, pass, "DB_PASS", future); err != nil {
		t.Fatalf("set future expiry: %v", err)
	}
	expired, err := ExpiredKeysList(path)
	if err != nil {
		t.Fatalf("ExpiredKeysList: %v", err)
	}
	if len(expired) != 1 || expired[0] != "API_KEY" {
		t.Errorf("expected [API_KEY], got %v", expired)
	}
}
