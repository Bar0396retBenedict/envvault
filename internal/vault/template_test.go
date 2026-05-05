package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func makeTemplateVault(t *testing.T, pairs map[string]string) *Vault {
	t.Helper()
	v := New()
	for k, val := range pairs {
		v.Set(k, val)
	}
	return v
}

func writeTmplFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "tmpl.env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("write template: %v", err)
	}
	return p
}

func TestRenderTemplateAllPresent(t *testing.T) {
	v := makeTemplateVault(t, map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	})
	p := writeTmplFile(t, "host={{DB_HOST}} port={{DB_PORT}}")
	res, err := RenderTemplate(v, p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Output != "host=localhost port=5432" {
		t.Errorf("got %q", res.Output)
	}
	if len(res.Missing) != 0 {
		t.Errorf("expected no missing keys, got %v", res.Missing)
	}
}

func TestRenderTemplateMissingKey(t *testing.T) {
	v := makeTemplateVault(t, map[string]string{"DB_HOST": "localhost"})
	p := writeTmplFile(t, "host={{DB_HOST}} pass={{DB_PASS}}")
	res, err := RenderTemplate(v, p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Missing) != 1 || res.Missing[0] != "DB_PASS" {
		t.Errorf("expected [DB_PASS] missing, got %v", res.Missing)
	}
	// Placeholder should remain in output
	if res.Output != "host=localhost pass={{DB_PASS}}" {
		t.Errorf("got %q", res.Output)
	}
}

func TestRenderTemplateDuplicatePlaceholder(t *testing.T) {
	v := makeTemplateVault(t, map[string]string{})
	p := writeTmplFile(t, "{{MISSING}} and {{MISSING}} again")
	res, err := RenderTemplate(v, p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Missing) != 1 {
		t.Errorf("expected deduplication, got %v", res.Missing)
	}
}

func TestRenderTemplateFileNotFound(t *testing.T) {
	v := makeTemplateVault(t, map[string]string{})
	_, err := RenderTemplate(v, "/nonexistent/path/tmpl.env")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRenderTemplateNoPlaceholders(t *testing.T) {
	v := makeTemplateVault(t, map[string]string{"KEY": "val"})
	p := writeTmplFile(t, "plain text with no placeholders")
	res, err := RenderTemplate(v, p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Output != "plain text with no placeholders" {
		t.Errorf("got %q", res.Output)
	}
	if len(res.Missing) != 0 {
		t.Errorf("expected no missing, got %v", res.Missing)
	}
}
