package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envvault/internal/vault"
)

func writeTempVault(t *testing.T, dir, name, passphrase string, entries map[string]string) string {
	t.Helper()
	v := vault.New()
	for k, val := range entries {
		if err := v.Set(k, val); err != nil {
			t.Fatalf("Set: %v", err)
		}
	}
	path := filepath.Join(dir, name)
	if err := v.Save(path, passphrase); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return path
}

func TestRunDiffNoArgs(t *testing.T) {
	err := runDiff([]string{})
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestRunDiffMissingPassphrase(t *testing.T) {
	t.Setenv("ENVVAULT_PASSPHRASE", "")
	err := runDiff([]string{"a.vault", "b.vault"})
	if err == nil {
		t.Fatal("expected error when passphrase missing")
	}
}

func TestRunDiffIdenticalVaults(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVVAULT_PASSPHRASE", "secret")

	entries := map[string]string{"FOO": "bar", "BAZ": "qux"}
	srcPath := writeTempVault(t, dir, "src.vault", "secret", entries)
	dstPath := writeTempVault(t, dir, "dst.vault", "secret", entries)

	if err := runDiff([]string{srcPath, dstPath}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunDiffDifferentVaults(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVVAULT_PASSPHRASE", "secret")

	srcPath := writeTempVault(t, dir, "src.vault", "secret", map[string]string{
		"KEEP": "same",
		"OLD":  "value",
	})
	dstPath := writeTempVault(t, dir, "dst.vault", "secret", map[string]string{
		"KEEP": "same",
		"NEW":  "added",
	})

	if err := runDiff([]string{srcPath, dstPath}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunDiffBadVaultFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVVAULT_PASSPHRASE", "secret")

	badPath := filepath.Join(dir, "bad.vault")
	if err := os.WriteFile(badPath, []byte("not-valid"), 0600); err != nil {
		t.Fatal(err)
	}
	goodPath := writeTempVault(t, dir, "good.vault", "secret", map[string]string{"X": "1"})

	if err := runDiff([]string{badPath, goodPath}); err == nil {
		t.Fatal("expected error for bad vault file")
	}
}
