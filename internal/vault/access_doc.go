// Package vault — access tracking.
//
// # Access Tracking
//
// The access subsystem records per-key read and write events for a vault file.
// Statistics are persisted in a sidecar JSON file named
// ".<vault-basename>.access.json" in the same directory as the vault.
//
// # Usage
//
// Record a read event after retrieving a key's value:
//
//	vault.RecordRead(vaultPath, "API_KEY")
//
// Record a write event after setting or deleting a key:
//
//	vault.RecordWrite(vaultPath, "DB_PASSWORD")
//
// Retrieve a sorted list of all tracked entries:
//
//	entries, err := vault.ListAccessEntries(vaultPath)
//
// Each AccessEntry exposes ReadCount, WriteCount, LastRead, and LastWrite
// fields so operators can audit which keys are accessed most frequently.
//
// # File Format
//
// The sidecar file is a JSON object with an "entries" map keyed by variable
// name. The file is written with mode 0600 to prevent other users from
// reading access metadata.
package vault
