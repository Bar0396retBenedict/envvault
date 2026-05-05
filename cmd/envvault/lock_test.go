package main

import (
	"os"
	"path/filepath"
	"testing"

	"envvault/internal/vault"
)

func writeLockVault(t *testing.T, passphrase string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "env.vault")
	v := vault.New()
	v.Set("KEY", "value")
	if err := v.Save(path, passphrase); err != nil {
		t.Fatalf("save vault: %v", err)
	}
	return path
}

func TestRunLockMissingVault(t *testing.T) {
	err := runLock("/nonexistent/path.vault", "")
	if err == nil {
		t.Error("expected error for missing vault")
	}
}

func TestRunUnlockMissingVault(t *testing.T) {
	err := runUnlock("/nonexistent/path.vault")
	if err == nil {
		t.Error("expected error for missing vault")
	}
}

func TestRunLockSuccess(t *testing.T) {
	vaultPath := writeLockVault(t, "secret")

	if err := runLock(vaultPath, "alice"); err != nil {
		t.Fatalf("runLock: %v", err)
	}

	locked, err := vault.IsLocked(vaultPath)
	if err != nil {
		t.Fatalf("IsLocked: %v", err)
	}
	if !locked {
		t.Error("expected vault to be locked after runLock")
	}
}

func TestRunUnlockSuccess(t *testing.T) {
	vaultPath := writeLockVault(t, "secret")

	if err := vault.LockVault(vaultPath, "ci"); err != nil {
		t.Fatalf("setup lock: %v", err)
	}
	if err := runUnlock(vaultPath); err != nil {
		t.Fatalf("runUnlock: %v", err)
	}

	locked, err := vault.IsLocked(vaultPath)
	if err != nil {
		t.Fatalf("IsLocked: %v", err)
	}
	if locked {
		t.Error("expected vault to be unlocked after runUnlock")
	}
}

func TestRunLockAlreadyLocked(t *testing.T) {
	vaultPath := writeLockVault(t, "secret")

	if err := runLock(vaultPath, "alice"); err != nil {
		t.Fatalf("first lock: %v", err)
	}
	if err := runLock(vaultPath, "bob"); err == nil {
		t.Error("expected error locking an already-locked vault")
	}
}

func TestRunUnlockNotLocked(t *testing.T) {
	vaultPath := writeLockVault(t, "secret")
	if err := runUnlock(vaultPath); err == nil {
		t.Error("expected error unlocking a vault that is not locked")
	}
}

func TestLockFileCreated(t *testing.T) {
	vaultPath := writeLockVault(t, "secret")
	if err := runLock(vaultPath, ""); err != nil {
		t.Fatalf("runLock: %v", err)
	}
	lockPath := vaultPath + ".lock"
	if _, err := os.Stat(lockPath); os.IsNotExist(err) {
		t.Errorf("expected lock file to exist at %s", lockPath)
	}
}
