package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envvault/internal/vault"
)

func captureAuditOutput(vaultPath string) (string, error) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := runAudit(vaultPath)
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String(), err
}

func TestRunAuditNoEntries(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "empty.vault")
	out, err := captureAuditOutput(vaultPath)
	if err != nil {
		t.Fatalf("runAudit: %v", err)
	}
	if !strings.Contains(out, "No audit entries") {
		t.Errorf("expected empty message, got: %q", out)
	}
}

func TestRunAuditWithEntries(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "prod.vault")

	if err := vault.AppendAuditEntry(vaultPath, vault.AuditSet, "SECRET"); err != nil {
		t.Fatalf("AppendAuditEntry: %v", err)
	}
	if err := vault.AppendAuditEntry(vaultPath, vault.AuditDelete, "OLD"); err != nil {
		t.Fatalf("AppendAuditEntry: %v", err)
	}

	out, err := captureAuditOutput(vaultPath)
	if err != nil {
		t.Fatalf("runAudit: %v", err)
	}
	if !strings.Contains(out, "set") {
		t.Errorf("expected 'set' in output, got: %q", out)
	}
	if !strings.Contains(out, "SECRET") {
		t.Errorf("expected 'SECRET' in output, got: %q", out)
	}
	if !strings.Contains(out, "delete") {
		t.Errorf("expected 'delete' in output, got: %q", out)
	}
	if !strings.Contains(out, "TIMESTAMP") {
		t.Errorf("expected header in output, got: %q", out)
	}
}

func TestRunAuditEmptyKeyShowsDash(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "test.vault")
	if err := vault.AppendAuditEntry(vaultPath, vault.AuditRotate, ""); err != nil {
		t.Fatalf("AppendAuditEntry: %v", err)
	}
	out, err := captureAuditOutput(vaultPath)
	if err != nil {
		t.Fatalf("runAudit: %v", err)
	}
	if !strings.Contains(out, "-") {
		t.Errorf("expected dash for empty key, got: %q", out)
	}
}
