// Package vault — alias subsystem
//
// # Overview
//
// Aliases provide short, human-friendly names that map to full vault key
// names. They are stored in a sidecar JSON file next to the vault file and
// are never encrypted, as they contain no secret material.
//
// # File layout
//
// Given a vault at /path/to/my.vault the alias record is stored at:
//
//	/path/to/.my.vault.aliases.json
//
// The file is written with mode 0600.
//
// # Usage
//
//	// Create an alias
//	vault.SetAlias(vaultPath, "db", "DATABASE_URL")
//
//	// Resolve before reading
//	key, _ := vault.ResolveAlias(vaultPath, "db")  // → "DATABASE_URL"
//
//	// List all aliases
//	pairs, _ := vault.ListAliases(vaultPath)
//
//	// Remove an alias
//	vault.RemoveAlias(vaultPath, "db")
package vault
