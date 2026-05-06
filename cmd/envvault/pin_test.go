package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"envvault/internal/vault"
)

func writePinVault(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "env.vault")
}

func TestRunPinAddAndList(t *testing.T) {
	vp := writePinVault(t)

	var buf bytes.Buffer
	cmd := buildRootCmd(&buf)
	cmd.SetArgs([]string{"pin", "add", vp, "API_KEY", "--note", "rotate quarterly"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("pin add: %v", err)
	}
	if !strings.Contains(buf.String(), "API_KEY") {
		t.Errorf("expected API_KEY in output, got: %s", buf.String())
	}

	buf.Reset()
	cmd = buildRootCmd(&buf)
	cmd.SetArgs([]string{"pin", "list", vp})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("pin list: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in list output, got: %s", out)
	}
	if !strings.Contains(out, "rotate quarterly") {
		t.Errorf("expected note in list output, got: %s", out)
	}
}

func TestRunPinListEmpty(t *testing.T) {
	vp := writePinVault(t)
	var buf bytes.Buffer
	cmd := buildRootCmd(&buf)
	cmd.SetArgs([]string{"pin", "list", vp})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("pin list: %v", err)
	}
	if !strings.Contains(buf.String(), "no pinned keys") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestRunPinRemove(t *testing.T) {
	vp := writePinVault(t)
	_ = vault.PinKey(vp, "DB_PASS", "")

	var buf bytes.Buffer
	cmd := buildRootCmd(&buf)
	cmd.SetArgs([]string{"pin", "remove", vp, "DB_PASS"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("pin remove: %v", err)
	}

	keys, _, _ := vault.ListPinnedKeys(vp)
	for _, k := range keys {
		if k == "DB_PASS" {
			t.Error("expected DB_PASS to be removed")
		}
	}
}

func TestRunPinRemoveNotPinned(t *testing.T) {
	vp := writePinVault(t)
	var buf bytes.Buffer
	cmd := buildRootCmd(&buf)
	cmd.SetArgs([]string{"pin", "remove", vp, "GHOST_KEY"})
	if err := cmd.Execute(); err == nil {
		t.Error("expected error removing unpinned key")
	}
}
