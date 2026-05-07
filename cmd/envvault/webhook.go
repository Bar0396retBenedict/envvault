package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envvault/internal/vault"
)

func init() {
	webhookCmd := &cobra.Command{
		Use:   "webhook",
		Short: "Manage vault webhooks",
	}

	addCmd := &cobra.Command{
		Use:   "add <name> <url> [events...]",
		Short: "Register a webhook (events: set, delete, rotate)",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			vaultPath, _ := cmd.Flags().GetString("vault")
			return runWebhookAdd(vaultPath, args[0], args[1], args[2:])
		},
	}
	addCmd.Flags().String("vault", "vault.env", "path to vault file")

	removeCmd := &cobra.Command{
		Use:   "remove <name>",
		Short: "Deregister a webhook by name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			vaultPath, _ := cmd.Flags().GetString("vault")
			return runWebhookRemove(vaultPath, args[0])
		},
	}
	removeCmd.Flags().String("vault", "vault.env", "path to vault file")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all registered webhooks",
		RunE: func(cmd *cobra.Command, args []string) error {
			vaultPath, _ := cmd.Flags().GetString("vault")
			return runWebhookList(vaultPath)
		},
	}
	listCmd.Flags().String("vault", "vault.env", "path to vault file")

	webhookCmd.AddCommand(addCmd, removeCmd, listCmd)
	rootCmd.AddCommand(webhookCmd)
}

func runWebhookAdd(vaultPath, name, url string, rawEvents []string) error {
	events := make([]vault.WebhookEvent, 0, len(rawEvents))
	for _, e := range rawEvents {
		switch strings.ToLower(e) {
		case "set":
			events = append(events, vault.EventSet)
		case "delete":
			events = append(events, vault.EventDelete)
		case "rotate":
			events = append(events, vault.EventRotate)
		default:
			return fmt.Errorf("unknown event %q (valid: set, delete, rotate)", e)
		}
	}
	if err := vault.RegisterWebhook(vaultPath, name, url, events); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "webhook %q registered\n", name)
	return nil
}

func runWebhookRemove(vaultPath, name string) error {
	if err := vault.DeregisterWebhook(vaultPath, name); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "webhook %q removed\n", name)
	return nil
}

func runWebhookList(vaultPath string) error {
	names, rec, err := vault.ListWebhooks(vaultPath)
	if err != nil {
		return err
	}
	if len(names) == 0 {
		fmt.Fprintln(os.Stdout, "no webhooks registered")
		return nil
	}
	fmt.Fprintf(os.Stdout, "%-20s %-40s %s\n", "NAME", "URL", "EVENTS")
	for _, n := range names {
		e := rec.Hooks[n]
		evStrs := make([]string, len(e.Events))
		for i, ev := range e.Events {
			evStrs[i] = string(ev)
		}
		fmt.Fprintf(os.Stdout, "%-20s %-40s %s\n", n, e.URL, strings.Join(evStrs, ","))
	}
	return nil
}
