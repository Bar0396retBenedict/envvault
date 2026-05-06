package vault

import (
	"testing"
)

func makeCompareVault(t *testing.T, pairs ...string) *Vault {
	t.Helper()
	v := &Vault{data: make(map[string]string)}
	for i := 0; i+1 < len(pairs); i += 2 {
		v.data[pairs[i]] = pairs[i+1]
	}
	return v
}

func TestCompareIdenticalVaults(t *testing.T) {
	a := makeCompareVault(t, "KEY", "val", "FOO", "bar")
	b := makeCompareVault(t, "KEY", "val", "FOO", "bar")

	r := Compare(a, b)

	if len(r.OnlyInA) != 0 || len(r.OnlyInB) != 0 || len(r.DiffValue) != 0 {
		t.Errorf("expected identical vaults, got %+v", r)
	}
	if len(r.SameValue) != 2 {
		t.Errorf("expected 2 same-value keys, got %d", len(r.SameValue))
	}
}

func TestCompareOnlyInA(t *testing.T) {
	a := makeCompareVault(t, "ALPHA", "1", "SHARED", "x")
	b := makeCompareVault(t, "SHARED", "x")

	r := Compare(a, b)

	if len(r.OnlyInA) != 1 || r.OnlyInA[0] != "ALPHA" {
		t.Errorf("expected [ALPHA] in OnlyInA, got %v", r.OnlyInA)
	}
	if len(r.OnlyInB) != 0 {
		t.Errorf("expected empty OnlyInB, got %v", r.OnlyInB)
	}
}

func TestCompareOnlyInB(t *testing.T) {
	a := makeCompareVault(t, "SHARED", "x")
	b := makeCompareVault(t, "BETA", "2", "SHARED", "x")

	r := Compare(a, b)

	if len(r.OnlyInB) != 1 || r.OnlyInB[0] != "BETA" {
		t.Errorf("expected [BETA] in OnlyInB, got %v", r.OnlyInB)
	}
	if len(r.OnlyInA) != 0 {
		t.Errorf("expected empty OnlyInA, got %v", r.OnlyInA)
	}
}

func TestCompareDiffValue(t *testing.T) {
	a := makeCompareVault(t, "KEY", "old")
	b := makeCompareVault(t, "KEY", "new")

	r := Compare(a, b)

	if len(r.DiffValue) != 1 || r.DiffValue[0] != "KEY" {
		t.Errorf("expected [KEY] in DiffValue, got %v", r.DiffValue)
	}
	if len(r.SameValue) != 0 {
		t.Errorf("expected empty SameValue, got %v", r.SameValue)
	}
}

func TestCompareSortedOutput(t *testing.T) {
	a := makeCompareVault(t, "Z", "1", "A", "1", "M", "1")
	b := makeCompareVault(t, "Z", "1", "A", "1", "M", "1")

	r := Compare(a, b)

	for i := 1; i < len(r.SameValue); i++ {
		if r.SameValue[i] < r.SameValue[i-1] {
			t.Errorf("SameValue not sorted: %v", r.SameValue)
		}
	}
}

func TestIsIdenticalTrue(t *testing.T) {
	a := makeCompareVault(t, "X", "1")
	b := makeCompareVault(t, "X", "1")
	if !IsIdentical(a, b) {
		t.Error("expected vaults to be identical")
	}
}

func TestIsIdenticalFalse(t *testing.T) {
	a := makeCompareVault(t, "X", "1")
	b := makeCompareVault(t, "X", "2")
	if IsIdentical(a, b) {
		t.Error("expected vaults to differ")
	}
}
