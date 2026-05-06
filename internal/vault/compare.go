package vault

import "sort"

// CompareResult holds the result of comparing two vaults.
type CompareResult struct {
	OnlyInA    []string // keys present only in vault A
	OnlyInB    []string // keys present only in vault B
	SameValue  []string // keys present in both with identical values
	DiffValue  []string // keys present in both but with different values
}

// Compare performs a detailed side-by-side comparison of two vaults,
// returning which keys are unique to each, shared with equal values,
// or shared with differing values.
func Compare(a, b *Vault) CompareResult {
	aKeys := make(map[string]string)
	bKeys := make(map[string]string)

	for k, v := range a.data {
		aKeys[k] = v
	}
	for k, v := range b.data {
		bKeys[k] = v
	}

	var result CompareResult

	for k, av := range aKeys {
		if bv, ok := bKeys[k]; ok {
			if av == bv {
				result.SameValue = append(result.SameValue, k)
			} else {
				result.DiffValue = append(result.DiffValue, k)
			}
		} else {
			result.OnlyInA = append(result.OnlyInA, k)
		}
	}

	for k := range bKeys {
		if _, ok := aKeys[k]; !ok {
			result.OnlyInB = append(result.OnlyInB, k)
		}
	}

	sort.Strings(result.OnlyInA)
	sort.Strings(result.OnlyInB)
	sort.Strings(result.SameValue)
	sort.Strings(result.DiffValue)

	return result
}

// IsIdentical returns true when both vaults contain exactly the same
// keys with exactly the same values.
func IsIdentical(a, b *Vault) bool {
	r := Compare(a, b)
	return len(r.OnlyInA) == 0 && len(r.OnlyInB) == 0 && len(r.DiffValue) == 0
}
