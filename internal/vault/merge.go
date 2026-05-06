package vault

import "fmt"

// MergeStrategy controls how conflicts are resolved during a merge.
type MergeStrategy int

const (
	// MergeStrategyOurs keeps the value from the destination vault on conflict.
	MergeStrategyOurs MergeStrategy = iota
	// MergeStrategyTheirs keeps the value from the source vault on conflict.
	MergeStrategyTheirs
	// MergeStrategyError returns an error when a conflict is detected.
	MergeStrategyError
)

// MergeResult summarises the outcome of a MergeInto call.
type MergeResult struct {
	Added    []string
	Updated  []string
	Skipped  []string
	Conflict []string
}

// MergeInto merges all keys from src into dst using the given strategy.
// src and dst must already be loaded (plaintext). The caller is responsible
// for saving dst afterward.
func MergeInto(dst, src *Vault, strategy MergeStrategy) (MergeResult, error) {
	var result MergeResult

	for key, srcVal := range src.data {
		dstVal, exists := dst.data[key]
		if !exists {
			dst.data[key] = srcVal
			result.Added = append(result.Added, key)
			continue
		}
		if dstVal == srcVal {
			result.Skipped = append(result.Skipped, key)
			continue
		}
		// Conflict: values differ.
		switch strategy {
		case MergeStrategyOurs:
			result.Skipped = append(result.Skipped, key)
		case MergeStrategyTheirs:
			dst.data[key] = srcVal
			result.Updated = append(result.Updated, key)
		case MergeStrategyError:
			result.Conflict = append(result.Conflict, key)
		default:
			return result, fmt.Errorf("unknown merge strategy: %d", strategy)
		}
	}

	if len(result.Conflict) > 0 {
		return result, fmt.Errorf("merge conflict on keys: %v", result.Conflict)
	}
	return result, nil
}
