package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envvault/internal/vault"
)

func init() {
	commentCmd := &cobra.Command{
		Use:   "comment",
		Short: "Manage comments on vault keys",
	}

	setCmd := &cobra.Command{
		Use:   "set <vault> <key> <comment>",
		Short: "Set a comment on a vault key",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommentSet(args[0], args[1], args[2])
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <vault> <key>",
		Short: "Remove the comment from a vault key",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommentRemove(args[0], args[1])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list <vault>",
		Short: "List all comments in a vault",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommentList(args[0])
		},
	}

	commentCmd.AddCommand(setCmd, removeCmd, listCmd)
	rootCmd.AddCommand(commentCmd)
}

func runCommentSet(vaultPath, key, comment string) error {
	if err := vault.SetComment(vaultPath, key, comment); err != nil {
		return fmt.Errorf("set comment: %w", err)
	}
	fmt.Fprintf(os.Stdout, "comment set for key %q\n", key)
	return nil
}

func runCommentRemove(vaultPath, key string) error {
	if err := vault.RemoveComment(vaultPath, key); err != nil {
		return fmt.Errorf("remove comment: %w", err)
	}
	fmt.Fprintf(os.Stdout, "comment removed for key %q\n", key)
	return nil
}

func runCommentList(vaultPath string) error {
	keys, comments, err := vault.ListComments(vaultPath)
	if err != nil {
		return fmt.Errorf("list comments: %w", err)
	}
	if len(keys) == 0 {
		fmt.Fprintln(os.Stdout, "no comments found")
		return nil
	}
	fmt.Fprintf(os.Stdout, "%-30s %s\n", "KEY", "COMMENT")
	fmt.Fprintf(os.Stdout, "%-30s %s\n", "---", "-------")
	for i, k := range keys {
		fmt.Fprintf(os.Stdout, "%-30s %s\n", k, comments[i])
	}
	return nil
}
