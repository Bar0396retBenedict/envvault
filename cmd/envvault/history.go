package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"envvault/internal/vault"
)

func init() {
	historyCmd := &cobra.Command{
		Use:   "history <vault-file>",
		Short: "Show key change history for a vault",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHistory(args[0])
		},
	}
	rootCmd.AddCommand(historyCmd)
}

func runHistory(vaultPath string) error {
	rec, err := vault.LoadHistory(vaultPath)
	if err != nil {
		return fmt.Errorf("loading history: %w", err)
	}

	if len(rec.Entries) == 0 {
		fmt.Println("No history recorded.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tACTION\tKEY\tOLD VALUE\tNEW VALUE")
	fmt.Fprintln(w, "---------\t------\t---\t---------\t---------")
	for _, e := range rec.Entries {
		oldVal := e.OldValue
		if oldVal == "" {
			oldVal = "-"
		}
		newVal := e.NewValue
		if newVal == "" {
			newVal = "-"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			e.Timestamp.Format("2006-01-02 15:04:05"),
			e.Action,
			e.Key,
			oldVal,
			newVal,
		)
	}
	return w.Flush()
}
