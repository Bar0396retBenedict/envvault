// Package vault — history module
//
// # Key Change History
//
// The history module records every mutation made to vault keys, providing an
// auditable trail of what changed, when, and what the previous value was.
//
// # Storage
//
// History is stored in a hidden JSON file alongside the vault:
//
//	.<vault-file>.history.json
//
// The file is written with 0600 permissions so only the owning user can read it.
//
// # Entry Format
//
// Each HistoryEntry contains:
//   - Timestamp — UTC time of the change
//   - Key       — the vault key that was modified
//   - OldValue  — previous plaintext value (empty for new keys)
//   - NewValue  — new plaintext value (empty for deletions)
//   - Action    — one of "set", "delete", or "rename"
//
// # Usage
//
//	entry := vault.HistoryEntry{Key: "API_KEY", OldValue: "x", NewValue: "y", Action: "set"}
//	err := vault.AppendHistory("/path/to/vault", entry)
//
//	rec, err := vault.LoadHistory("/path/to/vault")
//	for _, e := range rec.Entries {
//		fmt.Printf("%s  %s  %s -> %s\n", e.Timestamp, e.Key, e.OldValue, e.NewValue)
//	}
package vault
