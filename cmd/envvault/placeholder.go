package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"envvault/internal/vault"
)

func init() {
	placeholderCmd := &cobra.Command{
		Use:   "placeholder",
		Short: "Manage key placeholder descriptions",
	}

	setCmd := &cobra.Command{
		Use:   "set <vault> <key> <description>",
		Short: "Set a placeholder description for a vault key",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPlaceholderSet(args[0], args[1], args[2])
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <vault> <key>",
		Short: "Remove a placeholder description for a vault key",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPlaceholderRemove(args[0], args[1])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list <vault>",
		Short: "List all placeholder descriptions in a vault",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPlaceholderList(args[0])
		},
	}

	placeholderCmd.AddCommand(setCmd, removeCmd, listCmd)
	rootCmd.AddCommand(placeholderCmd)
}

func runPlaceholderSet(vaultPath, key, description string) error {
	if err := vault.SetPlaceholder(vaultPath, key, description); err != nil {
		return fmt.Errorf("set placeholder: %w", err)
	}
	fmt.Printf("placeholder set for %q\n", key)
	return nil
}

func runPlaceholderRemove(vaultPath, key string) error {
	if err := vault.RemovePlaceholder(vaultPath, key); err != nil {
		return fmt.Errorf("remove placeholder: %w", err)
	}
	fmt.Printf("placeholder removed for %q\n", key)
	return nil
}

func runPlaceholderList(vaultPath string) error {
	entries, err := vault.ListPlaceholders(vaultPath)
	if err != nil {
		return fmt.Errorf("list placeholders: %w", err)
	}
	if len(entries) == 0 {
		fmt.Println("no placeholders defined")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KEY\tDESCRIPTION")
	for _, e := range entries {
		fmt.Fprintf(w, "%s\t%s\n", e.Key, e.Description)
	}
	return w.Flush()
}
