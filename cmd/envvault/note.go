package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envvault/internal/vault"
)

func init() {
	noteCmd := &cobra.Command{
		Use:   "note",
		Short: "Manage annotations attached to vault keys",
	}

	setCmd := &cobra.Command{
		Use:   "set <vault> <key> <note>",
		Short: "Attach or replace a note for a key",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runNoteSet(args[0], args[1], args[2])
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <vault> <key>",
		Short: "Remove the note for a key",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runNoteRemove(args[0], args[1])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list <vault>",
		Short: "List all notes in a vault",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runNoteList(args[0])
		},
	}

	noteCmd.AddCommand(setCmd, removeCmd, listCmd)
	rootCmd.AddCommand(noteCmd)
}

func runNoteSet(vaultPath, key, note string) error {
	if err := vault.SetNote(vaultPath, key, note); err != nil {
		return fmt.Errorf("set note: %w", err)
	}
	fmt.Fprintf(os.Stdout, "note set for key %q\n", key)
	return nil
}

func runNoteRemove(vaultPath, key string) error {
	if err := vault.RemoveNote(vaultPath, key); err != nil {
		return fmt.Errorf("remove note: %w", err)
	}
	fmt.Fprintf(os.Stdout, "note removed for key %q\n", key)
	return nil
}

func runNoteList(vaultPath string) error {
	rec, err := vault.LoadNoteRecord(vaultPath)
	if err != nil {
		return fmt.Errorf("load notes: %w", err)
	}
	if len(rec.Notes) == 0 {
		fmt.Fprintln(os.Stdout, "no notes found")
		return nil
	}
	fmt.Fprintf(os.Stdout, "%-30s  %-20s  %s\n", "KEY", "UPDATED", "NOTE")
	for _, e := range rec.Notes {
		fmt.Fprintf(os.Stdout, "%-30s  %-20s  %s\n",
			e.Key,
			e.UpdatedAt.Format("2006-01-02 15:04:05"),
			e.Note,
		)
	}
	return nil
}
