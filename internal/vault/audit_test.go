package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAuditLogMissing(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "test.vault")
	log, err := LoadAuditLog(vaultPath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(log.Entries) != 0 {
		t.Fatalf("expected empty log, got %d entries", len(log.Entries))
	}
}

func TestAppendAuditEntry(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "test.vault")

	if err := AppendAuditEntry(vaultPath, AuditSet, "API_KEY"); err != nil {
		t.Fatalf("AppendAuditEntry: %v", err)
	}
	if err := AppendAuditEntry(vaultPath, AuditDelete, "OLD_KEY"); err != nil {
		t.Fatalf("AppendAuditEntry: %v", err)
	}

	log, err := LoadAuditLog(vaultPath)
	if err != nil {
		t.Fatalf("LoadAuditLog: %v", err)
	}
	if len(log.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(log.Entries))
	}
	if log.Entries[0].Action != AuditSet || log.Entries[0].Key != "API_KEY" {
		t.Errorf("unexpected first entry: %+v", log.Entries[0])
	}
	if log.Entries[1].Action != AuditDelete || log.Entries[1].Key != "OLD_KEY" {
		t.Errorf("unexpected second entry: %+v", log.Entries[1])
	}
}

func TestAuditLogPath(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"/home/user/prod.vault", "/home/user/prod.audit.json"},
		{"local.vault", "local.audit.json"},
	}
	for _, c := range cases {
		got := auditLogPath(c.input)
		if got != c.want {
			t.Errorf("auditLogPath(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}

func TestAuditLogFilePermissions(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "test.vault")
	if err := AppendAuditEntry(vaultPath, AuditRotate, ""); err != nil {
		t.Fatalf("AppendAuditEntry: %v", err)
	}
	info, err := os.Stat(auditLogPath(vaultPath))
	if err != nil {
		t.Fatalf("stat audit log: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected mode 0600, got %v", info.Mode().Perm())
	}
}

func TestAuditEntryTimestampNonZero(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "test.vault")
	if err := AppendAuditEntry(vaultPath, AuditImport, "DB_URL"); err != nil {
		t.Fatalf("AppendAuditEntry: %v", err)
	}
	log, _ := LoadAuditLog(vaultPath)
	if log.Entries[0].Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}
