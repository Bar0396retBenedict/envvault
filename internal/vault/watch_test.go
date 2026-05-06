package vault

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeWatchVault(t *testing.T, passphrase, key, value string) (string, *Vault) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.vault")
	v := New(passphrase)
	v.Set(key, value)
	if err := v.SaveToFile(path); err != nil {
		t.Fatalf("save: %v", err)
	}
	return path, v
}

func TestFileHashConsistent(t *testing.T) {
	path, _ := writeWatchVault(t, "pw", "KEY", "val")
	h1, err := fileHash(path)
	if err != nil {
		t.Fatalf("fileHash: %v", err)
	}
	h2, err := fileHash(path)
	if err != nil {
		t.Fatalf("fileHash: %v", err)
	}
	if h1 != h2 {
		t.Errorf("expected stable hash, got %s vs %s", h1, h2)
	}
}

func TestFileHashChangesOnWrite(t *testing.T) {
	path, v := writeWatchVault(t, "pw", "KEY", "val")
	h1, _ := fileHash(path)

	v.Set("KEY", "changed")
	if err := v.SaveToFile(path); err != nil {
		t.Fatalf("save: %v", err)
	}
	h2, _ := fileHash(path)
	if h1 == h2 {
		t.Error("expected hash to change after vault update")
	}
}

func TestFileHashMissingFile(t *testing.T) {
	_, err := fileHash("/nonexistent/path/vault.enc")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestWatchVaultDetectsChange(t *testing.T) {
	path, v := writeWatchVault(t, "pw", "FOO", "bar")

	ch := make(chan WatchEvent, 1)
	done := make(chan struct{})
	defer close(done)

	go WatchVault(path, 20*time.Millisecond, ch, done)

	// Allow watcher to record initial hash.
	time.Sleep(30 * time.Millisecond)

	v.Set("FOO", "updated")
	if err := v.SaveToFile(path); err != nil {
		t.Fatalf("save: %v", err)
	}

	select {
	case evt := <-ch:
		if evt.Path != path {
			t.Errorf("unexpected path %s", evt.Path)
		}
		if evt.OldHash == evt.NewHash {
			t.Error("old and new hash should differ")
		}
		if evt.At.IsZero() {
			t.Error("event timestamp should not be zero")
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("timed out waiting for watch event")
	}
}

func TestWatchVaultNoEventWhenUnchanged(t *testing.T) {
	path, _ := writeWatchVault(t, "pw", "STABLE", "value")

	ch := make(chan WatchEvent, 1)
	done := make(chan struct{})
	defer close(done)

	go WatchVault(path, 20*time.Millisecond, ch, done)

	select {
	case evt := <-ch:
		t.Errorf("unexpected event for unchanged file: %+v", evt)
	case <-time.After(150 * time.Millisecond):
		// expected: no event
	}
}

func TestWatchVaultMissingFileNoBlock(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "missing.vault")

	ch := make(chan WatchEvent, 1)
	done := make(chan struct{})
	defer close(done)

	go WatchVault(path, 20*time.Millisecond, ch, done)

	// Write the file after the watcher has started.
	time.Sleep(40 * time.Millisecond)
	if err := os.WriteFile(path, []byte("data"), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}

	select {
	case <-ch:
		// event received after file appeared — acceptable
	case <-time.After(300 * time.Millisecond):
		// also acceptable: watcher skips errors silently
	}
}
