package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envvault/internal/vault"
)

var exportCmd = &cobra.Command{
	Use:   "export [vault-file]",
	Short: "Export environment variables from a vault file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runExport(cmd, args)
	},
}

func init() {
	exportCmd.Flags().StringP("passphrase", "p", "", "Passphrase to decrypt the vault")
	exportCmd.Flags().StringP("format", "f", "dotenv", "Output format: dotenv or shell")
	rootCmd.AddCommand(exportCmd)
}

func runExport(cmd *cobra.Command, args []string) error {
	vaultPath := args[0]

	passphrase, err := cmd.Flags().GetString("passphrase")
	if err != nil {
		return err
	}
	if passphrase == "" {
		passphrase = os.Getenv("ENVVAULT_PASSPHRASE")
	}
	if passphrase == "" {
		return fmt.Errorf("passphrase is required: use --passphrase or set ENVVAULT_PASSPHRASE")
	}

	format, err := cmd.Flags().GetString("format")
	if err != nil {
		return err
	}
	format = strings.ToLower(strings.TrimSpace(format))

	v, err := vault.Load(vaultPath, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load vault: %w", err)
	}

	output, err := vault.Export(v, format)
	if err != nil {
		return fmt.Errorf("failed to export vault: %w", err)
	}

	fmt.Print(output)
	return nil
}
