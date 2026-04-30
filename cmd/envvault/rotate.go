package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envvault/internal/vault"
)

func runRotate(args []string) error {
	fs := flag.NewFlagSet("rotate", flag.ContinueOnError)
	oldPass := fs.String("old-passphrase", os.Getenv("ENVVAULT_OLD_PASSPHRASE"), "current passphrase (env: ENVVAULT_OLD_PASSPHRASE)")
	newPass := fs.String("new-passphrase", os.Getenv("ENVVAULT_NEW_PASSPHRASE"), "replacement passphrase (env: ENVVAULT_NEW_PASSPHRASE)")
	vaultPath := fs.String("vault", "envvault.enc", "path to the vault file")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *oldPass == "" {
		return fmt.Errorf("rotate: --old-passphrase is required")
	}
	if *newPass == "" {
		return fmt.Errorf("rotate: --new-passphrase is required")
	}

	rec, err := vault.Rotate(*vaultPath, *oldPass, *newPass)
	if err != nil {
		return fmt.Errorf("rotate failed: %w", err)
	}

	fmt.Printf("✓ Vault re-encrypted successfully\n")
	fmt.Printf("  Profile : %s\n", rec.Profile)
	fmt.Printf("  Rotated : %s\n", rec.RotatedAt.Format("2006-01-02 15:04:05 UTC"))
	return nil
}

func init() {
	registerCommand("rotate", "Re-encrypt the vault with a new passphrase", runRotate)
}
