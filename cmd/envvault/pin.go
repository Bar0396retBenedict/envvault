package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envvault/internal/vault"
)

func init() {
	pinCmd := &cobra.Command{
		Use:   "pin",
		Short: "Pin or unpin keys to prevent accidental modification",
	}

	addCmd := &cobra.Command{
		Use:   "add <vault> <key>",
		Short: "Pin a key",
		Args:  cobra.ExactArgs(2),
		RunE:  runPinAdd,
	}
	addCmd.Flags().String("note", "", "Optional note describing why the key is pinned")

	removeCmd := &cobra.Command{
		Use:   "remove <vault> <key>",
		Short: "Unpin a key",
		Args:  cobra.ExactArgs(2),
		RunE:  runPinRemove,
	}

	listCmd := &cobra.Command{
		Use:   "list <vault>",
		Short: "List pinned keys",
		Args:  cobra.ExactArgs(1),
		RunE:  runPinList,
	}

	pinCmd.AddCommand(addCmd, removeCmd, listCmd)
	rootCmd.AddCommand(pinCmd)
}

func runPinAdd(cmd *cobra.Command, args []string) error {
	vaultPath, key := args[0], args[1]
	note, _ := cmd.Flags().GetString("note")
	if err := vault.PinKey(vaultPath, key, note); err != nil {
		return fmt.Errorf("pin key: %w", err)
	}
	fmt.Fprintf(os.Stdout, "pinned %q in %s\n", key, vaultPath)
	return nil
}

func runPinRemove(cmd *cobra.Command, args []string) error {
	vaultPath, key := args[0], args[1]
	if err := vault.UnpinKey(vaultPath, key); err != nil {
		return fmt.Errorf("unpin key: %w", err)
	}
	fmt.Fprintf(os.Stdout, "unpinned %q from %s\n", key, vaultPath)
	return nil
}

func runPinList(cmd *cobra.Command, args []string) error {
	vaultPath := args[0]
	keys, rec, err := vault.ListPinnedKeys(vaultPath)
	if err != nil {
		return fmt.Errorf("list pins: %w", err)
	}
	if len(keys) == 0 {
		fmt.Fprintln(os.Stdout, "no pinned keys")
		return nil
	}
	for _, k := range keys {
		entry := rec.Pins[k]
		if entry.Note != "" {
			fmt.Fprintf(os.Stdout, "%-30s  pinned at %s  note: %s\n", k, entry.PinnedAt.Format("2006-01-02 15:04:05"), entry.Note)
		} else {
			fmt.Fprintf(os.Stdout, "%-30s  pinned at %s\n", k, entry.PinnedAt.Format("2006-01-02 15:04:05"))
		}
	}
	return nil
}
