// Package vault provides audit logging for vault operations.
//
// # Audit Log
//
// Every mutating operation on a vault (set, delete, rotate, import, copy)
// can be recorded in a companion audit log file. The log is stored alongside
// the vault file with an ".audit.json" extension.
//
// Example vault file:  prod.vault
// Companion audit log: prod.audit.json
//
// Each entry captures:
//   - Timestamp: UTC time of the operation
//   - Action:    one of set, delete, rotate, import, copy
//   - Key:       the affected variable name (empty for bulk operations)
//   - VaultPath: absolute or relative path to the vault file
//
// # Usage
//
//	// Record that API_KEY was set in prod.vault
//	err := vault.AppendAuditEntry("prod.vault", vault.AuditSet, "API_KEY")
//
//	// Read the full history
//	log, err := vault.LoadAuditLog("prod.vault")
//	for _, e := range log.Entries {
//	    fmt.Println(e.Timestamp, e.Action, e.Key)
//	}
//
// The audit log is append-only from the perspective of the API; entries are
// never removed by envvault itself.
package vault
