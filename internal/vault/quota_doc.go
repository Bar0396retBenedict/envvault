// Package vault — quota subsystem
//
// # Overview
//
// The quota subsystem allows operators to place hard limits on the number of
// keys and the total byte size of key+value data stored inside a vault file.
// Limits are stored in a sidecar JSON file named .<vault>.quota.json and are
// checked on demand via CheckQuota.
//
// # File layout
//
//	<vault>.vault          — encrypted vault data
//	.<vault>.vault.quota.json — quota configuration (0600)
//
// # Usage
//
//	rec, err := vault.LoadQuotaRecord(vaultPath)
//	violations := vault.CheckQuota(v, rec)
//	for _, viol := range violations {
//	    fmt.Println(viol.Message)
//	}
//
// # Limits
//
//   - MaxKeys  — maximum number of keys (0 = unlimited)
//   - MaxBytes — maximum combined byte length of all keys and values (0 = unlimited)
package vault
