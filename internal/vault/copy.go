package vault

import "fmt"

// CopyOptions controls the behaviour of CopyKeys.
type CopyOptions struct {
	// Overwrite existing keys in the destination vault.
	Overwrite bool
}

// CopyResult summarises what happened during a copy operation.
type CopyResult struct {
	Copied  []string
	Skipped []string
}

// CopyKeys copies the given keys (or all keys when keys is empty) from src
// into dst.  It returns a CopyResult describing which keys were copied and
// which were skipped because they already existed (and Overwrite is false).
func CopyKeys(src, dst *Vault, keys []string, opts CopyOptions) (CopyResult, error) {
	if src == nil {
		return CopyResult{}, fmt.Errorf("copy: source vault is nil")
	}
	if dst == nil {
		return CopyResult{}, fmt.Errorf("copy: destination vault is nil")
	}

	targets := keys
	if len(targets) == 0 {
		for k := range src.data {
			targets = append(targets, k)
		}
	}

	var result CopyResult
	for _, k := range targets {
		v, ok := src.Get(k)
		if !ok {
			return result, fmt.Errorf("copy: key %q not found in source vault", k)
		}

		_, exists := dst.Get(k)
		if exists && !opts.Overwrite {
			result.Skipped = append(result.Skipped, k)
			continue
		}

		dst.Set(k, v)
		result.Copied = append(result.Copied, k)
	}

	return result, nil
}
