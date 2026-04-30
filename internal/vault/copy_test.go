package vault

import (
	"testing"
)

func makeCopyVault(t *testing.T, pairs map[string]string) *Vault {
	t.Helper()
	v := New()
	for k, val := range pairs {
		v.Set(k, val)
	}
	return v
}

func TestCopyAllKeys(t *testing.T) {
	src := makeCopyVault(t, map[string]string{"A": "1", "B": "2"})
	dst := New()

	res, err := CopyKeys(src, dst, nil, CopyOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Copied) != 2 {
		t.Errorf("expected 2 copied, got %d", len(res.Copied))
	}
	if len(res.Skipped) != 0 {
		t.Errorf("expected 0 skipped, got %d", len(res.Skipped))
	}
	for _, k := range []string{"A", "B"} {
		if v, ok := dst.Get(k); !ok || v == "" {
			t.Errorf("key %q missing or empty in destination", k)
		}
	}
}

func TestCopySelectedKeys(t *testing.T) {
	src := makeCopyVault(t, map[string]string{"A": "1", "B": "2", "C": "3"})
	dst := New()

	res, err := CopyKeys(src, dst, []string{"A", "C"}, CopyOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Copied) != 2 {
		t.Errorf("expected 2 copied, got %d", len(res.Copied))
	}
	if _, ok := dst.Get("B"); ok {
		t.Error("key B should not have been copied")
	}
}

func TestCopySkipsExistingWithoutOverwrite(t *testing.T) {
	src := makeCopyVault(t, map[string]string{"X": "new"})
	dst := makeCopyVault(t, map[string]string{"X": "old"})

	res, err := CopyKeys(src, dst, nil, CopyOptions{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(res.Skipped))
	}
	if v, _ := dst.Get("X"); v != "old" {
		t.Errorf("expected original value 'old', got %q", v)
	}
}

func TestCopyOverwritesExisting(t *testing.T) {
	src := makeCopyVault(t, map[string]string{"X": "new"})
	dst := makeCopyVault(t, map[string]string{"X": "old"})

	_, err := CopyKeys(src, dst, nil, CopyOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, _ := dst.Get("X"); v != "new" {
		t.Errorf("expected overwritten value 'new', got %q", v)
	}
}

func TestCopyMissingKeyReturnsError(t *testing.T) {
	src := New()
	dst := New()

	_, err := CopyKeys(src, dst, []string{"MISSING"}, CopyOptions{})
	if err == nil {
		t.Error("expected error for missing key, got nil")
	}
}

func TestCopyNilSourceReturnsError(t *testing.T) {
	_, err := CopyKeys(nil, New(), nil, CopyOptions{})
	if err == nil {
		t.Error("expected error for nil source")
	}
}

func TestCopyNilDestinationReturnsError(t *testing.T) {
	_, err := CopyKeys(New(), nil, nil, CopyOptions{})
	if err == nil {
		t.Error("expected error for nil destination")
	}
}
