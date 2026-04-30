// Package vault — diff.go
//
// # Vault Diff
//
// The diff module provides utilities for comparing two Vault instances and
// producing a structured list of changes. This is useful when syncing
// environment variables between environments (local → staging → production)
// and you need a human-readable summary before committing changes.
//
// # Usage
//
//	src, _ := vault.Load("local.env.vault", passphrase)
//	dst, _ := vault.Load("staging.env.vault", passphrase)
//
//	changes := vault.Diff(src, dst)
//	for _, c := range changes {
//		switch c.Kind {
//		case vault.ChangeAdded:
//			fmt.Printf("+ %s=%s\n", c.Key, c.NewValue)
//		case vault.ChangeRemoved:
//			fmt.Printf("- %s\n", c.Key)
//		case vault.ChangeUpdated:
//			fmt.Printf("~ %s: %s → %s\n", c.Key, c.OldValue, c.NewValue)
//		}
//	}
//
// Changes are always returned in lexicographic key order for deterministic
// output in CLI summaries and logs.
package vault
