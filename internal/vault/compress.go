package vault

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

// CompressVault reads the vault file at path, compresses its contents using
// gzip, and writes the result to destPath. The vault must be loadable with
// passphrase before compression so that only a valid vault is archived.
//
// The compressed file is a raw gzip stream of the encrypted vault bytes and
// can be decompressed with DecompressVault.
func CompressVault(path, destPath, passphrase string) error {
	v, err := Load(path, passphrase)
	if err != nil {
		return fmt.Errorf("compress: load vault: %w", err)
	}

	// Re-save to a buffer so we capture the canonical encrypted bytes.
	tmpPath := destPath + ".tmp"
	if err := v.Save(tmpPath, passphrase); err != nil {
		return fmt.Errorf("compress: stage vault: %w", err)
	}

	raw, err := readFile(tmpPath)
	if err != nil {
		return fmt.Errorf("compress: read staged vault: %w", err)
	}
	_ = removeFile(tmpPath)

	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	if _, err := w.Write(raw); err != nil {
		return fmt.Errorf("compress: gzip write: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("compress: gzip close: %w", err)
	}

	if err := writeFileAtomic(destPath, buf.Bytes(), 0o600); err != nil {
		return fmt.Errorf("compress: write dest: %w", err)
	}
	return nil
}

// DecompressVault decompresses a gzip-compressed vault file at srcPath,
// writes the raw encrypted vault to destPath, and verifies it can be loaded
// with passphrase. Returns the loaded Vault on success.
func DecompressVault(srcPath, destPath, passphrase string) (*Vault, error) {
	compressed, err := readFile(srcPath)
	if err != nil {
		return nil, fmt.Errorf("decompress: read source: %w", err)
	}

	r, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, fmt.Errorf("decompress: gzip reader: %w", err)
	}
	defer r.Close()

	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("decompress: gzip read: %w", err)
	}

	if err := writeFileAtomic(destPath, raw, 0o600); err != nil {
		return nil, fmt.Errorf("decompress: write dest: %w", err)
	}

	v, err := Load(destPath, passphrase)
	if err != nil {
		return nil, fmt.Errorf("decompress: load vault: %w", err)
	}
	return v, nil
}
