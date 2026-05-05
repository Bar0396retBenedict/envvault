package vault

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func makeTTLVault(t *testing.T) (v *Vault, path string) {
	t.Helper()
	dir := t.TempDir()
	path = filepath.Join(dir, "test.vault")
	v = New()
	v.Set("KEY_A", "alpha")
	v.Set("KEY_B", "beta")
	if err := v.Save(path, "passphrase"); err != nil {
		t.Fatalf("save vault: %v", err)
	}
	return v, path
}

func TestLoadTTLRecordMissing(t *testing.T) {
	dir := t.TempDir()
	rec, err := LoadTTLRecord(filepath.Join(dir, "missing.vault"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rec.Entries) != 0 {
		t.Errorf("expected empty entries, got %d", len(rec.Entries))
	}
}

func TestSetAndLoadTTL(t *testing.T) {
	_, path := makeTTLVault(t)
	if err := SetTTL(path, "KEY_A", 10*time.Minute); err != nil {
		t.Fatalf("SetTTL: %v", err)
	}
	rec, err := LoadTTLRecord(path)
	if err != nil {
		t.Fatalf("LoadTTLRecord: %v", err)
	}
	if len(rec.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(rec.Entries))
	}
	if rec.Entries[0].Key != "KEY_A" {
		t.Errorf("expected KEY_A, got %s", rec.Entries[0].Key)
	}
	if time.Until(rec.Entries[0].ExpiresAt) <= 0 {
		t.Error("expiry should be in the future")
	}
}

func TestSetTTLUpdatesExisting(t *testing.T) {
	_, path := makeTTLVault(t)
	_ = SetTTL(path, "KEY_A", 5*time.Minute)
	_ = SetTTL(path, "KEY_A", 30*time.Minute)
	rec, _ := LoadTTLRecord(path)
	if len(rec.Entries) != 1 {
		t.Errorf("expected 1 entry after update, got %d", len(rec.Entries))
	}
}

func TestExpiredKeys(t *testing.T) {
	_, path := makeTTLVault(t)
	_ = SetTTL(path, "KEY_A", -1*time.Second) // already expired
	_ = SetTTL(path, "KEY_B", 10*time.Minute)  // still valid

	expired, err := ExpiredKeys(path)
	if err != nil {
		t.Fatalf("ExpiredKeys: %v", err)
	}
	if len(expired) != 1 || expired[0] != "KEY_A" {
		t.Errorf("expected [KEY_A], got %v", expired)
	}
}

func TestExpiredKeysNone(t *testing.T) {
	_, path := makeTTLVault(t)
	_ = SetTTL(path, "KEY_A", 10*time.Minute)
	expired, err := ExpiredKeys(path)
	if err != nil {
		t.Fatalf("ExpiredKeys: %v", err)
	}
	if len(expired) != 0 {
		t.Errorf("expected no expired keys, got %v", expired)
	}
}

func TestPurgeTTLEntry(t *testing.T) {
	_, path := makeTTLVault(t)
	_ = SetTTL(path, "KEY_A", 10*time.Minute)
	_ = SetTTL(path, "KEY_B", 10*time.Minute)
	if err := PurgeTTLEntry(path, "KEY_A"); err != nil {
		t.Fatalf("PurgeTTLEntry: %v", err)
	}
	rec, _ := LoadTTLRecord(path)
	if len(rec.Entries) != 1 || rec.Entries[0].Key != "KEY_B" {
		t.Errorf("expected only KEY_B after purge, got %v", rec.Entries)
	}
}

func TestTTLFilePermissions(t *testing.T) {
	_, path := makeTTLVault(t)
	_ = SetTTL(path, "KEY_A", time.Hour)
	ttlPath := filepath.Join(filepath.Dir(path), "."+filepath.Base(path)+".ttl.json")
	info, err := os.Stat(ttlPath)
	if err != nil {
		t.Fatalf("stat ttl file: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %o", info.Mode().Perm())
	}
}
