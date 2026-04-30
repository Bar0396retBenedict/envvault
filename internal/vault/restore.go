package vault

import (
	"fmt"
	"os"
	"path/filepath"
)

// Restore loads a snapshot by label and writes it back to the target vault path,
// re-encrypting with the provided passphrase. If the snapshot was created with a
// different passphrase, snapshotPassphrase must match that original passphrase.
func Restore(vaultPath, label, snapshotPassphrase, targetPassphrase string) error {
	if label == "" {
		return fmt.Errorf("restore: label must not be empty")
	}
	if snapshotPassphrase == "" || targetPassphrase == "" {
		return fmt.Errorf("restore: passphrases must not be empty")
	}

	dir := snapshotDir(vaultPath)
	safe := sanitizeLabel(label)
	snapshotPath := filepath.Join(dir, safe+".vault")

	if _, err := os.Stat(snapshotPath); os.IsNotExist(err) {
		return fmt.Errorf("restore: snapshot %q not found", label)
	}

	src, err := Load(snapshotPath, snapshotPassphrase)
	if err != nil {
		return fmt.Errorf("restore: failed to load snapshot: %w", err)
	}

	// Re-create a vault at the target path with the target passphrase.
	dst := New(targetPassphrase)
	for _, k := range src.Keys() {
		v, _ := src.Get(k)
		dst.Set(k, v)
	}

	if err := dst.Save(vaultPath); err != nil {
		return fmt.Errorf("restore: failed to save vault: %w", err)
	}
	return nil
}
