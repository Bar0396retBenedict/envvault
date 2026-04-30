package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envvault/internal/vault"
)

func writeSnapshotVault(t *testing.T, pass string, entries map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "env.vault")
	v := vault.New(path, pass)
	for k, val := range entries {
		v.Set(k, val)
	}
	if err := v.Save(); err != nil {
		t.Fatalf("save vault: %v", err)
	}
	return path
}

func TestRunSnapshotTakeMissingPassphrase(t *testing.T) {
	path := writeSnapshotVault(t, "secret", map[string]string{"A": "1"})
	os.Unsetenv("ENVVAULT_PASSPHRASE")
	err := runSnapshotTake(path, "test")
	if err == nil || !strings.Contains(err.Error(), "ENVVAULT_PASSPHRASE") {
		t.Errorf("expected passphrase error, got %v", err)
	}
}

func TestRunSnapshotTakeSuccess(t *testing.T) {
	path := writeSnapshotVault(t, "secret", map[string]string{"FOO": "bar"})
	t.Setenv("ENVVAULT_PASSPHRASE", "secret")
	if err := runSnapshotTake(path, "pre-release"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	snaps, err := vault.ListSnapshots(path)
	if err != nil {
		t.Fatalf("list snapshots: %v", err)
	}
	if len(snaps) != 1 {
		t.Errorf("expected 1 snapshot, got %d", len(snaps))
	}
}

func TestRunSnapshotListEmpty(t *testing.T) {
	path := writeSnapshotVault(t, "secret", map[string]string{})
	var buf bytes.Buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runSnapshotList(path)
	w.Close()
	os.Stdout = old
	buf.ReadFrom(r)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No snapshots") {
		t.Errorf("expected 'No snapshots' message, got %q", buf.String())
	}
}

func TestRunSnapshotListWithEntries(t *testing.T) {
	path := writeSnapshotVault(t, "pass", map[string]string{"K": "v"})
	t.Setenv("ENVVAULT_PASSPHRASE", "pass")
	if err := runSnapshotTake(path, "snap-a"); err != nil {
		t.Fatalf("take: %v", err)
	}

	var buf bytes.Buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runSnapshotList(path)
	w.Close()
	os.Stdout = old
	buf.ReadFrom(r)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "snap-a") {
		t.Errorf("expected 'snap-a' in output, got %q", buf.String())
	}
}
