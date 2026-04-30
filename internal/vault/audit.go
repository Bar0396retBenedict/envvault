package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// AuditAction represents the type of operation performed on a vault.
type AuditAction string

const (
	AuditSet    AuditAction = "set"
	AuditDelete AuditAction = "delete"
	AuditRotate AuditAction = "rotate"
	AuditImport AuditAction = "import"
	AuditCopy   AuditAction = "copy"
)

// AuditEntry records a single operation performed on a vault.
type AuditEntry struct {
	Timestamp time.Time   `json:"timestamp"`
	Action    AuditAction `json:"action"`
	Key       string      `json:"key,omitempty"`
	VaultPath string      `json:"vault_path"`
}

// AuditLog holds a list of audit entries for a vault.
type AuditLog struct {
	Entries []AuditEntry `json:"entries"`
}

// auditLogPath returns the path to the audit log file for a given vault file.
func auditLogPath(vaultPath string) string {
	ext := filepath.Ext(vaultPath)
	base := vaultPath[:len(vaultPath)-len(ext)]
	return base + ".audit.json"
}

// LoadAuditLog loads the audit log for the given vault path.
// Returns an empty log if the file does not exist.
func LoadAuditLog(vaultPath string) (*AuditLog, error) {
	path := auditLogPath(vaultPath)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &AuditLog{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("audit: read log: %w", err)
	}
	var log AuditLog
	if err := json.Unmarshal(data, &log); err != nil {
		return nil, fmt.Errorf("audit: parse log: %w", err)
	}
	return &log, nil
}

// AppendAuditEntry adds an entry to the audit log and persists it.
func AppendAuditEntry(vaultPath string, action AuditAction, key string) error {
	log, err := LoadAuditLog(vaultPath)
	if err != nil {
		return err
	}
	log.Entries = append(log.Entries, AuditEntry{
		Timestamp: time.Now().UTC(),
		Action:    action,
		Key:       key,
		VaultPath: vaultPath,
	})
	return saveAuditLog(vaultPath, log)
}

func saveAuditLog(vaultPath string, log *AuditLog) error {
	data, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		return fmt.Errorf("audit: marshal log: %w", err)
	}
	path := auditLogPath(vaultPath)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("audit: write log: %w", err)
	}
	return nil
}
