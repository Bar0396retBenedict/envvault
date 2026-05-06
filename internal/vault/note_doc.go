// Package vault — note subsystem
//
// # Key Notes
//
// Notes allow operators to attach human-readable annotations to individual
// vault keys. Annotations are stored in a sidecar JSON file alongside the
// encrypted vault file and are never encrypted themselves — they are
// intended to be informational metadata only (e.g. "rotated on 2024-01-15",
// "used by the payments service").
//
// # File Location
//
// Given a vault at /path/to/prod.vault, the note record is stored at
// /path/to/.prod.vault.notes.json with permissions 0600.
//
// # Functions
//
//   - SetNote(vaultPath, key, note) — attach or replace a note for a key.
//   - RemoveNote(vaultPath, key)    — delete the note for a key.
//   - GetNote(vaultPath, key)       — retrieve the note entry for a key.
//   - LoadNoteRecord(vaultPath)     — load the full note record from disk.
package vault
