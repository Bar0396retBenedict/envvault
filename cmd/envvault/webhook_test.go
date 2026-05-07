package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envvault/internal/vault"
)

func writeWebhookVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	vp := filepath.Join(dir, "test.vault")
	if err := os.WriteFile(vp, []byte(`{}`), 0600); err != nil {
		t.Fatal(err)
	}
	return vp
}

func captureWebhookOutput(t *testing.T, fn func() error) (string, error) {
	t.Helper()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := fn()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	return buf.String(), err
}

func TestRunWebhookAddSuccess(t *testing.T) {
	vp := writeWebhookVault(t)
	out, err := captureWebhookOutput(t, func() error {
		return runWebhookAdd(vp, "deploy", "https://ci.test/hook", []string{"set", "rotate"})
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "deploy") {
		t.Errorf("expected name in output, got: %s", out)
	}
}

func TestRunWebhookAddUnknownEvent(t *testing.T) {
	vp := writeWebhookVault(t)
	err := runWebhookAdd(vp, "hook", "https://x.test/", []string{"unknown"})
	if err == nil {
		t.Error("expected error for unknown event")
	}
}

func TestRunWebhookListEmpty(t *testing.T) {
	vp := writeWebhookVault(t)
	out, err := captureWebhookOutput(t, func() error {
		return runWebhookList(vp)
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "no webhooks") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestRunWebhookListWithEntries(t *testing.T) {
	vp := writeWebhookVault(t)
	_ = vault.RegisterWebhook(vp, "ci", "https://ci.example.com/", []vault.WebhookEvent{vault.EventSet})
	out, err := captureWebhookOutput(t, func() error {
		return runWebhookList(vp)
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "ci") {
		t.Errorf("expected webhook name in output, got: %s", out)
	}
	if !strings.Contains(out, "https://ci.example.com/") {
		t.Errorf("expected URL in output, got: %s", out)
	}
}

func TestRunWebhookRemoveSuccess(t *testing.T) {
	vp := writeWebhookVault(t)
	_ = vault.RegisterWebhook(vp, "old", "https://old.test/", nil)
	out, err := captureWebhookOutput(t, func() error {
		return runWebhookRemove(vp, "old")
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "removed") {
		t.Errorf("expected removed message, got: %s", out)
	}
}

func TestRunWebhookRemoveNotFound(t *testing.T) {
	vp := writeWebhookVault(t)
	err := runWebhookRemove(vp, "ghost")
	if err == nil {
		t.Error("expected error for missing webhook")
	}
}
