package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func makeRestoreVault(t *testing.T, pass string, entries map[string]string) (string, *Vault) {
	t.Helper()
	dir := t.TempDir()
	v := New(pass)
	for k, val := range entries {
		v.Set(k, val)
	}
	path := filepath.Join(dir, "test.vault")
	if err := v.Save(path); err != nil {
		t.Fatalf("save: %v", err)
	}
	return path, v
}

func TestRestoreSuccess(t *testing.T) {
	pass := "restorepass"
	path, _ := makeRestoreVault(t, pass, map[string]string{"KEY": "value", "FOO": "bar"})

	if err := TakeSnapshot(path, pass, "before-restore"); err != nil {
		t.Fatalf("TakeSnapshot: %v", err)
	}

	// Overwrite the vault with different data.
	v2 := New(pass)
	v2.Set("KEY", "changed")
	if err := v2.Save(path); err != nil {
		t.Fatalf("overwrite save: %v", err)
	}

	if err := Restore(path, "before-restore", pass, pass); err != nil {
		t.Fatalf("Restore: %v", err)
	}

	restored, err := Load(path, pass)
	if err != nil {
		t.Fatalf("Load after restore: %v", err)
	}
	if v, _ := restored.Get("FOO"); v != "bar" {
		t.Errorf("expected FOO=bar, got %q", v)
	}
	if v, _ := restored.Get("KEY"); v != "value" {
		t.Errorf("expected KEY=value, got %q", v)
	}
}

func TestRestoreMissingSnapshot(t *testing.T) {
	path, _ := makeRestoreVault(t, "pass", nil)
	err := Restore(path, "nonexistent", "pass", "pass")
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestRestoreEmptyLabel(t *testing.T) {
	path, _ := makeRestoreVault(t, "pass", nil)
	err := Restore(path, "", "pass", "pass")
	if err == nil {
		t.Fatal("expected error for empty label")
	}
}

func TestRestoreWrongSnapshotPassphrase(t *testing.T) {
	pass := "correct"
	path, _ := makeRestoreVault(t, pass, map[string]string{"X": "1"})
	if err := TakeSnapshot(path, pass, "snap"); err != nil {
		t.Fatalf("TakeSnapshot: %v", err)
	}
	err := Restore(path, "snap", "wrongpass", pass)
	if err == nil {
		t.Fatal("expected error with wrong snapshot passphrase")
	}
}

func TestRestoreChangesPassphrase(t *testing.T) {
	oldPass := "oldpass"
	newPass := "newpass"
	path, _ := makeRestoreVault(t, oldPass, map[string]string{"ENV": "prod"})
	if err := TakeSnapshot(path, oldPass, "migrate"); err != nil {
		t.Fatalf("TakeSnapshot: %v", err)
	}

	if err := Restore(path, "migrate", oldPass, newPass); err != nil {
		t.Fatalf("Restore: %v", err)
	}

	// Must be loadable with new passphrase.
	v, err := Load(path, newPass)
	if err != nil {
		t.Fatalf("Load with new passphrase: %v", err)
	}
	if val, _ := v.Get("ENV"); val != "prod" {
		t.Errorf("expected ENV=prod, got %q", val)
	}
	// Must NOT be loadable with old passphrase.
	if _, err := Load(path, oldPass); err == nil {
		t.Error("expected error loading with old passphrase after re-encrypt")
	}
	_ = os.Remove(path)
}
