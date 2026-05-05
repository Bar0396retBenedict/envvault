package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envvault/internal/vault"
)

func init() {
	var vaultFile string
	var passphrase string
	var matchValue bool
	var caseSensitive bool

	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search for keys (and optionally values) in a vault",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSearch(vaultFile, passphrase, args[0], matchValue, caseSensitive)
		},
	}

	cmd.Flags().StringVarP(&vaultFile, "file", "f", "vault.env", "path to vault file")
	cmd.Flags().StringVarP(&passphrase, "passphrase", "p", "", "vault passphrase")
	cmd.Flags().BoolVar(&matchValue, "match-value", false, "also search within values")
	cmd.Flags().BoolVar(&caseSensitive, "case-sensitive", false, "use case-sensitive matching")

	rootCmd.AddCommand(cmd)
}

func runSearch(vaultFile, passphrase, query string, matchValue, caseSensitive bool) error {
	if passphrase == "" {
		return fmt.Errorf("passphrase is required (use -p)")
	}

	v, err := vault.Load(vaultFile, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load vault: %w", err)
	}

	opts := vault.SearchOptions{
		MatchValue:    matchValue,
		CaseSensitive: caseSensitive,
	}

	results := vault.Search(v, query, opts)
	if len(results) == 0 {
		fmt.Fprintln(os.Stderr, "no matching keys found")
		return nil
	}

	for _, r := range results {
		fmt.Printf("%-40s %s\n", r.Key, r.Value)
	}
	return nil
}
