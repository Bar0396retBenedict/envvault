package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"envvault/internal/vault"
)

func init() {
	var strategy string

	cmd := &cobra.Command{
		Use:   "merge <src-vault> <dst-vault>",
		Short: "Merge keys from a source vault into a destination vault",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMerge(args[0], args[1], strategy)
		},
	}
	cmd.Flags().StringVarP(&strategy, "strategy", "s", "ours",
		"conflict resolution strategy: ours | theirs | error")
	rootCmd.AddCommand(cmd)
}

func runMerge(srcPath, dstPath, strategyFlag string) error {
	srcPass := os.Getenv("ENVVAULT_SRC_PASSPHRASE")
	if srcPass == "" {
		srcPass = os.Getenv("ENVVAULT_PASSPHRASE")
	}
	dstPass := os.Getenv("ENVVAULT_DST_PASSPHRASE")
	if dstPass == "" {
		dstPass = os.Getenv("ENVVAULT_PASSPHRASE")
	}
	if srcPass == "" || dstPass == "" {
		return fmt.Errorf("passphrase not set: use ENVVAULT_PASSPHRASE or ENVVAULT_SRC_PASSPHRASE / ENVVAULT_DST_PASSPHRASE")
	}

	src, err := vault.New(srcPath, srcPass)
	if err != nil {
		return fmt.Errorf("open src vault: %w", err)
	}
	dst, err := vault.New(dstPath, dstPass)
	if err != nil {
		return fmt.Errorf("open dst vault: %w", err)
	}

	var strategy vault.MergeStrategy
	switch strategyFlag {
	case "ours":
		strategy = vault.MergeStrategyOurs
	case "theirs":
		strategy = vault.MergeStrategyTheirs
	case "error":
		strategy = vault.MergeStrategyError
	default:
		return fmt.Errorf("unknown strategy %q: choose ours, theirs, or error", strategyFlag)
	}

	result, err := vault.MergeInto(dst, src, strategy)
	if err != nil {
		return fmt.Errorf("merge failed: %w", err)
	}

	if err := dst.Save(); err != nil {
		return fmt.Errorf("save dst vault: %w", err)
	}

	sort.Strings(result.Added)
	sort.Strings(result.Updated)
	sort.Strings(result.Skipped)

	for _, k := range result.Added {
		fmt.Printf("+ %s\n", k)
	}
	for _, k := range result.Updated {
		fmt.Printf("~ %s\n", k)
	}
	for _, k := range result.Skipped {
		fmt.Printf("= %s\n", k)
	}
	fmt.Printf("\nmerge complete: %d added, %d updated, %d skipped\n",
		len(result.Added), len(result.Updated), len(result.Skipped))
	return nil
}
