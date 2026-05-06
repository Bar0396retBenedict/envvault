package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func makeAliasVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "test.vault")
}

func TestLoadAliasRecordMissing(t *testing.T) {
	vp := makeAliasVault(t)
	rec, err := LoadAliasRecord(vp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rec.Aliases) != 0 {
		t.Fatalf("expected empty aliases, got %v", rec.Aliases)
	}
}

func TestSetAndResolveAlias(t *testing.T) {
	vp := makeAliasVault(t)
	if err := SetAlias(vp, "db", "DATABASE_URL"); err != nil {
		t.Fatalf("SetAlias: %v", err)
	}
	key, err := ResolveAlias(vp, "db")
	if err != nil {
		t.Fatalf("ResolveAlias: %v", err)
	}
	if key != "DATABASE_URL" {
		t.Fatalf("expected DATABASE_URL, got %s", key)
	}
}

func TestResolveAliasPassthrough(t *testing.T) {
	vp := makeAliasVault(t)
	key, err := ResolveAlias(vp, "UNKNOWN_KEY")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if key != "UNKNOWN_KEY" {
		t.Fatalf("expected pass-through, got %s", key)
	}
}

func TestRemoveAlias(t *testing.T) {
	vp := makeAliasVault(t)
	_ = SetAlias(vp, "db", "DATABASE_URL")
	if err := RemoveAlias(vp, "db"); err != nil {
		t.Fatalf("RemoveAlias: %v", err)
	}
	rec, _ := LoadAliasRecord(vp)
	if _, ok := rec.Aliases["db"]; ok {
		t.Fatal("alias should have been removed")
	}
}

func TestRemoveAliasNotFound(t *testing.T) {
	vp := makeAliasVault(t)
	err := RemoveAlias(vp, "ghost")
	if err == nil {
		t.Fatal("expected error for missing alias")
	}
}

func TestListAliasesSorted(t *testing.T) {
	vp := makeAliasVault(t)
	_ = SetAlias(vp, "z_key", "Z_VAR")
	_ = SetAlias(vp, "a_key", "A_VAR")
	_ = SetAlias(vp, "m_key", "M_VAR")
	list, err := ListAliases(vp)
	if err != nil {
		t.Fatalf("ListAliases: %v", err)
	}
	if len(list) != 3 {
		t.Fatalf("expected 3 aliases, got %d", len(list))
	}
	if list[0][0] != "a_key" || list[1][0] != "m_key" || list[2][0] != "z_key" {
		t.Fatalf("unexpected order: %v", list)
	}
}

func TestSetAliasEmptyName(t *testing.T) {
	vp := makeAliasVault(t)
	if err := SetAlias(vp, "", "SOME_KEY"); err == nil {
		t.Fatal("expected error for empty alias name")
	}
}

func TestAliasFilePermissions(t *testing.T) {
	vp := makeAliasVault(t)
	_ = SetAlias(vp, "key", "VAL")
	info, err := os.Stat(aliasFilePath(vp))
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("expected 0600, got %v", info.Mode().Perm())
	}
}
