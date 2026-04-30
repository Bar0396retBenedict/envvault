package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func writeEnvFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestImportDotEnv(t *testing.T) {
	dir := t.TempDir()
	path := writeEnvFile(t, dir, ".env", "FOO=bar\nBAZ=qux\n")
	v := newTestVault(t)
	res, err := ImportFromFile(v, path, ImportDotEnv, false)
	if err != nil {
		t.Fatal(err)
	}
	if res.Added != 2 || res.Skipped != 0 {
		t.Fatalf("want added=2 skipped=0, got %+v", res)
	}
	if val, _ := v.Get("FOO"); val != "bar" {
		t.Fatalf("expected bar, got %q", val)
	}
}

func TestImportShellFormat(t *testing.T) {
	dir := t.TempDir()
	path := writeEnvFile(t, dir, "vars.sh", "export ALPHA=one\nexport BETA=two\n")
	v := newTestVault(t)
	res, err := ImportFromFile(v, path, ImportShell, false)
	if err != nil {
		t.Fatal(err)
	}
	if res.Added != 2 {
		t.Fatalf("expected 2 added, got %d", res.Added)
	}
	if val, _ := v.Get("ALPHA"); val != "one" {
		t.Fatalf("expected one, got %q", val)
	}
}

func TestImportSkipsExisting(t *testing.T) {
	dir := t.TempDir()
	path := writeEnvFile(t, dir, ".env", "KEY=new\n")
	v := newTestVault(t)
	v.Set("KEY", "old")
	res, err := ImportFromFile(v, path, ImportDotEnv, false)
	if err != nil {
		t.Fatal(err)
	}
	if res.Skipped != 1 || res.Added != 0 {
		t.Fatalf("want skipped=1, got %+v", res)
	}
	if val, _ := v.Get("KEY"); val != "old" {
		t.Fatal("existing value should not be overwritten")
	}
}

func TestImportOverwrite(t *testing.T) {
	dir := t.TempDir()
	path := writeEnvFile(t, dir, ".env", "KEY=new\n")
	v := newTestVault(t)
	v.Set("KEY", "old")
	res, err := ImportFromFile(v, path, ImportDotEnv, true)
	if err != nil {
		t.Fatal(err)
	}
	if res.Overwritten != 1 {
		t.Fatalf("want overwritten=1, got %+v", res)
	}
	if val, _ := v.Get("KEY"); val != "new" {
		t.Fatal("value should be overwritten")
	}
}

func TestImportIgnoresCommentsAndBlanks(t *testing.T) {
	dir := t.TempDir()
	content := "# comment\n\nVALID=yes\n"
	path := writeEnvFile(t, dir, ".env", content)
	v := newTestVault(t)
	res, err := ImportFromFile(v, path, ImportDotEnv, false)
	if err != nil {
		t.Fatal(err)
	}
	if res.Added != 1 {
		t.Fatalf("expected 1 added, got %d", res.Added)
	}
}

func TestImportMissingFile(t *testing.T) {
	v := newTestVault(t)
	_, err := ImportFromFile(v, "/nonexistent/.env", ImportDotEnv, false)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
