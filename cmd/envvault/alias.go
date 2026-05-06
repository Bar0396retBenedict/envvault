package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"envvault/internal/vault"
)

func init() {
	aliasCmd := &cobra.Command{
		Use:   "alias",
		Short: "Manage key aliases",
	}

	setCmd := &cobra.Command{
		Use:   "set <alias> <key>",
		Short: "Create or update an alias pointing to a vault key",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			vaultPath, _ := cmd.Flags().GetString("vault")
			return runAliasSet(vaultPath, args[0], args[1])
		},
	}
	setCmd.Flags().String("vault", "envvault.vault", "path to vault file")

	removeCmd := &cobra.Command{
		Use:   "remove <alias>",
		Short: "Remove an alias",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			vaultPath, _ := cmd.Flags().GetString("vault")
			return runAliasRemove(vaultPath, args[0])
		},
	}
	removeCmd.Flags().String("vault", "envvault.vault", "path to vault file")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all aliases",
		RunE: func(cmd *cobra.Command, args []string) error {
			vaultPath, _ := cmd.Flags().GetString("vault")
			return runAliasList(vaultPath)
		},
	}
	listCmd.Flags().String("vault", "envvault.vault", "path to vault file")

	aliasCmd.AddCommand(setCmd, removeCmd, listCmd)
	rootCmd.AddCommand(aliasCmd)
}

func runAliasSet(vaultPath, alias, key string) error {
	if err := vault.SetAlias(vaultPath, alias, key); err != nil {
		return err
	}
	fmt.Printf("alias %q → %q saved\n", alias, key)
	return nil
}

func runAliasRemove(vaultPath, alias string) error {
	if err := vault.RemoveAlias(vaultPath, alias); err != nil {
		return err
	}
	fmt.Printf("alias %q removed\n", alias)
	return nil
}

func runAliasList(vaultPath string) error {
	pairs, err := vault.ListAliases(vaultPath)
	if err != nil {
		return err
	}
	if len(pairs) == 0 {
		fmt.Println("no aliases defined")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ALIAS\tKEY")
	for _, p := range pairs {
		fmt.Fprintf(w, "%s\t%s\n", p[0], p[1])
	}
	return w.Flush()
}
