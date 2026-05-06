package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envvault/internal/vault"
)

func init() {
	groupCmd := &cobra.Command{
		Use:   "group",
		Short: "Manage key groups within a vault",
	}

	addCmd := &cobra.Command{
		Use:   "add <vault> <group> <key>",
		Short: "Add a key to a group",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGroupAdd(args[0], args[1], args[2])
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <vault> <group> <key>",
		Short: "Remove a key from a group",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGroupRemove(args[0], args[1], args[2])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list <vault> [group]",
		Short: "List groups or keys in a group",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 2 {
				return runGroupKeys(args[0], args[1])
			}
			return runGroupList(args[0])
		},
	}

	groupCmd.AddCommand(addCmd, removeCmd, listCmd)
	rootCmd.AddCommand(groupCmd)
}

func runGroupAdd(vaultPath, group, key string) error {
	if err := vault.AddToGroup(vaultPath, group, key); err != nil {
		return fmt.Errorf("group add: %w", err)
	}
	fmt.Fprintf(os.Stdout, "Added %q to group %q\n", key, group)
	return nil
}

func runGroupRemove(vaultPath, group, key string) error {
	if err := vault.RemoveFromGroup(vaultPath, group, key); err != nil {
		return fmt.Errorf("group remove: %w", err)
	}
	fmt.Fprintf(os.Stdout, "Removed %q from group %q\n", key, group)
	return nil
}

func runGroupList(vaultPath string) error {
	groups, err := vault.ListGroups(vaultPath)
	if err != nil {
		return fmt.Errorf("group list: %w", err)
	}
	if len(groups) == 0 {
		fmt.Fprintln(os.Stdout, "No groups defined.")
		return nil
	}
	for _, g := range groups {
		fmt.Fprintln(os.Stdout, g)
	}
	return nil
}

func runGroupKeys(vaultPath, group string) error {
	keys, err := vault.KeysForGroup(vaultPath, group)
	if err != nil {
		return fmt.Errorf("group keys: %w", err)
	}
	if len(keys) == 0 {
		fmt.Fprintf(os.Stdout, "Group %q is empty.\n", group)
		return nil
	}
	for _, k := range keys {
		fmt.Fprintln(os.Stdout, k)
	}
	return nil
}
