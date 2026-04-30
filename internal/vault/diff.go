package vault

import "sort"

// ChangeKind describes the type of change between two vaults.
type ChangeKind string

const (
	ChangeAdded   ChangeKind = "added"
	ChangeRemoved ChangeKind = "removed"
	ChangeUpdated ChangeKind = "updated"
)

// Change represents a single key-level difference between two vaults.
type Change struct {
	Key      string
	Kind     ChangeKind
	OldValue string // empty for added
	NewValue string // empty for removed
}

// Diff compares two vaults and returns the ordered list of changes
// needed to transform src into dst.
func Diff(src, dst *Vault) []Change {
	srcData := src.All()
	dstData := dst.All()

	var changes []Change

	for k, dv := range dstData {
		if sv, ok := srcData[k]; !ok {
			changes = append(changes, Change{Key: k, Kind: ChangeAdded, NewValue: dv})
		} else if sv != dv {
			changes = append(changes, Change{Key: k, Kind: ChangeUpdated, OldValue: sv, NewValue: dv})
		}
	}

	for k, sv := range srcData {
		if _, ok := dstData[k]; !ok {
			changes = append(changes, Change{Key: k, Kind: ChangeRemoved, OldValue: sv})
		}
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Key < changes[j].Key
	})

	return changes
}

// HasChanges returns true when Diff finds at least one difference.
func HasChanges(src, dst *Vault) bool {
	return len(Diff(src, dst)) > 0
}
