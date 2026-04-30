package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envvault/internal/vault"
)

var (
	importFormat   string
	importOverwrite bool
)

func init() {
	importCmd := &cobra.Command{
		Use:   "import <vault-file> <env-file>",
		Short: "Import key=value pairs from a .env or shell file into a vault",
		Args:  cobra.ExactArgs(2),
		RunE:  runImport,
	}
	importCmd.Flags().StringVarP(&importFormat, "format", "f", "dotenv", "Input format: dotenv or shell")
	importCmd.Flags().BoolVar(&importOverwrite, "overwrite", false, "Overwrite existing keys")
	rootCmd.AddCommand(importCmd)
}

func runImport(cmd *cobra.Command, args []string) error {
	vaultPath := args[0]
	envPath := args[1]

	passphrase := os.Getenv("ENVVAULT_PASSPHRASE")
	if passphrase == "" {
		return errors.New("ENVVAULT_PASSPHRASE environment variable is required")
	}

	v, err := vault.Load(vaultPath, passphrase)
	if err != nil {
		return fmt.Errorf("load vault: %w", err)
	}

	var fmt_ vault.ImportFormat
	switch importFormat {
	case "dotenv", "":
		fmt_ = vault.ImportDotEnv
	case "shell":
		fmt_ = vault.ImportShell
	default:
		return fmt.Errorf("unknown format %q: use dotenv or shell", importFormat)
	}

	res, err := vault.ImportFromFile(v, envPath, fmt_, importOverwrite)
	if err != nil {
		return fmt.Errorf("import: %w", err)
	}

	if err := v.Save(vaultPath, passphrase); err != nil {
		return fmt.Errorf("save vault: %w", err)
	}

	cmd.Printf("Imported: %d added, %d overwritten, %d skipped\n",
		res.Added, res.Overwritten, res.Skipped)
	return nil
}
