package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func makeQuotaVault(t *testing.T) (string, *Vault) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.vault")
	v := New()
	if err := v.Save(path, "pass"); err != nil {
		t.Fatalf("save: %v", err)
	}
	return path, v
}

func TestLoadQuotaRecordMissing(t *testing.T) {
	path, _ := makeQuotaVault(t)
	rec, err := LoadQuotaRecord(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.MaxKeys != 0 || rec.MaxBytes != 0 {
		t.Errorf("expected zero record, got %+v", rec)
	}
}

func TestSaveAndLoadQuotaRecord(t *testing.T) {
	path, _ := makeQuotaVault(t)
	want := QuotaRecord{MaxKeys: 10, MaxBytes: 4096}
	if err := SaveQuotaRecord(path, want); err != nil {
		t.Fatalf("save: %v", err)
	}
	got, err := LoadQuotaRecord(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if got != want {
		t.Errorf("want %+v, got %+v", want, got)
	}
}

func TestQuotaFilePermissions(t *testing.T) {
	path, _ := makeQuotaVault(t)
	if err := SaveQuotaRecord(path, QuotaRecord{MaxKeys: 5}); err != nil {
		t.Fatalf("save: %v", err)
	}
	qPath := quotaFilePath(path)
	info, err := os.Stat(qPath)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}

func TestCheckQuotaNoViolations(t *testing.T) {
	_, v := makeQuotaVault(t)
	v.Set("KEY1", "val1")
	v.Set("KEY2", "val2")
	rec := QuotaRecord{MaxKeys: 5, MaxBytes: 1024}
	violations := CheckQuota(v, rec)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %v", violations)
	}
}

func TestCheckQuotaKeyLimitExceeded(t *testing.T) {
	_, v := makeQuotaVault(t)
	for i := 0; i < 5; i++ {
		v.Set(string(rune('A'+i))+"_KEY", "value")
	}
	rec := QuotaRecord{MaxKeys: 3}
	violations := CheckQuota(v, rec)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Field != "max_keys" {
		t.Errorf("expected max_keys violation, got %s", violations[0].Field)
	}
}

func TestCheckQuotaByteLimitExceeded(t *testing.T) {
	_, v := makeQuotaVault(t)
	v.Set("LARGE_KEY", "this is a fairly long value that pushes over the byte limit")
	rec := QuotaRecord{MaxBytes: 10}
	violations := CheckQuota(v, rec)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Field != "max_bytes" {
		t.Errorf("expected max_bytes violation, got %s", violations[0].Field)
	}
}

func TestCheckQuotaZeroLimitsIgnored(t *testing.T) {
	_, v := makeQuotaVault(t)
	for i := 0; i < 100; i++ {
		v.Set(string(rune('A'+i%26))+"_KEY", "val")
	}
	rec := QuotaRecord{} // zero means unlimited
	violations := CheckQuota(v, rec)
	if len(violations) != 0 {
		t.Errorf("expected no violations with zero limits, got %v", violations)
	}
}
