package vault

import (
	"path/filepath"
	"testing"
)

func makeMergeVault(t *testing.T, data map[string]string) *Vault {
	t.Helper()
	dir := t.TempDir()
	v, err := New(filepath.Join(dir, "vault.env"), "passphrase")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	for k, val := range data {
		v.Set(k, val)
	}
	return v
}

func TestMergeIntoAddsNewKeys(t *testing.T) {
	dst := makeMergeVault(t, map[string]string{"A": "1"})
	src := makeMergeVault(t, map[string]string{"B": "2", "C": "3"})

	res, err := MergeInto(dst, src, MergeStrategyOurs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Added) != 2 {
		t.Errorf("expected 2 added, got %d", len(res.Added))
	}
	if v, _ := dst.Get("B"); v != "2" {
		t.Errorf("expected B=2, got %q", v)
	}
}

func TestMergeIntoStrategyOursKeepsDst(t *testing.T) {
	dst := makeMergeVault(t, map[string]string{"KEY": "dst-value"})
	src := makeMergeVault(t, map[string]string{"KEY": "src-value"})

	res, err := MergeInto(dst, src, MergeStrategyOurs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(res.Skipped))
	}
	if v, _ := dst.Get("KEY"); v != "dst-value" {
		t.Errorf("expected dst-value, got %q", v)
	}
}

func TestMergeIntoStrategyTheirsOverwrites(t *testing.T) {
	dst := makeMergeVault(t, map[string]string{"KEY": "old"})
	src := makeMergeVault(t, map[string]string{"KEY": "new"})

	res, err := MergeInto(dst, src, MergeStrategyTheirs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Updated) != 1 {
		t.Errorf("expected 1 updated, got %d", len(res.Updated))
	}
	if v, _ := dst.Get("KEY"); v != "new" {
		t.Errorf("expected new, got %q", v)
	}
}

func TestMergeIntoStrategyErrorOnConflict(t *testing.T) {
	dst := makeMergeVault(t, map[string]string{"KEY": "a"})
	src := makeMergeVault(t, map[string]string{"KEY": "b"})

	res, err := MergeInto(dst, src, MergeStrategyError)
	if err == nil {
		t.Fatal("expected error on conflict, got nil")
	}
	if len(res.Conflict) != 1 || res.Conflict[0] != "KEY" {
		t.Errorf("expected conflict on KEY, got %v", res.Conflict)
	}
}

func TestMergeIntoIdenticalSkipsAll(t *testing.T) {
	dst := makeMergeVault(t, map[string]string{"X": "1", "Y": "2"})
	src := makeMergeVault(t, map[string]string{"X": "1", "Y": "2"})

	res, err := MergeInto(dst, src, MergeStrategyOurs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 2 {
		t.Errorf("expected 2 skipped, got %d", len(res.Skipped))
	}
	if len(res.Added) != 0 || len(res.Updated) != 0 {
		t.Errorf("expected no adds/updates")
	}
}
