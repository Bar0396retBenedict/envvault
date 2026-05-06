package vault

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func makeAccessVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "test.vault")
}

func TestLoadAccessRecordMissing(t *testing.T) {
	vp := makeAccessVault(t)
	rec, err := LoadAccessRecord(vp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rec.Entries) != 0 {
		t.Errorf("expected empty entries, got %d", len(rec.Entries))
	}
}

func TestRecordRead(t *testing.T) {
	vp := makeAccessVault(t)
	if err := RecordRead(vp, "API_KEY"); err != nil {
		t.Fatalf("RecordRead: %v", err)
	}
	if err := RecordRead(vp, "API_KEY"); err != nil {
		t.Fatalf("RecordRead second: %v", err)
	}
	rec, err := LoadAccessRecord(vp)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	e := rec.Entries["API_KEY"]
	if e == nil {
		t.Fatal("entry not found")
	}
	if e.ReadCount != 2 {
		t.Errorf("expected ReadCount=2, got %d", e.ReadCount)
	}
	if e.LastRead.IsZero() {
		t.Error("LastRead should not be zero")
	}
}

func TestRecordWrite(t *testing.T) {
	vp := makeAccessVault(t)
	if err := RecordWrite(vp, "DB_PASS"); err != nil {
		t.Fatalf("RecordWrite: %v", err)
	}
	rec, err := LoadAccessRecord(vp)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	e := rec.Entries["DB_PASS"]
	if e == nil {
		t.Fatal("entry not found")
	}
	if e.WriteCount != 1 {
		t.Errorf("expected WriteCount=1, got %d", e.WriteCount)
	}
	if e.LastWrite.IsZero() {
		t.Error("LastWrite should not be zero")
	}
}

func TestListAccessEntries(t *testing.T) {
	vp := makeAccessVault(t)
	keys := []string{"Z_KEY", "A_KEY", "M_KEY"}
	for _, k := range keys {
		_ = RecordRead(vp, k)
	}
	entries, err := ListAccessEntries(vp)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Key != "A_KEY" || entries[1].Key != "M_KEY" || entries[2].Key != "Z_KEY" {
		t.Errorf("entries not sorted: %v", entries)
	}
}

func TestAccessFilePermissions(t *testing.T) {
	vp := makeAccessVault(t)
	_ = RecordRead(vp, "SECRET")
	p := accessFilePath(vp)
	info, err := os.Stat(p)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %o", info.Mode().Perm())
	}
}

func TestAccessTimestampRecency(t *testing.T) {
	vp := makeAccessVault(t)
	before := time.Now().UTC().Add(-time.Second)
	_ = RecordWrite(vp, "TOKEN")
	after := time.Now().UTC().Add(time.Second)
	rec, _ := LoadAccessRecord(vp)
	e := rec.Entries["TOKEN"]
	if e.LastWrite.Before(before) || e.LastWrite.After(after) {
		t.Errorf("LastWrite %v not in expected range", e.LastWrite)
	}
}
