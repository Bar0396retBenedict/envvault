package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/yourorg/envvault/internal/vault"
)

func init() {
	quotaCmd := &cobra.Command{
		Use:   "quota",
		Short: "Manage vault key/size quotas",
	}

	setCmd := &cobra.Command{
		Use:   "set <vault> [--max-keys N] [--max-bytes N]",
		Short: "Set quota limits for a vault",
		Args:  cobra.ExactArgs(1),
		RunE:  runQuotaSet,
	}
	setCmd.Flags().Int("max-keys", 0, "Maximum number of keys (0 = unlimited)")
	setCmd.Flags().Int64("max-bytes", 0, "Maximum total bytes of key+value data (0 = unlimited)")

	checkCmd := &cobra.Command{
		Use:   "check <vault>",
		Short: "Check a vault against its quota",
		Args:  cobra.ExactArgs(1),
		RunE:  runQuotaCheck,
	}
	checkCmd.Flags().String("passphrase", "", "Vault passphrase")

	quotaCmd.AddCommand(setCmd, checkCmd)
	rootCmd.AddCommand(quotaCmd)
}

func runQuotaSet(cmd *cobra.Command, args []string) error {
	vaultPath := args[0]
	maxKeys, _ := cmd.Flags().GetInt("max-keys")
	maxBytes, _ := cmd.Flags().GetInt64("max-bytes")

	rec := vault.QuotaRecord{
		MaxKeys:  maxKeys,
		MaxBytes: maxBytes,
	}
	if err := vault.SaveQuotaRecord(vaultPath, rec); err != nil {
		return fmt.Errorf("quota set: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "quota saved: max_keys=%s max_bytes=%s\n",
		formatLimit(int64(maxKeys)), formatLimit(maxBytes))
	return nil
}

func runQuotaCheck(cmd *cobra.Command, args []string) error {
	vaultPath := args[0]
	passphrase, _ := cmd.Flags().GetString("passphrase")
	if passphrase == "" {
		return fmt.Errorf("--passphrase is required")
	}

	v, err := vault.Load(vaultPath, passphrase)
	if err != nil {
		return fmt.Errorf("load vault: %w", err)
	}
	rec, err := vault.LoadQuotaRecord(vaultPath)
	if err != nil {
		return fmt.Errorf("load quota: %w", err)
	}
	violations := vault.CheckQuota(v, rec)
	if len(violations) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "OK: vault is within quota limits")
		return nil
	}
	for _, viol := range violations {
		fmt.Fprintf(cmd.OutOrStdout(), "VIOLATION [%s]: %s\n", viol.Field, viol.Message)
	}
	os.Exit(1)
	return nil
}

func formatLimit(n int64) string {
	if n == 0 {
		return "unlimited"
	}
	return strconv.FormatInt(n, 10)
}
