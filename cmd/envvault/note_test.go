package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envvault/internal/vault"
)

func captureNoteOutput(fn func() error) (string, error) {
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

func writeNoteVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "test.vault")
}

func TestRunNoteSetSuccess(t *testing.T) {
	vp := writeNoteVault(t)
	out, err := captureNoteOutput(func() error {
		return runNoteSet(vp, "DB_URL", "primary db")
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "DB_URL") {
		t.Errorf("expected key in output, got: %q", out)
	}
}

func TestRunNoteListEmpty(t *testing.T) {
	vp := writeNoteVault(t)
	out, err := captureNoteOutput(func() error {
		return runNoteList(vp)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "no notes found") {
		t.Errorf("expected empty message, got: %q", out)
	}
}

func TestRunNoteListWithEntries(t *testing.T) {
	vp := writeNoteVault(t)
	_ = vault.SetNote(vp, "API_KEY", "rotated monthly")
	_ = vault.SetNote(vp, "SECRET", "do not share")
	out, err := captureNoteOutput(func() error {
		return runNoteList(vp)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{"API_KEY", "SECRET", "rotated monthly", "do not share"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output:\n%s", want, out)
		}
	}
}

func TestRunNoteRemove(t *testing.T) {
	vp := writeNoteVault(t)
	_ = vault.SetNote(vp, "TOKEN", "temp token")
	out, err := captureNoteOutput(func() error {
		return runNoteRemove(vp, "TOKEN")
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "TOKEN") {
		t.Errorf("expected key in output: %q", out)
	}
	_, ok := vault.GetNote(vp, "TOKEN")
	if ok {
		t.Error("note should have been removed")
	}
}

func TestRunNoteListHeaderPresent(t *testing.T) {
	vp := writeNoteVault(t)
	_ = vault.SetNote(vp, "X", fmt.Sprintf("note %d", 1))
	out, _ := captureNoteOutput(func() error {
		return runNoteList(vp)
	})
	if !strings.Contains(out, "KEY") || !strings.Contains(out, "NOTE") {
		t.Errorf("header not present in output: %q", out)
	}
}
