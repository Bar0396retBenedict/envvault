package vault

import "fmt"

// RenameKey renames an existing key within a vault to a new name.
// It returns an error if the source key does not exist, if the destination
// key already exists (and overwrite is false), or if either key name is empty.
func RenameKey(v *Vault, oldKey, newKey string, overwrite bool) error {
	if oldKey == "" {
		return fmt.Errorf("rename: source key must not be empty")
	}
	if newKey == "" {
		return fmt.Errorf("rename: destination key must not be empty")
	}
	if oldKey == newKey {
		return fmt.Errorf("rename: source and destination keys are identical")
	}

	val, ok := v.Get(oldKey)
	if !ok {
		return fmt.Errorf("rename: key %q not found", oldKey)
	}

	if _, exists := v.Get(newKey); exists && !overwrite {
		return fmt.Errorf("rename: destination key %q already exists; use --overwrite to replace it", newKey)
	}

	v.Set(newKey, val)
	v.Delete(oldKey)
	return nil
}
