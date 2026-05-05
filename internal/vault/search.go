package vault

import (
	"sort"
	"strings"
)

// SearchResult holds a matched key and its value from a vault.
type SearchResult struct {
	Key   string
	Value string
}

// SearchOptions controls how Search behaves.
type SearchOptions struct {
	// CaseSensitive disables case-folding when matching keys and values.
	CaseSensitive bool
	// MatchValue also checks whether the query appears in the value.
	MatchValue bool
}

// Search returns all entries in v whose key (or optionally value) contains
// query. Results are returned in lexicographic key order.
func Search(v *Vault, query string, opts SearchOptions) []SearchResult {
	if !opts.CaseSensitive {
		query = strings.ToLower(query)
	}

	var results []SearchResult

	for _, k := range v.Keys() {
		val, _ := v.Get(k)

		keyToMatch := k
		valToMatch := val
		if !opts.CaseSensitive {
			keyToMatch = strings.ToLower(k)
			valToMatch = strings.ToLower(val)
		}

		keyMatches := strings.Contains(keyToMatch, query)
		valMatches := opts.MatchValue && strings.Contains(valToMatch, query)

		if keyMatches || valMatches {
			results = append(results, SearchResult{Key: k, Value: val})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})

	return results
}
