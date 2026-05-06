package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func makeEnvVault(t *testing.T) (string, *Vault) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.vault")
	v := New()
	v.Set("DB_PASS", "s3cr3t")
	v.Set("API_KEY", "abc123")
	if err := v.Save(path, "passphrase"); err != nil {
		t.Fatalf("save vault: %v", err)
	}
	return path, v
}

func TestLoadEnvRecordMissing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ghost.vault")
	rec, err := LoadEnvRecord(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rec.Bindings) != 0 {
		t.Errorf("expected empty bindings, got %v", rec.Bindings)
	}
}

func TestBindAndUnbindEnvVar(t *testing.T) {
	path, _ := makeEnvVault(t)
	if err := BindEnvVar(path, "MY_DB", "DB_PASS"); err != nil {
		t.Fatalf("bind: %v", err)
	}
	rec, err := LoadEnvRecord(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if rec.Bindings["MY_DB"] != "DB_PASS" {
		t.Errorf("expected DB_PASS, got %q", rec.Bindings["MY_DB"])
	}
	if err := UnbindEnvVar(path, "MY_DB"); err != nil {
		t.Fatalf("unbind: %v", err)
	}
	rec, _ = LoadEnvRecord(path)
	if _, ok := rec.Bindings["MY_DB"]; ok {
		t.Error("binding should have been removed")
	}
}

func TestUnbindNotBound(t *testing.T) {
	path, _ := makeEnvVault(t)
	err := UnbindEnvVar(path, "NONEXISTENT")
	if err == nil {
		t.Fatal("expected error for unbound key")
	}
}

func TestBindEmptyEnvKey(t *testing.T) {
	path, _ := makeEnvVault(t)
	if err := BindEnvVar(path, "", "DB_PASS"); err == nil {
		t.Fatal("expected error for empty env key")
	}
}

func TestApplyEnvBindings(t *testing.T) {
	path, _ := makeEnvVault(t)
	if err := BindEnvVar(path, "TEST_APPLY_DB", "DB_PASS"); err != nil {
		t.Fatalf("bind: %v", err)
	}
	applied, err := ApplyEnvBindings(path, "passphrase")
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	if len(applied) != 1 || applied[0] != "TEST_APPLY_DB" {
		t.Errorf("unexpected applied list: %v", applied)
	}
	if got := os.Getenv("TEST_APPLY_DB"); got != "s3cr3t" {
		t.Errorf("expected s3cr3t, got %q", got)
	}
	os.Unsetenv("TEST_APPLY_DB")
}

func TestApplyEnvBindingsMissingVaultKey(t *testing.T) {
	path, _ := makeEnvVault(t)
	if err := BindEnvVar(path, "MISSING_KEY_ENV", "NO_SUCH_KEY"); err != nil {
		t.Fatalf("bind: %v", err)
	}
	_, err := ApplyEnvBindings(path, "passphrase")
	if err == nil {
		t.Fatal("expected error for missing vault key")
	}
}
