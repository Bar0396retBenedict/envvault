package vault

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"time"
)

// WatchEvent describes a change detected in a vault file.
type WatchEvent struct {
	Path    string
	OldHash string
	NewHash string
	At      time.Time
}

// fileHash returns the SHA-256 hex digest of the file at path.
func fileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("watch: open %s: %w", path, err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("watch: hash %s: %w", path, err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// WatchVault polls the vault file at path every interval and sends a
// WatchEvent on ch whenever the file content changes. The goroutine
// exits when done is closed. Any I/O errors are silently skipped so
// that transient filesystem issues do not abort a long-running watch.
func WatchVault(path string, interval time.Duration, ch chan<- WatchEvent, done <-chan struct{}) {
	last, _ := fileHash(path)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			current, err := fileHash(path)
			if err != nil {
				continue
			}
			if current != last {
				ch <- WatchEvent{
					Path:    path,
					OldHash: last,
					NewHash: current,
					At:      time.Now().UTC(),
				}
				last = current
			}
		}
	}
}
