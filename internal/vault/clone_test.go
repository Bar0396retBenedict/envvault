package vault_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nicholasgasior/envvault/internal/vault"
)

func makeCloneVault(t *testing.T, dir string, keys map[string]string, pass string) string {
	t.Helper()
	v := vault.New()
	for k, val := range keys {
		if err := v.Set(k, val); err != nil {
			t.Fatalf("set %s: %v", k, err)
		}
	}
	path := filepath.Join(dir, "src.vault")
	if err := v.Save(path, pass); err != nil {
		t.Fatalf("save: %v", err)
	}
	return path
}

func TestCloneSuccess(t *testing.T) {
	dir := t.TempDir()
	src := makeCloneVault(t, dir, map[string]string{"API_KEY": "abc", "DB_URL": "postgres://"}, "src-pass")
	dst := filepath.Join(dir, "dst.vault")

	if err := vault.CloneVault(src, "src-pass", dst, "dst-pass", vault.CloneOptions{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	v, err := vault.Load(dst, "dst-pass")
	if err != nil {
		t.Fatalf("load cloned vault: %v", err)
	}
	if val, ok := v.Get("API_KEY"); !ok || val != "abc" {
		t.Errorf("API_KEY: got %q, want %q", val, "abc")
	}
	if val, ok := v.Get("DB_URL"); !ok || val != "postgres://" {
		t.Errorf("DB_URL: got %q, want %q", val, "postgres://")
	}
}

func TestCloneDestinationExistsNoOverwrite(t *testing.T) {
	dir := t.TempDir()
	src := makeCloneVault(t, dir, map[string]string{"X": "1"}, "pass")
	dst := filepath.Join(dir, "dst.vault")
	_ = os.WriteFile(dst, []byte("existing"), 0o600)

	err := vault.CloneVault(src, "pass", dst, "pass", vault.CloneOptions{})
	if err == nil {
		t.Fatal("expected error when destination exists without Overwrite")
	}
}

func TestCloneDestinationExistsWithOverwrite(t *testing.T) {
	dir := t.TempDir()
	src := makeCloneVault(t, dir, map[string]string{"X": "1"}, "pass")
	dst := filepath.Join(dir, "dst.vault")
	_ = os.WriteFile(dst, []byte("existing"), 0o600)

	if err := vault.CloneVault(src, "pass", dst, "new-pass", vault.CloneOptions{Overwrite: true}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, err := vault.Load(dst, "new-pass")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if val, ok := v.Get("X"); !ok || val != "1" {
		t.Errorf("X: got %q, want %q", val, "1")
	}
}

func TestCloneEmptySourcePassphrase(t *testing.T) {
	dir := t.TempDir()
	err := vault.CloneVault(filepath.Join(dir, "src.vault"), "", filepath.Join(dir, "dst.vault"), "pass", vault.CloneOptions{})
	if err == nil {
		t.Fatal("expected error for empty source passphrase")
	}
}

func TestCloneEmptyDestinationPassphrase(t *testing.T) {
	dir := t.TempDir()
	src := makeCloneVault(t, dir, map[string]string{"Y": "2"}, "pass")
	err := vault.CloneVault(src, "pass", filepath.Join(dir, "dst.vault"), "", vault.CloneOptions{})
	if err == nil {
		t.Fatal("expected error for empty destination passphrase")
	}
}

func TestCloneMissingSource(t *testing.T) {
	dir := t.TempDir()
	err := vault.CloneVault(filepath.Join(dir, "missing.vault"), "pass", filepath.Join(dir, "dst.vault"), "pass", vault.CloneOptions{})
	if err == nil {
		t.Fatal("expected error for missing source vault")
	}
}
