package vault

import (
	"testing"
)

func newTestVault(t *testing.T, entries map[string]string) *Vault {
	t.Helper()
	v := New()
	for k, val := range entries {
		if err := v.Set(k, val); err != nil {
			t.Fatalf("Set(%q): %v", k, err)
		}
	}
	return v
}

func TestDiffNoChanges(t *testing.T) {
	src := newTestVault(t, map[string]string{"A": "1", "B": "2"})
	dst := newTestVault(t, map[string]string{"A": "1", "B": "2"})

	changes := Diff(src, dst)
	if len(changes) != 0 {
		t.Fatalf("expected 0 changes, got %d", len(changes))
	}
	if HasChanges(src, dst) {
		t.Fatal("HasChanges should be false")
	}
}

func TestDiffAdded(t *testing.T) {
	src := newTestVault(t, map[string]string{"A": "1"})
	dst := newTestVault(t, map[string]string{"A": "1", "B": "2"})

	changes := Diff(src, dst)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != ChangeAdded || changes[0].Key != "B" || changes[0].NewValue != "2" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestDiffRemoved(t *testing.T) {
	src := newTestVault(t, map[string]string{"A": "1", "B": "2"})
	dst := newTestVault(t, map[string]string{"A": "1"})

	changes := Diff(src, dst)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != ChangeRemoved || changes[0].Key != "B" || changes[0].OldValue != "2" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestDiffUpdated(t *testing.T) {
	src := newTestVault(t, map[string]string{"A": "old"})
	dst := newTestVault(t, map[string]string{"A": "new"})

	changes := Diff(src, dst)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	c := changes[0]
	if c.Kind != ChangeUpdated || c.OldValue != "old" || c.NewValue != "new" {
		t.Errorf("unexpected change: %+v", c)
	}
}

func TestDiffSorted(t *testing.T) {
	src := newTestVault(t, map[string]string{})
	dst := newTestVault(t, map[string]string{"Z": "1", "A": "2", "M": "3"})

	changes := Diff(src, dst)
	if len(changes) != 3 {
		t.Fatalf("expected 3 changes, got %d", len(changes))
	}
	keys := []string{changes[0].Key, changes[1].Key, changes[2].Key}
	expected := []string{"A", "M", "Z"}
	for i, k := range expected {
		if keys[i] != k {
			t.Errorf("position %d: want %q got %q", i, k, keys[i])
		}
	}
}
