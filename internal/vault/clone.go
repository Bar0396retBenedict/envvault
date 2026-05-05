package vault

import (
	"fmt"
	"os"
)

// CloneOptions controls the behaviour of CloneVault.
type CloneOptions struct {
	// Overwrite allows the destination file to be replaced if it already exists.
	Overwrite bool
}

// CloneVault reads the vault at srcPath (decrypted with srcPassphrase), creates
// a new vault at dstPath encrypted with dstPassphrase, and copies every key
// from the source into it. The destination vault is saved atomically; if
// dstPath already exists and Overwrite is false an error is returned.
//
// CloneVault is useful when promoting a local vault to staging/production with
// a different passphrase, or when creating a sanitised copy for a new team
// member.
func CloneVault(srcPath, srcPassphrase, dstPath, dstPassphrase string, opts CloneOptions) error {
	if srcPassphrase == "" {
		return fmt.Errorf("clone: source passphrase must not be empty")
	}
	if dstPassphrase == "" {
		return fmt.Errorf("clone: destination passphrase must not be empty")
	}

	src, err := Load(srcPath, srcPassphrase)
	if err != nil {
		return fmt.Errorf("clone: load source vault: %w", err)
	}

	if !opts.Overwrite {
		if _, err := os.Stat(dstPath); err == nil {
			return fmt.Errorf("clone: destination %q already exists (use Overwrite to replace)", dstPath)
		}
	}

	dst := New()

	for _, key := range src.Keys() {
		val, ok := src.Get(key)
		if !ok {
			continue
		}
		if err := dst.Set(key, val); err != nil {
			return fmt.Errorf("clone: set key %q: %w", key, err)
		}
	}

	if err := dst.Save(dstPath, dstPassphrase); err != nil {
		return fmt.Errorf("clone: save destination vault: %w", err)
	}

	return nil
}
