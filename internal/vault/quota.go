package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// QuotaRecord holds per-vault key count and size limits.
type QuotaRecord struct {
	MaxKeys  int   `json:"max_keys,omitempty"`
	MaxBytes int64 `json:"max_bytes,omitempty"`
}

// QuotaViolation describes a single quota breach.
type QuotaViolation struct {
	Field   string
	Limit   int64
	Actual  int64
	Message string
}

func quotaFilePath(vaultPath string) string {
	dir := filepath.Dir(vaultPath)
	base := filepath.Base(vaultPath)
	return filepath.Join(dir, "."+base+".quota.json")
}

// LoadQuotaRecord reads the quota config for the given vault file.
// Returns a zero-value record if the file does not exist.
func LoadQuotaRecord(vaultPath string) (QuotaRecord, error) {
	path := quotaFilePath(vaultPath)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return QuotaRecord{}, nil
	}
	if err != nil {
		return QuotaRecord{}, fmt.Errorf("quota: read: %w", err)
	}
	var rec QuotaRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return QuotaRecord{}, fmt.Errorf("quota: parse: %w", err)
	}
	return rec, nil
}

// SaveQuotaRecord persists a QuotaRecord next to the vault file.
func SaveQuotaRecord(vaultPath string, rec QuotaRecord) error {
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return fmt.Errorf("quota: marshal: %w", err)
	}
	path := quotaFilePath(vaultPath)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("quota: write: %w", err)
	}
	return nil
}

// CheckQuota validates the vault against its quota record.
// Returns a slice of violations (empty means all limits satisfied).
func CheckQuota(v *Vault, rec QuotaRecord) []QuotaViolation {
	var violations []QuotaViolation

	keys := v.Keys()
	keyCount := int64(len(keys))

	if rec.MaxKeys > 0 && keyCount > int64(rec.MaxKeys) {
		violations = append(violations, QuotaViolation{
			Field:   "max_keys",
			Limit:   int64(rec.MaxKeys),
			Actual:  keyCount,
			Message: fmt.Sprintf("key count %d exceeds limit %d", keyCount, rec.MaxKeys),
		})
	}

	if rec.MaxBytes > 0 {
		var total int64
		for _, k := range keys {
			val, _ := v.Get(k)
			total += int64(len(k)) + int64(len(val))
		}
		if total > rec.MaxBytes {
			violations = append(violations, QuotaViolation{
				Field:   "max_bytes",
				Limit:   rec.MaxBytes,
				Actual:  total,
				Message: fmt.Sprintf("data size %d bytes exceeds limit %d", total, rec.MaxBytes),
			})
		}
	}

	return violations
}
