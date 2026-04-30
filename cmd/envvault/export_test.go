package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envvault/internal/vault"
)

func TestRunExportNoArgs(t *testing.T) {
	cmd := exportCmd
	cmd.ResetFlags()
	cmd.Flags().StringP("passphrase", "p", "", "")
	cmd.Flags().StringP("format", "f", "dotenv", "")

	err := cmd.RunE(cmd, []string{})
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestRunExportMissingPassphrase(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "test.vault")

	v := vault.New()
	if err := vault.Save(v, vaultPath, "secret"); err != nil {
		t.Fatalf("failed to create vault: %v", err)
	}

	os.Unsetenv("ENVVAULT_PASSPHRASE")

	cmd := exportCmd
	cmd.ResetFlags()
	cmd.Flags().StringP("passphrase", "p", "", "")
	cmd.Flags().StringP("format", "f", "dotenv", "")

	err := cmd.RunE(cmd, []string{vaultPath})
	if err == nil || !strings.Contains(err.Error(), "passphrase is required") {
		t.Fatalf("expected passphrase error, got: %v", err)
	}
}

func TestRunExportDotenv(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "test.vault")

	v := vault.New()
	v.Set("APP_ENV", "production")
	v.Set("DB_HOST", "localhost")
	if err := vault.Save(v, vaultPath, "secret"); err != nil {
		t.Fatalf("failed to save vault: %v", err)
	}

	cmd := exportCmd
	cmd.ResetFlags()
	cmd.Flags().StringP("passphrase", "p", "", "")
	cmd.Flags().StringP("format", "f", "dotenv", "")
	_ = cmd.Flags().Set("passphrase", "secret")
	_ = cmd.Flags().Set("format", "dotenv")

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := cmd.RunE(cmd, []string{vaultPath})
	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "APP_ENV=") {
		t.Errorf("expected APP_ENV in output, got:\n%s", output)
	}
	if !strings.Contains(output, "DB_HOST=") {
		t.Errorf("expected DB_HOST in output, got:\n%s", output)
	}
}

func TestRunExportWrongPassphrase(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "test.vault")

	v := vault.New()
	v.Set("KEY", "value")
	if err := vault.Save(v, vaultPath, "correct"); err != nil {
		t.Fatalf("failed to save vault: %v", err)
	}

	cmd := exportCmd
	cmd.ResetFlags()
	cmd.Flags().StringP("passphrase", "p", "", "")
	cmd.Flags().StringP("format", "f", "dotenv", "")
	_ = cmd.Flags().Set("passphrase", "wrong")

	err := cmd.RunE(cmd, []string{vaultPath})
	if err == nil {
		t.Fatal("expected error for wrong passphrase")
	}
}
