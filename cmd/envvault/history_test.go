package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"envvault/internal/vault"
)

func captureHistoryOutput(vaultPath string) (string, error) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := runHistory(vaultPath)
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String(), err
}

func TestRunHistoryNoEntries(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "test.vault")
	out, err := captureHistoryOutput(vaultPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No history") {
		t.Errorf("expected 'No history', got: %s", out)
	}
}

func TestRunHistoryWithEntries(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "test.vault")
	_ = vault.AppendHistory(vaultPath, vault.HistoryEntry{
		Key:      "SECRET",
		OldValue: "old",
		NewValue: "new",
		Action:   "set",
		Timestamp: time.Now().UTC(),
	})
	out, err := captureHistoryOutput(vaultPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "SECRET") {
		t.Errorf("expected SECRET in output, got: %s", out)
	}
	if !strings.Contains(out, "set") {
		t.Errorf("expected action 'set' in output, got: %s", out)
	}
}

func TestRunHistoryHeaderPresent(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "test.vault")
	_ = vault.AppendHistory(vaultPath, vault.HistoryEntry{Key: "K", Action: "delete"})
	out, _ := captureHistoryOutput(vaultPath)
	if !strings.Contains(out, "TIMESTAMP") || !strings.Contains(out, "ACTION") {
		t.Errorf("header missing from output: %s", out)
	}
}

func TestRunHistoryEmptyOldValueShowsDash(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "test.vault")
	_ = vault.AppendHistory(vaultPath, vault.HistoryEntry{
		Key:      "NEW_KEY",
		NewValue: "value",
		Action:   "set",
	})
	out, _ := captureHistoryOutput(vaultPath)
	if !strings.Contains(out, "-") {
		t.Errorf("expected dash for empty old value, got: %s", out)
	}
}
