package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func makeWebhookVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "test.vault")
	// create a minimal vault file so the directory exists
	if err := os.WriteFile(vaultPath, []byte(`{}`), 0600); err != nil {
		t.Fatal(err)
	}
	return vaultPath
}

func TestLoadWebhookRecordMissing(t *testing.T) {
	vp := makeWebhookVault(t)
	rec, err := LoadWebhookRecord(vp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rec.Hooks) != 0 {
		t.Errorf("expected empty hooks, got %d", len(rec.Hooks))
	}
}

func TestRegisterAndListWebhook(t *testing.T) {
	vp := makeWebhookVault(t)
	err := RegisterWebhook(vp, "deploy", "https://example.com/hook", []WebhookEvent{EventSet, EventDelete})
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	names, rec, err := ListWebhooks(vp)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(names) != 1 || names[0] != "deploy" {
		t.Errorf("expected [deploy], got %v", names)
	}
	entry := rec.Hooks["deploy"]
	if entry.URL != "https://example.com/hook" {
		t.Errorf("unexpected URL: %s", entry.URL)
	}
	if len(entry.Events) != 2 {
		t.Errorf("expected 2 events, got %d", len(entry.Events))
	}
}

func TestRegisterWebhookEventsSorted(t *testing.T) {
	vp := makeWebhookVault(t)
	events := []WebhookEvent{EventRotate, EventSet, EventDelete}
	if err := RegisterWebhook(vp, "hook", "https://h.test/", events); err != nil {
		t.Fatal(err)
	}
	rec, err := LoadWebhookRecord(vp)
	if err != nil {
		t.Fatal(err)
	}
	got := rec.Hooks["hook"].Events
	for i := 1; i < len(got); i++ {
		if got[i] < got[i-1] {
			t.Errorf("events not sorted: %v", got)
		}
	}
}

func TestRegisterWebhookUpdatesExisting(t *testing.T) {
	vp := makeWebhookVault(t)
	_ = RegisterWebhook(vp, "ci", "https://old.url/", []WebhookEvent{EventSet})
	_ = RegisterWebhook(vp, "ci", "https://new.url/", []WebhookEvent{EventDelete})
	rec, _ := LoadWebhookRecord(vp)
	if rec.Hooks["ci"].URL != "https://new.url/" {
		t.Errorf("expected updated URL")
	}
}

func TestDeregisterWebhook(t *testing.T) {
	vp := makeWebhookVault(t)
	_ = RegisterWebhook(vp, "hook", "https://x.test/", []WebhookEvent{EventSet})
	if err := DeregisterWebhook(vp, "hook"); err != nil {
		t.Fatalf("deregister: %v", err)
	}
	names, _, _ := ListWebhooks(vp)
	if len(names) != 0 {
		t.Errorf("expected empty list after deregister")
	}
}

func TestDeregisterWebhookNotFound(t *testing.T) {
	vp := makeWebhookVault(t)
	err := DeregisterWebhook(vp, "ghost")
	if err == nil {
		t.Error("expected error for missing webhook")
	}
}

func TestRegisterWebhookEmptyName(t *testing.T) {
	vp := makeWebhookVault(t)
	err := RegisterWebhook(vp, "", "https://x.test/", nil)
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestRegisterWebhookEmptyURL(t *testing.T) {
	vp := makeWebhookVault(t)
	err := RegisterWebhook(vp, "hook", "", nil)
	if err == nil {
		t.Error("expected error for empty URL")
	}
}
