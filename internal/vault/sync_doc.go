// Package vault provides sync utilities for reconciling vault state
// across multiple environments (local, staging, production).
//
// # Sync Metadata
//
// The SyncRecord type tracks when each environment's vault was last
// synchronized, along with a checksum of the vault contents at that time.
// This allows the CLI to detect drift between environments.
//
// Sync records are stored in a separate JSON file (e.g. .envvault.sync)
// and are never encrypted, since they contain no secret material.
//
// # Merging Vaults
//
// MergeVaults copies keys from a source vault into a destination vault.
// It is intentionally non-destructive: keys present only in the destination
// are preserved. Callers are responsible for saving the destination vault
// after a merge.
//
// Typical sync workflow:
//
//	// 1. Load both vaults with their respective passphrases.
//	// 2. Call MergeVaults(local, remote) to pull remote changes locally.
//	// 3. Save the updated local vault.
//	// 4. Update the SyncRecord and call SaveSyncRecord.
package vault
