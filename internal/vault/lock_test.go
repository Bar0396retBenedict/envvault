package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func makeLockVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "test.vault")
}

func TestLoadLockRecordMissing(t *testing.T) {
	vaultPath := makeLockVault(t)
	rec, err := LoadLockRecord(vaultPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Locked {
		t.Error("expected unlocked record when no lock file exists")
	}
}

func TestLockAndUnlock(t *testing.T) {
	vaultPath := makeLockVault(t)

	if err := LockVault(vaultPath, "alice"); err != nil {
		t.Fatalf("LockVault: %v", err)
	}

	locked, err := IsLocked(vaultPath)
	if err != nil {
		t.Fatalf("IsLocked: %v", err)
	}
	if !locked {
		t.Error("expected vault to be locked")
	}

	if err := UnlockVault(vaultPath); err != nil {
		t.Fatalf("UnlockVault: %v", err)
	}

	locked, err = IsLocked(vaultPath)
	if err != nil {
		t.Fatalf("IsLocked after unlock: %v", err)
	}
	if locked {
		t.Error("expected vault to be unlocked after UnlockVault")
	}
}

func TestLockAlreadyLocked(t *testing.T) {
	vaultPath := makeLockVault(t)

	if err := LockVault(vaultPath, "alice"); err != nil {
		t.Fatalf("first lock: %v", err)
	}
	if err := LockVault(vaultPath, "bob"); err == nil {
		t.Error("expected error when locking an already-locked vault")
	}
}

func TestUnlockNotLocked(t *testing.T) {
	vaultPath := makeLockVault(t)
	if err := UnlockVault(vaultPath); err == nil {
		t.Error("expected error when unlocking a vault that is not locked")
	}
}

func TestLockRecordLockedBy(t *testing.T) {
	vaultPath := makeLockVault(t)
	if err := LockVault(vaultPath, "ci-pipeline"); err != nil {
		t.Fatalf("LockVault: %v", err)
	}
	rec, err := LoadLockRecord(vaultPath)
	if err != nil {
		t.Fatalf("LoadLockRecord: %v", err)
	}
	if rec.LockedBy != "ci-pipeline" {
		t.Errorf("expected LockedBy=ci-pipeline, got %q", rec.LockedBy)
	}
	if rec.LockedAt.IsZero() {
		t.Error("expected non-zero LockedAt timestamp")
	}
}

func TestLockFilePermissions(t *testing.T) {
	vaultPath := makeLockVault(t)
	if err := LockVault(vaultPath, "ops"); err != nil {
		t.Fatalf("LockVault: %v", err)
	}
	info, err := os.Stat(lockFilePath(vaultPath))
	if err != nil {
		t.Fatalf("stat lock file: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("expected 0600 permissions, got %04o", perm)
	}
}
