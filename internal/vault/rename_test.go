package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func makeRenameVault(t *testing.T) (*Vault, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "rename.vault")
	v, err := New(path, "passphrase")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	v.Set("ALPHA", "one")
	v.Set("BETA", "two")
	if err := v.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return v, path
}

func TestRenameSuccess(t *testing.T) {
	v, _ := makeRenameVault(t)
	if err := RenameKey(v, "ALPHA", "ALPHA_RENAMED", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := v.Get("ALPHA"); ok {
		t.Error("old key ALPHA should no longer exist")
	}
	val, ok := v.Get("ALPHA_RENAMED")
	if !ok {
		t.Fatal("new key ALPHA_RENAMED should exist")
	}
	if val != "one" {
		t.Errorf("expected value %q, got %q", "one", val)
	}
}

func TestRenameSourceMissing(t *testing.T) {
	v, _ := makeRenameVault(t)
	err := RenameKey(v, "MISSING", "NEW_KEY", false)
	if err == nil {
		t.Fatal("expected error for missing source key")
	}
}

func TestRenameDestinationExistsNoOverwrite(t *testing.T) {
	v, _ := makeRenameVault(t)
	err := RenameKey(v, "ALPHA", "BETA", false)
	if err == nil {
		t.Fatal("expected error when destination exists and overwrite is false")
	}
}

func TestRenameDestinationExistsWithOverwrite(t *testing.T) {
	v, _ := makeRenameVault(t)
	if err := RenameKey(v, "ALPHA", "BETA", true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	val, ok := v.Get("BETA")
	if !ok {
		t.Fatal("key BETA should exist after overwrite rename")
	}
	if val != "one" {
		t.Errorf("expected value %q, got %q", "one", val)
	}
	if _, ok := v.Get("ALPHA"); ok {
		t.Error("old key ALPHA should no longer exist")
	}
}

func TestRenameSameKey(t *testing.T) {
	v, _ := makeRenameVault(t)
	err := RenameKey(v, "ALPHA", "ALPHA", false)
	if err == nil {
		t.Fatal("expected error when source and destination are identical")
	}
}

func TestRenameEmptyKeys(t *testing.T) {
	v, _ := makeRenameVault(t)
	if err := RenameKey(v, "", "NEW", false); err == nil {
		t.Error("expected error for empty source key")
	}
	if err := RenameKey(v, "ALPHA", "", false); err == nil {
		t.Error("expected error for empty destination key")
	}
}

func TestRenamePersists(t *testing.T) {
	v, path := makeRenameVault(t)
	if err := RenameKey(v, "ALPHA", "ALPHA_NEW", false); err != nil {
		t.Fatalf("RenameKey: %v", err)
	}
	if err := v.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}
	v2, err := New(path, "passphrase")
	if err != nil {
		t.Fatalf("reload New: %v", err)
	}
	if err := v2.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if _, ok := v2.Get("ALPHA"); ok {
		t.Error("old key ALPHA should not persist after save/load")
	}
	val, ok := v2.Get("ALPHA_NEW")
	if !ok {
		t.Fatal("renamed key ALPHA_NEW should persist after save/load")
	}
	if val != "one" {
		t.Errorf("expected %q, got %q", "one", val)
	}
	_ = os.Remove(path)
}
