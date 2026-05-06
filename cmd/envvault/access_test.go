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

func captureAccessOutput(vaultPath string) (string, error) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := runAccess(vaultPath)
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String(), err
}

func TestRunAccessNoRecords(t *testing.T) {
	dir := t.TempDir()
	vp := filepath.Join(dir, "test.vault")
	out, err := captureAccessOutput(vp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No access records") {
		t.Errorf("expected empty message, got: %q", out)
	}
}

func TestRunAccessWithEntries(t *testing.T) {
	dir := t.TempDir()
	vp := filepath.Join(dir, "test.vault")
	_ = vault.RecordRead(vp, "API_KEY")
	_ = vault.RecordRead(vp, "API_KEY")
	_ = vault.RecordWrite(vp, "DB_PASS")

	out, err := captureAccessOutput(vp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in output, got: %q", out)
	}
	if !strings.Contains(out, "DB_PASS") {
		t.Errorf("expected DB_PASS in output, got: %q", out)
	}
	if !strings.Contains(out, "2") {
		t.Errorf("expected read count 2 in output, got: %q", out)
	}
}

func TestRunAccessSortedOutput(t *testing.T) {
	dir := t.TempDir()
	vp := filepath.Join(dir, "test.vault")
	_ = vault.RecordRead(vp, "Z_VAR")
	_ = vault.RecordRead(vp, "A_VAR")

	out, err := captureAccessOutput(vp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	idxA := strings.Index(out, "A_VAR")
	idxZ := strings.Index(out, "Z_VAR")
	if idxA < 0 || idxZ < 0 || idxA > idxZ {
		t.Errorf("expected A_VAR before Z_VAR in output: %q", out)
	}
}

func TestRunAccessHeaderPresent(t *testing.T) {
	dir := t.TempDir()
	vp := filepath.Join(dir, "test.vault")
	_ = vault.RecordWrite(vp, "TOKEN")

	out, err := captureAccessOutput(vp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, col := range []string{"KEY", "READS", "WRITES", "LAST READ", "LAST WRITE"} {
		if !strings.Contains(out, col) {
			t.Errorf("missing column header %q in output: %q", col, out)
		}
	}
}
