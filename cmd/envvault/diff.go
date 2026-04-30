package main

import (
	"fmt"
	"os"

	"github.com/user/envvault/internal/vault"
)

// runDiff loads two vault files and prints a human-readable diff to stdout.
// Usage: envvault diff <src-vault> <dst-vault>
func runDiff(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envvault diff <src-vault> <dst-vault>")
	}

	passphrase := os.Getenv("ENVVAULT_PASSPHRASE")
	if passphrase == "" {
		return fmt.Errorf("ENVVAULT_PASSPHRASE environment variable is not set")
	}

	src, err := vault.Load(args[0], passphrase)
	if err != nil {
		return fmt.Errorf("loading src vault %q: %w", args[0], err)
	}

	dst, err := vault.Load(args[1], passphrase)
	if err != nil {
		return fmt.Errorf("loading dst vault %q: %w", args[1], err)
	}

	changes := vault.Diff(src, dst)
	if len(changes) == 0 {
		fmt.Println("No differences found.")
		return nil
	}

	fmt.Printf("Comparing %s → %s\n\n", args[0], args[1])
	for _, c := range changes {
		switch c.Kind {
		case vault.ChangeAdded:
			fmt.Printf("  + %-30s = %s\n", c.Key, c.NewValue)
		case vault.ChangeRemoved:
			fmt.Printf("  - %-30s   (was %s)\n", c.Key, c.OldValue)
		case vault.ChangeUpdated:
			fmt.Printf("  ~ %-30s : %s → %s\n", c.Key, c.OldValue, c.NewValue)
		}
	}
	fmt.Printf("\n%d change(s) total.\n", len(changes))
	return nil
}
