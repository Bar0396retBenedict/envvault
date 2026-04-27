package vault_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envvault/internal/vault"
)

const testPassphrase = "test-passphrase-123"

func TestNewVault(t *testing.T) {
	v := vault.New()
	if v == nil {
		t.Fatal("expected non-nil vault")
	}
	if v.Version != 1 {
		t.Errorf("expected version 1, got %d", v.Version)
	}
	if len(v.Env) != 0 {
		t.Errorf("expected empty env map, got %d entries", len(v.Env))
	}
}

func TestSetAndGet(t *testing.T) {
	v := vault.New()
	v.Set("FOO", "bar")

	val, ok := v.Get("FOO")
	if !ok {
		t.Fatal("expected key FOO to exist")
	}
	if val != "bar" {
		t.Errorf("expected 'bar', got '%s'", val)
	}
}

func TestDelete(t *testing.T) {
	v := vault.New()
	v.Set("FOO", "bar")
	v.Delete("FOO")

	_, ok := v.Get("FOO")
	if ok {
		t.Error("expected FOO to be deleted")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.vault")

	v := vault.New()
	v.Set("DATABASE_URL", "postgres://localhost/mydb")
	v.Set("API_KEY", "secret-key-value")

	if err := v.Save(path, testPassphrase); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := vault.Load(path, testPassphrase)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	for key, expected := range v.Env {
		got, ok := loaded.Get(key)
		if !ok {
			t.Errorf("missing key %s after load", key)
			continue
		}
		if got != expected {
			t.Errorf("key %s: expected %q, got %q", key, expected, got)
		}
	}
}

func TestLoadWrongPassphrase(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.vault")

	v := vault.New()
	v.Set("SECRET", "value")

	if err := v.Save(path, testPassphrase); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	_, err := vault.Load(path, "wrong-passphrase")
	if err == nil {
		t.Error("expected error when loading with wrong passphrase")
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := vault.Load("/nonexistent/path/test.vault", testPassphrase)
	if err == nil {
		t.Error("expected error when loading missing file")
	}
}

func TestSaveCreatesFileWithRestrictedPermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.vault")

	v := vault.New()
	v.Set("KEY", "value")

	if err := v.Save(path, testPassphrase); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected file permissions 0600, got %o", info.Mode().Perm())
	}
}
