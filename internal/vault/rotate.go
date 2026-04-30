package vault

import (
	"fmt"
	"time"
)

// RotateRecord holds metadata about a passphrase rotation event.
type RotateRecord struct {
	RotatedAt time.Time `json:"rotated_at"`
	Profile   string    `json:"profile"`
}

// Rotate re-encrypts the vault with a new passphrase, replacing the old one.
// It loads all secrets using the old passphrase, then saves them under the
// new passphrase. Returns an error if decryption or re-encryption fails.
func Rotate(path, oldPassphrase, newPassphrase string) (*RotateRecord, error) {
	if oldPassphrase == "" {
		return nil, fmt.Errorf("rotate: old passphrase must not be empty")
	}
	if newPassphrase == "" {
		return nil, fmt.Errorf("rotate: new passphrase must not be empty")
	}
	if oldPassphrase == newPassphrase {
		return nil, fmt.Errorf("rotate: new passphrase must differ from old passphrase")
	}

	v, err := New(path, oldPassphrase)
	if err != nil {
		return nil, fmt.Errorf("rotate: load vault: %w", err)
	}

	// Snapshot all current secrets.
	secrets := make(map[string]string, len(v.data))
	for k, val := range v.data {
		secrets[k] = val
	}

	// Create a new vault at the same path with the new passphrase.
	nv, err := New(path, newPassphrase)
	if err != nil {
		return nil, fmt.Errorf("rotate: init new vault: %w", err)
	}

	for k, val := range secrets {
		if err := nv.Set(k, val); err != nil {
			return nil, fmt.Errorf("rotate: set key %q: %w", k, err)
		}
	}

	if err := nv.Save(); err != nil {
		return nil, fmt.Errorf("rotate: save vault: %w", err)
	}

	return &RotateRecord{
		RotatedAt: time.Now().UTC(),
		Profile:   path,
	}, nil
}
