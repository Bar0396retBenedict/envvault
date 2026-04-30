package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envvault/internal/vault"
)

func init() {
	snapshotCmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Manage vault snapshots",
	}

	takeCmd := &cobra.Command{
		Use:   "take <vault-file> <label>",
		Short: "Take a snapshot of the vault",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSnapshotTake(args[0], args[1])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list <vault-file>",
		Short: "List all snapshots for a vault",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSnapshotList(args[0])
		},
	}

	snapshotCmd.AddCommand(takeCmd, listCmd)
	rootCmd.AddCommand(snapshotCmd)
}

func runSnapshotTake(vaultPath, label string) error {
	pass := os.Getenv("ENVVAULT_PASSPHRASE")
	if pass == "" {
		return fmt.Errorf("ENVVAULT_PASSPHRASE environment variable is not set")
	}
	snap, err := vault.TakeSnapshot(vaultPath, pass, label)
	if err != nil {
		return err
	}
	fmt.Printf("Snapshot '%s' taken at %s (%d entries)\n",
		snap.Label, snap.CreatedAt.Format("2006-01-02 15:04:05 UTC"), len(snap.Entries))
	return nil
}

func runSnapshotList(vaultPath string) error {
	snaps, err := vault.ListSnapshots(vaultPath)
	if err != nil {
		return err
	}
	if len(snaps) == 0 {
		fmt.Println("No snapshots found.")
		return nil
	}
	fmt.Printf("%-30s  %-20s  %s\n", "TIMESTAMP", "LABEL", "KEYS")
	for _, s := range snaps {
		fmt.Printf("%-30s  %-20s  %d\n",
			s.CreatedAt.Format("2006-01-02 15:04:05 UTC"),
			s.Label,
			len(s.Entries),
		)
	}
	return nil
}
