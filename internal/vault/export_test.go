package vault_test

import (
	"strings"
	"testing"

	"github.com/user/envvault/internal/vault"
)

func TestExportDotEnv(t *testing.T) {
	v := vault.New()
	v.Set("APP_ENV", "production")
	v.Set("DATABASE_URL", "postgres://localhost/db")
	v.Set("PORT", "8080")

	var sb strings.Builder
	if err := v.Export(&sb, vault.FormatDotEnv); err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	output := sb.String()
	for _, expected := range []string{
		"APP_ENV=production\n",
		"DATABASE_URL=postgres://localhost/db\n",
		"PORT=8080\n",
	} {
		if !strings.Contains(output, expected) {
			t.Errorf("expected output to contain %q, got:\n%s", expected, output)
		}
	}
}

func TestExportShell(t *testing.T) {
	v := vault.New()
	v.Set("API_KEY", "abc123")

	var sb strings.Builder
	if err := v.Export(&sb, vault.FormatShell); err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	output := sb.String()
	if !strings.Contains(output, "export API_KEY=abc123") {
		t.Errorf("expected shell export format, got: %s", output)
	}
}

func TestExportSortedKeys(t *testing.T) {
	v := vault.New()
	v.Set("ZEBRA", "z")
	v.Set("ALPHA", "a")
	v.Set("MIDDLE", "m")

	var sb strings.Builder
	if err := v.Export(&sb, vault.FormatDotEnv); err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(sb.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "ALPHA") {
		t.Errorf("expected first line to be ALPHA, got: %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "ZEBRA") {
		t.Errorf("expected last line to be ZEBRA, got: %s", lines[2])
	}
}

func TestExportValueWithSpaces(t *testing.T) {
	v := vault.New()
	v.Set("MSG", "hello world")

	var sb strings.Builder
	if err := v.Export(&sb, vault.FormatDotEnv); err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	output := sb.String()
	if !strings.Contains(output, "MSG='hello world'") {
		t.Errorf("expected quoted value for space-containing string, got: %s", output)
	}
}

func TestExportUnknownFormat(t *testing.T) {
	v := vault.New()
	v.Set("KEY", "val")

	var sb strings.Builder
	err := v.Export(&sb, vault.ExportFormat("xml"))
	if err == nil {
		t.Error("expected error for unknown format")
	}
}
