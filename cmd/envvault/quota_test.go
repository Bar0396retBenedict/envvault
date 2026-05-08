package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/envvault/internal/vault"
)

func writeQuotaVault(t *testing.T, keys map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.vault")
	v := vault.New()
	for k, val := range keys {
		v.Set(k, val)
	}
	if err := v.Save(path, "secret"); err != nil {
		t.Fatalf("save vault: %v", err)
	}
	return path
}

func TestRunQuotaSetSuccess(t *testing.T) {
	path := writeQuotaVault(t, nil)
	cmd := newRootCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"quota", "set", path, "--max-keys", "50", "--max-bytes", "8192"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "max_keys=50") {
		t.Errorf("expected max_keys=50 in output, got: %s", out)
	}
	if !strings.Contains(out, "max_bytes=8192") {
		t.Errorf("expected max_bytes=8192 in output, got: %s", out)
	}
}

func TestRunQuotaSetUnlimited(t *testing.T) {
	path := writeQuotaVault(t, nil)
	cmd := newRootCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"quota", "set", path})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "unlimited") {
		t.Errorf("expected 'unlimited' in output, got: %s", out)
	}
}

func TestRunQuotaCheckMissingPassphrase(t *testing.T) {
	path := writeQuotaVault(t, nil)
	cmd := newRootCmd()
	cmd.SetArgs([]string{"quota", "check", path})
	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "passphrase") {
		t.Errorf("expected passphrase error, got: %v", err)
	}
}

func TestRunQuotaCheckNoViolations(t *testing.T) {
	path := writeQuotaVault(t, map[string]string{"A": "1", "B": "2"})
	if err := vault.SaveQuotaRecord(path, vault.QuotaRecord{MaxKeys: 10, MaxBytes: 1024}); err != nil {
		t.Fatalf("save quota: %v", err)
	}
	cmd := newRootCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"quota", "check", path, "--passphrase", "secret"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "OK") {
		t.Errorf("expected OK, got: %s", buf.String())
	}
}

func TestRunQuotaCheckViolation(t *testing.T) {
	path := writeQuotaVault(t, map[string]string{"A": "1", "B": "2", "C": "3"})
	if err := vault.SaveQuotaRecord(path, vault.QuotaRecord{MaxKeys: 2}); err != nil {
		t.Fatalf("save quota: %v", err)
	}
	cmd := newRootCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"quota", "check", path, "--passphrase", "secret"})
	// exits 1 on violation — check output regardless
	cmd.Execute() //nolint:errcheck
	out := buf.String()
	if !strings.Contains(out, "VIOLATION") {
		t.Errorf("expected VIOLATION in output, got: %s", out)
	}
	_ = os.Stderr // suppress unused import
}
