// Package vault provides snapshot functionality for capturing point-in-time
// copies of a vault's decrypted contents.
//
// # Overview
//
// Snapshots allow users to save the current state of a vault before making
// bulk changes, rotating passphrases, or deploying to a new environment.
// Each snapshot is stored as a JSON file inside a hidden `.snapshots`
// directory adjacent to the vault file.
//
// # Storage Layout
//
// Given a vault at `/path/to/prod.vault`, snapshots are written to:
//
//	/path/to/.snapshots/prod/<unix_nano>_<label>.json
//
// Files are named with a nanosecond Unix timestamp prefix so that directory
// listing naturally returns them in chronological order without additional
// sorting.
//
// # Security
//
// Snapshot files are written with mode 0600 and their parent directory with
// mode 0700. The entries stored inside a snapshot are plaintext, so the
// `.snapshots` directory should be treated with the same care as the vault
// file itself and should be excluded from version control.
package vault
