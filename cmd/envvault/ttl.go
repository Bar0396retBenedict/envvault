package main

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"envvault/internal/vault"
)

func init() {
	ttlCmd := &cobra.Command{
		Use:   "ttl",
		Short: "Manage key expiry (time-to-live) within a vault",
	}

	setCmd := &cobra.Command{
		Use:   "set <vault> <key> <duration>",
		Short: "Set a TTL on a vault key (e.g. 24h, 30m)",
		Args:  cobra.ExactArgs(3),
		RunE:  runTTLSet,
	}

	listCmd := &cobra.Command{
		Use:   "list <vault>",
		Short: "List all TTL entries for a vault",
		Args:  cobra.ExactArgs(1),
		RunE:  runTTLList,
	}

	purgeCmd := &cobra.Command{
		Use:   "purge <vault> <key>",
		Short: "Remove the TTL entry for a vault key",
		Args:  cobra.ExactArgs(2),
		RunE:  runTTLPurge,
	}

	ttlCmd.AddCommand(setCmd, listCmd, purgeCmd)
	rootCmd.AddCommand(ttlCmd)
}

func runTTLSet(cmd *cobra.Command, args []string) error {
	vaultPath, key, rawDur := args[0], args[1], args[2]
	ttl, err := time.ParseDuration(rawDur)
	if err != nil {
		return fmt.Errorf("invalid duration %q: %w", rawDur, err)
	}
	if err := vault.SetTTL(vaultPath, key, ttl); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "TTL set: %s expires in %s\n", key, ttl)
	return nil
}

func runTTLList(cmd *cobra.Command, args []string) error {
	vaultPath := args[0]
	rec, err := vault.LoadTTLRecord(vaultPath)
	if err != nil {
		return err
	}
	if len(rec.Entries) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No TTL entries.")
		return nil
	}
	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KEY\tEXPIRES AT\tSTATUS")
	now := time.Now().UTC()
	for _, e := range rec.Entries {
		status := "valid"
		if now.After(e.ExpiresAt) {
			status = "EXPIRED"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", e.Key, e.ExpiresAt.Format(time.RFC3339), status)
	}
	w.Flush()
	return nil
}

func runTTLPurge(cmd *cobra.Command, args []string) error {
	vaultPath, key := args[0], args[1]
	if err := vault.PurgeTTLEntry(vaultPath, key); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "TTL entry removed for key: %s\n", key)
	_ = os.Stderr // suppress unused import warning
	return nil
}
