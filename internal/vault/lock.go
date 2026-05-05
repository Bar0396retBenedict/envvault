package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// LockRecord holds the lock state for a vault file.
type LockRecord struct {
	Locked    bool      `json:"locked"`
	LockedAt  time.Time `json:"locked_at,omitempty"`
	LockedBy  string    `json:"locked_by,omitempty"`
}

func lockFilePath(vaultPath string) string {
	base := vaultPath + ".lock"
	return filepath.Clean(base)
}

// LoadLockRecord reads the lock record for the given vault file.
// If no lock file exists, it returns an unlocked record.
func LoadLockRecord(vaultPath string) (LockRecord, error) {
	path := lockFilePath(vaultPath)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return LockRecord{Locked: false}, nil
	}
	if err != nil {
		return LockRecord{}, fmt.Errorf("load lock record: %w", err)
	}
	var rec LockRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return LockRecord{}, fmt.Errorf("parse lock record: %w", err)
	}
	return rec, nil
}

func saveLockRecord(vaultPath string, rec LockRecord) error {
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal lock record: %w", err)
	}
	path := lockFilePath(vaultPath)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("write lock record: %w", err)
	}
	return nil
}

// LockVault marks the vault as locked, preventing writes.
// lockedBy is an optional identifier (e.g. username or process).
func LockVault(vaultPath, lockedBy string) error {
	rec, err := LoadLockRecord(vaultPath)
	if err != nil {
		return err
	}
	if rec.Locked {
		return fmt.Errorf("vault is already locked by %q at %s", rec.LockedBy, rec.LockedAt.Format(time.RFC3339))
	}
	rec = LockRecord{
		Locked:   true,
		LockedAt: time.Now().UTC(),
		LockedBy: lockedBy,
	}
	return saveLockRecord(vaultPath, rec)
}

// UnlockVault removes the lock from the vault.
func UnlockVault(vaultPath string) error {
	rec, err := LoadLockRecord(vaultPath)
	if err != nil {
		return err
	}
	if !rec.Locked {
		return fmt.Errorf("vault is not locked")
	}
	return saveLockRecord(vaultPath, LockRecord{Locked: false})
}

// IsLocked returns true when the vault has an active lock.
func IsLocked(vaultPath string) (bool, error) {
	rec, err := LoadLockRecord(vaultPath)
	if err != nil {
		return false, err
	}
	return rec.Locked, nil
}
