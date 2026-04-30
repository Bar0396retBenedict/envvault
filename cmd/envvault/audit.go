package main

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"envvault/internal/vault"
)

func init() {
	auditCmd := &cobra.Command{
		Use:   "audit <vault-file>",
		Short: "Show the audit log for a vault file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAudit(args[0])
		},
	}
	rootCmd.AddCommand(auditCmd)
}

func runAudit(vaultPath string) error {
	log, err := vault.LoadAuditLog(vaultPath)
	if err != nil {
		return fmt.Errorf("failed to load audit log: %w", err)
	}
	if len(log.Entries) == 0 {
		fmt.Fprintln(os.Stdout, "No audit entries found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tACTION\tKEY\tVAULT")
	for _, e := range log.Entries {
		key := e.Key
		if key == "" {
			key = "-"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			e.Timestamp.Format(time.RFC3339),
			string(e.Action),
			key,
			e.VaultPath,
		)
	}
	return w.Flush()
}
