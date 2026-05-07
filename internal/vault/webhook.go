package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// WebhookEvent represents the type of vault event that triggered a webhook.
type WebhookEvent string

const (
	EventSet    WebhookEvent = "set"
	EventDelete WebhookEvent = "delete"
	EventRotate WebhookEvent = "rotate"
)

// WebhookEntry stores a registered webhook URL and the events it subscribes to.
type WebhookEntry struct {
	URL       string         `json:"url"`
	Events    []WebhookEvent `json:"events"`
	CreatedAt time.Time      `json:"created_at"`
}

// WebhookRecord maps webhook names to their entries.
type WebhookRecord struct {
	Hooks map[string]WebhookEntry `json:"hooks"`
}

func webhookFilePath(vaultPath string) string {
	return filepath.Join(filepath.Dir(vaultPath), ".webhooks.json")
}

// LoadWebhookRecord loads the webhook registry from disk.
// Returns an empty record if the file does not exist.
func LoadWebhookRecord(vaultPath string) (WebhookRecord, error) {
	path := webhookFilePath(vaultPath)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return WebhookRecord{Hooks: make(map[string]WebhookEntry)}, nil
	}
	if err != nil {
		return WebhookRecord{}, fmt.Errorf("read webhook record: %w", err)
	}
	var rec WebhookRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return WebhookRecord{}, fmt.Errorf("parse webhook record: %w", err)
	}
	if rec.Hooks == nil {
		rec.Hooks = make(map[string]WebhookEntry)
	}
	return rec, nil
}

func saveWebhookRecord(vaultPath string, rec WebhookRecord) error {
	path := webhookFilePath(vaultPath)
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal webhook record: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// RegisterWebhook adds or updates a named webhook with the given URL and events.
func RegisterWebhook(vaultPath, name, url string, events []WebhookEvent) error {
	if name == "" {
		return fmt.Errorf("webhook name must not be empty")
	}
	if url == "" {
		return fmt.Errorf("webhook URL must not be empty")
	}
	rec, err := LoadWebhookRecord(vaultPath)
	if err != nil {
		return err
	}
	sorted := make([]WebhookEvent, len(events))
	copy(sorted, events)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	rec.Hooks[name] = WebhookEntry{URL: url, Events: sorted, CreatedAt: time.Now().UTC()}
	return saveWebhookRecord(vaultPath, rec)
}

// DeregisterWebhook removes a named webhook. Returns an error if not found.
func DeregisterWebhook(vaultPath, name string) error {
	rec, err := LoadWebhookRecord(vaultPath)
	if err != nil {
		return err
	}
	if _, ok := rec.Hooks[name]; !ok {
		return fmt.Errorf("webhook %q not found", name)
	}
	delete(rec.Hooks, name)
	return saveWebhookRecord(vaultPath, rec)
}

// ListWebhooks returns all webhook names in sorted order.
func ListWebhooks(vaultPath string) ([]string, WebhookRecord, error) {
	rec, err := LoadWebhookRecord(vaultPath)
	if err != nil {
		return nil, rec, err
	}
	names := make([]string, 0, len(rec.Hooks))
	for k := range rec.Hooks {
		names = append(names, k)
	}
	sort.Strings(names)
	return names, rec, nil
}
