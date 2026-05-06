package main

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"envvault/internal/vault"
)

func init() {
	envCmd := &cobra.Command{
		Use:   "env",
		Short: "Manage OS environment variable bindings for a vault",
	}

	bindCmd := &cobra.Command{
		Use:   "bind <vault> <ENV_VAR> <vault-key>",
		Short: "Bind an OS env var to a vault key",
		Args:  cobra.ExactArgs(3),
		RunE:  runEnvBind,
	}

	unbindCmd := &cobra.Command{
		Use:   "unbind <vault> <ENV_VAR>",
		Short: "Remove an OS env var binding",
		Args:  cobra.ExactArgs(2),
		RunE:  runEnvUnbind,
	}

	listCmd := &cobra.Command{
		Use:   "list <vault>",
		Short: "List all env var bindings",
		Args:  cobra.ExactArgs(1),
		RunE:  runEnvList,
	}

	envCmd.AddCommand(bindCmd, unbindCmd, listCmd)
	rootCmd.AddCommand(envCmd)
}

func runEnvBind(cmd *cobra.Command, args []string) error {
	vaultPath, envKey, vaultKey := args[0], args[1], args[2]
	if err := vault.BindEnvVar(vaultPath, envKey, vaultKey); err != nil {
		return fmt.Errorf("bind: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "bound %s -> %s\n", envKey, vaultKey)
	return nil
}

func runEnvUnbind(cmd *cobra.Command, args []string) error {
	vaultPath, envKey := args[0], args[1]
	if err := vault.UnbindEnvVar(vaultPath, envKey); err != nil {
		return fmt.Errorf("unbind: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "unbound %s\n", envKey)
	return nil
}

func runEnvList(cmd *cobra.Command, args []string) error {
	vaultPath := args[0]
	rec, err := vault.LoadEnvRecord(vaultPath)
	if err != nil {
		return fmt.Errorf("list: %w", err)
	}
	if len(rec.Bindings) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no env bindings configured")
		return nil
	}
	keys := make([]string, 0, len(rec.Bindings))
	for k := range rec.Bindings {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ENV VAR\tVAULT KEY")
	for _, k := range keys {
		fmt.Fprintf(w, "%s\t%s\n", k, rec.Bindings[k])
	}
	_ = w.Flush()
	return nil
}

var _ = os.Stderr // suppress unused import
