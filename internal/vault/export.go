package vault

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// ExportFormat represents the output format for exporting vault contents.
type ExportFormat string

const (
	// FormatDotEnv exports variables in KEY=VALUE format.
	FormatDotEnv ExportFormat = "dotenv"
	// FormatShell exports variables as shell export statements.
	FormatShell ExportFormat = "shell"
)

// Export writes the vault's environment variables to w in the specified format.
// Keys are written in sorted order for deterministic output.
func (v *Vault) Export(w io.Writer, format ExportFormat) error {
	keys := make([]string, 0, len(v.Env))
	for k := range v.Env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	switch format {
	case FormatDotEnv:
		for _, k := range keys {
			_, err := fmt.Fprintf(w, "%s=%s\n", k, quoteValue(v.Env[k]))
			if err != nil {
				return fmt.Errorf("vault: export dotenv: %w", err)
			}
		}
	case FormatShell:
		for _, k := range keys {
			_, err := fmt.Fprintf(w, "export %s=%s\n", k, quoteValue(v.Env[k]))
			if err != nil {
				return fmt.Errorf("vault: export shell: %w", err)
			}
		}
	default:
		return fmt.Errorf("vault: unknown export format: %q", format)
	}
	return nil
}

// quoteValue wraps a value in single quotes if it contains spaces or
// special shell characters, escaping any existing single quotes.
func quoteValue(v string) string {
	if v == "" {
		return `""`
	}
	special := " \t\n$`&|;<>(){}!#~"
	for _, ch := range special {
		if strings.ContainsRune(v, ch) {
			escaped := strings.ReplaceAll(v, "'", `'\''`)
			return "'" + escaped + "'"
		}
	}
	return v
}
