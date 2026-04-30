// Package vault – import subsystem
//
// # Importing Environment Variables
//
// The import subsystem allows loading key=value pairs from existing .env or
// shell-export files into an open Vault without requiring manual entry.
//
// Supported formats:
//
//   - ImportDotEnv – plain KEY=VALUE lines, optionally quoted, with # comments.
//   - ImportShell  – same as dotenv but lines may be prefixed with "export ".
//
// # Conflict Resolution
//
// When a key already exists in the target vault the caller chooses behaviour
// via the overwrite flag:
//
//   - overwrite=false: existing keys are skipped and counted in Skipped.
//   - overwrite=true:  existing keys are updated and counted in Overwritten.
//
// New keys are always counted in Added.
//
// # Example
//
//	res, err := vault.ImportFromFile(v, ".env.local", vault.ImportDotEnv, false)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("added %d, skipped %d\n", res.Added, res.Skipped)
package vault
