package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envvault/internal/vault"
)

func init() {
	labelCmd := &cobra.Command{
		Use:   "label",
		Short: "Manage human-readable labels for vault keys",
	}

	setCmd := &cobra.Command{
		Use:   "set <vault> <key> <label>",
		Short: "Assign a label to a vault key",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLabelSet(args[0], args[1], args[2])
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <vault> <key>",
		Short: "Remove the label from a vault key",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLabelRemove(args[0], args[1])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list <vault>",
		Short: "List all labels in a vault",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLabelList(args[0])
		},
	}

	labelCmd.AddCommand(setCmd, removeCmd, listCmd)
	rootCmd.AddCommand(labelCmd)
}

func runLabelSet(vaultPath, key, label string) error {
	if err := vault.SetLabel(vaultPath, key, label); err != nil {
		return fmt.Errorf("set label: %w", err)
	}
	fmt.Fprintf(os.Stdout, "label set for %q\n", key)
	return nil
}

func runLabelRemove(vaultPath, key string) error {
	if err := vault.RemoveLabel(vaultPath, key); err != nil {
		return fmt.Errorf("remove label: %w", err)
	}
	fmt.Fprintf(os.Stdout, "label removed for %q\n", key)
	return nil
}

func runLabelList(vaultPath string) error {
	keys, labels, err := vault.ListLabels(vaultPath)
	if err != nil {
		return fmt.Errorf("list labels: %w", err)
	}
	if len(keys) == 0 {
		fmt.Fprintln(os.Stdout, "no labels defined")
		return nil
	}
	fmt.Fprintf(os.Stdout, "%-30s  %s\n", "KEY", "LABEL")
	fmt.Fprintf(os.Stdout, "%-30s  %s\n", "---", "-----")
	for i, k := range keys {
		fmt.Fprintf(os.Stdout, "%-30s  %s\n", k, labels[i])
	}
	return nil
}
