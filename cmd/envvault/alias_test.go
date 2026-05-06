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

func captureAliasOutput(fn func() error) (string, error) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := fn()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String(), err
}

func writeAliasVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "test.vault")
}

func TestRunAliasSetSuccess(t *testing.T) {
	vp := writeAliasVault(t)
	out, err := captureAliasOutput(func() error {
		return runAliasSet(vp, "db", "DATABASE_URL")
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "db") || !strings.Contains(out, "DATABASE_URL") {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestRunAliasSetEmptyAlias(t *testing.T) {
	vp := writeAliasVault(t)
	err := runAliasSet(vp, "", "DATABASE_URL")
	if err == nil {
		t.Fatal("expected error for empty alias")
	}
}

func TestRunAliasListEmpty(t *testing.T) {
	vp := writeAliasVault(t)
	out, err := captureAliasOutput(func() error {
		return runAliasList(vp)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "no aliases") {
		t.Fatalf("expected empty message, got: %s", out)
	}
}

func TestRunAliasListWithEntries(t *testing.T) {
	vp := writeAliasVault(t)
	_ = vault.SetAlias(vp, "db", "DATABASE_URL")
	_ = vault.SetAlias(vp, "secret", "API_SECRET")
	out, err := captureAliasOutput(func() error {
		return runAliasList(vp)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "db") || !strings.Contains(out, "DATABASE_URL") {
		t.Fatalf("missing db entry in output: %s", out)
	}
	if !strings.Contains(out, "ALIAS") {
		t.Fatalf("missing header in output: %s", out)
	}
}

func TestRunAliasRemoveSuccess(t *testing.T) {
	vp := writeAliasVault(t)
	_ = vault.SetAlias(vp, "db", "DATABASE_URL")
	out, err := captureAliasOutput(func() error {
		return runAliasRemove(vp, "db")
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "removed") {
		t.Fatalf("expected removed message, got: %s", out)
	}
}

func TestRunAliasRemoveNotFound(t *testing.T) {
	vp := writeAliasVault(t)
	err := runAliasRemove(vp, "ghost")
	if err == nil {
		t.Fatal("expected error for non-existent alias")
	}
}
