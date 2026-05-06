package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func makePinVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "test.vault")
}

func TestLoadPinRecordMissing(t *testing.T) {
	vp := makePinVault(t)
	rec, err := LoadPinRecord(vp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rec.Pins) != 0 {
		t.Errorf("expected empty pins, got %d", len(rec.Pins))
	}
}

func TestPinAndUnpin(t *testing.T) {
	vp := makePinVault(t)

	if err := PinKey(vp, "DB_PASSWORD", "do not rotate"); err != nil {
		t.Fatalf("PinKey: %v", err)
	}

	rec, err := LoadPinRecord(vp)
	if err != nil {
		t.Fatalf("LoadPinRecord: %v", err)
	}
	entry, ok := rec.Pins["DB_PASSWORD"]
	if !ok {
		t.Fatal("expected DB_PASSWORD to be pinned")
	}
	if entry.Note != "do not rotate" {
		t.Errorf("note mismatch: got %q", entry.Note)
	}
	if entry.PinnedAt.IsZero() {
		t.Error("expected non-zero PinnedAt")
	}

	if err := UnpinKey(vp, "DB_PASSWORD"); err != nil {
		t.Fatalf("UnpinKey: %v", err)
	}
	rec, _ = LoadPinRecord(vp)
	if _, ok := rec.Pins["DB_PASSWORD"]; ok {
		t.Error("expected DB_PASSWORD to be unpinned")
	}
}

func TestUnpinNotPinned(t *testing.T) {
	vp := makePinVault(t)
	err := UnpinKey(vp, "MISSING_KEY")
	if err == nil {
		t.Fatal("expected error unpinning non-existent key")
	}
}

func TestListPinnedKeys(t *testing.T) {
	vp := makePinVault(t)
	_ = PinKey(vp, "Z_KEY", "")
	_ = PinKey(vp, "A_KEY", "first")
	_ = PinKey(vp, "M_KEY", "middle")

	keys, _, err := ListPinnedKeys(vp)
	if err != nil {
		t.Fatalf("ListPinnedKeys: %v", err)
	}
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
	if keys[0] != "A_KEY" || keys[1] != "M_KEY" || keys[2] != "Z_KEY" {
		t.Errorf("keys not sorted: %v", keys)
	}
}

func TestPinFilePermissions(t *testing.T) {
	vp := makePinVault(t)
	_ = PinKey(vp, "SECRET", "")

	dir := filepath.Dir(vp)
	base := filepath.Base(vp)
	pinPath := filepath.Join(dir, "."+base+".pins.json")

	info, err := os.Stat(pinPath)
	if err != nil {
		t.Fatalf("stat pin file: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}
