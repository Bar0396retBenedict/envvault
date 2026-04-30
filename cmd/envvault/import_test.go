package main

import (
	"os"
	"path/filepath"
	"testing"

	"envvault/internal/vault"
)

func writeEnvSource(t *testing.T, dir, content string) string {
	t.Helper()
	p := filepath.Join(dir, "source.env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunImportNoArgs(t *testing.T) {
	err := runImport(newCmd(), []string{})
	if err == nil {
		t.Fatal("expected error with no args")
	}
}

func TestRunImportMissingPassphrase(t *testing.T) {
	t.Setenv("ENVVAULT_PASSPHRASE", "")
	dir := t.TempDir()
	src := writeEnvSource(t, dir, "A=1\n")
	vp := filepath.Join(dir, "v.vault")
	err := runImport(newCmd(), []string{vp, src})
	if err == nil || err.Error() != "ENVVAULT_PASSPHRASE environment variable is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunImportAddsKeys(t *testing.T) {
	pass := "s3cret"
	t.Setenv("ENVVAULT_PASSPHRASE", pass)
	dir := t.TempDir()
	vp := filepath.Join(dir, "v.vault")
	v := vault.New()
	if err := v.Save(vp, pass); err != nil {
		t.Fatal(err)
	}
	src := writeEnvSource(t, dir, "FOO=bar\nBAZ=qux\n")
	cmd := newCmd()
	if err := runImport(cmd, []string{vp, src}); err != nil {
		t.Fatal(err)
	}
	v2, err := vault.Load(vp, pass)
	if err != nil {
		t.Fatal(err)
	}
	if val, _ := v2.Get("FOO"); val != "bar" {
		t.Fatalf("expected bar, got %q", val)
	}
}

func TestRunImportUnknownFormat(t *testing.T) {
	pass := "s3cret"
	t.Setenv("ENVVAULT_PASSPHRASE", pass)
	dir := t.TempDir()
	vp := filepath.Join(dir, "v.vault")
	v := vault.New()
	if err := v.Save(vp, pass); err != nil {
		t.Fatal(err)
	}
	src := writeEnvSource(t, dir, "X=1\n")
	importFormat = "toml"
	defer func() { importFormat = "dotenv" }()
	err := runImport(newCmd(), []string{vp, src})
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestRunImportWrongPassphrase(t *testing.T) {
	pass := "correct"
	t.Setenv("ENVVAULT_PASSPHRASE", "wrong")
	dir := t.TempDir()
	vp := filepath.Join(dir, "v.vault")
	v := vault.New()
	if err := v.Save(vp, pass); err != nil {
		t.Fatal(err)
	}
	src := writeEnvSource(t, dir, "X=1\n")
	err := runImport(newCmd(), []string{vp, src})
	if err == nil {
		t.Fatal("expected error with wrong passphrase")
	}
}
