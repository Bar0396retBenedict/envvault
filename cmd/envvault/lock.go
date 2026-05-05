package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envvault/internal/vault"
)

func init() {
	lockCmd := &cobra.Command{
		Use:   "lock <vault-file>",
		Short: "Lock a vault to prevent modifications",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			by, _ := cmd.Flags().GetString("by")
			return runLock(args[0], by)
		},
	}
	lockCmd.Flags().String("by", "", "identifier recorded as the lock owner (e.g. username)")

	unlockCmd := &cobra.Command{
		Use:   "unlock <vault-file>",
		Short: "Unlock a previously locked vault",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUnlock(args[0])
		},
	}

	rootCmd.AddCommand(lockCmd)
	rootCmd.AddCommand(unlockCmd)
}

func runLock(vaultPath, lockedBy string) error {
	if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
		return fmt.Errorf("vault file not found: %s", vaultPath)
	}
	if err := vault.LockVault(vaultPath, lockedBy); err != nil {
		return fmt.Errorf("lock: %w", err)
	}
	fmt.Printf("Vault locked: %s\n", vaultPath)
	if lockedBy != "" {
		fmt.Printf("Locked by: %s\n", lockedBy)
	}
	return nil
}

func runUnlock(vaultPath string) error {
	if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
		return fmt.Errorf("vault file not found: %s", vaultPath)
	}
	rec, err := vault.LoadLockRecord(vaultPath)
	if err != nil {
		return fmt.Errorf("unlock: %w", err)
	}
	if err := vault.UnlockVault(vaultPath); err != nil {
		return fmt.Errorf("unlock: %w", err)
	}
	fmt.Printf("Vault unlocked: %s\n", vaultPath)
	if rec.LockedBy != "" {
		fmt.Printf("Was locked by: %s\n", rec.LockedBy)
	}
	return nil
}
