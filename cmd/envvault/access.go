package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"envvault/internal/vault"
)

func init() {
	accessCmd := &cobra.Command{
		Use:   "access <vault-file>",
		Short: "Show key access statistics for a vault",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAccess(args[0])
		},
	}
	rootCmd.AddCommand(accessCmd)
}

func runAccess(vaultPath string) error {
	entries, err := vault.ListAccessEntries(vaultPath)
	if err != nil {
		return fmt.Errorf("access: %w", err)
	}
	if len(entries) == 0 {
		fmt.Println("No access records found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KEY\tREADS\tWRITES\tLAST READ\tLAST WRITE")
	for _, e := range entries {
		lr := "-"
		if !e.LastRead.IsZero() {
			lr = e.LastRead.Format("2006-01-02 15:04:05")
		}
		lw := "-"
		if !e.LastWrite.IsZero() {
			lw = e.LastWrite.Format("2006-01-02 15:04:05")
		}
		fmt.Fprintf(w, "%s\t%d\t%d\t%s\t%s\n",
			e.Key, e.ReadCount, e.WriteCount, lr, lw)
	}
	return w.Flush()
}
