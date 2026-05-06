package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"envvault/internal/vault"
)

func writeEnvVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "env.vault")
	v := vault.New()
	v.Set("DB_PASS", "hunter2")
	v.Set("API_KEY", "key-xyz")
	if err := v.Save(path, "pass"); err != nil {
		t.Fatalf("save: %v", err)
	}
	return path
}

func TestRunEnvBindSuccess(t *testing.T) {
	path := writeEnvVault(t)
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"env", "bind", path, "MY_API", "API_KEY"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.Contains(buf.String(), "MY_API") {
		t.Errorf("expected MY_API in output, got: %s", buf.String())
	}
}

func TestRunEnvListEmpty(t *testing.T) {
	path := writeEnvVault(t)
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"env", "list", path})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.Contains(buf.String(), "no env bindings") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestRunEnvListWithEntries(t *testing.T) {
	path := writeEnvVault(t)
	if err := vault.BindEnvVar(path, "MY_DB", "DB_PASS"); err != nil {
		t.Fatalf("bind: %v", err)
	}
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"env", "list", path})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "MY_DB") || !strings.Contains(out, "DB_PASS") {
		t.Errorf("expected MY_DB and DB_PASS in output, got: %s", out)
	}
}

func TestRunEnvUnbindSuccess(t *testing.T) {
	path := writeEnvVault(t)
	if err := vault.BindEnvVar(path, "MY_DB", "DB_PASS"); err != nil {
		t.Fatalf("bind: %v", err)
	}
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"env", "unbind", path, "MY_DB"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.Contains(buf.String(), "unbound") {
		t.Errorf("expected 'unbound' in output, got: %s", buf.String())
	}
}

func TestRunEnvUnbindNotBound(t *testing.T) {
	path := writeEnvVault(t)
	rootCmd.SetArgs([]string{"env", "unbind", path, "GHOST"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error for unbound key")
	}
}
