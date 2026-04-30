package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRotateSuccess(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.enc")

	v, err := New(path, "old-secret")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	v.Set("KEY1", "value1")
	v.Set("KEY2", "value2")
	if err := v.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	rec, err := Rotate(path, "old-secret", "new-secret")
	if err != nil {
		t.Fatalf("Rotate: %v", err)
	}
	if rec == nil {
		t.Fatal("expected non-nil RotateRecord")
	}

	// Old passphrase should no longer work.
	_, err = New(path, "old-secret")
	if err == nil {
		t.Fatal("expected error loading vault with old passphrase after rotation")
	}

	// New passphrase should work and data must be intact.
	nv, err := New(path, "new-secret")
	if err != nil {
		t.Fatalf("New with new passphrase: %v", err)
	}
	for _, key := range []string{"KEY1", "KEY2"} {
		if _, err := nv.Get(key); err != nil {
			t.Errorf("Get(%q) after rotate: %v", key, err)
		}
	}
}

func TestRotateSamePassphrase(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.enc")

	v, _ := New(path, "pass")
	v.Save()

	_, err := Rotate(path, "pass", "pass")
	if err == nil {
		t.Fatal("expected error when old and new passphrases are identical")
	}
}

func TestRotateEmptyPassphrase(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.enc")

	_, err := Rotate(path, "", "newpass")
	if err == nil {
		t.Fatal("expected error for empty old passphrase")
	}

	_, err = Rotate(path, "oldpass", "")
	if err == nil {
		t.Fatal("expected error for empty new passphrase")
	}
}

func TestRotateMissingVault(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent.enc")
	_, err := Rotate(path, "old", "new")
	if err == nil {
		t.Fatal("expected error rotating non-existent vault")
	}
	_ = os.Remove(path)
}
