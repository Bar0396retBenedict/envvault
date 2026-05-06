package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func makeCompressVault(t *testing.T, pass string) (string, *Vault) {
	t.Helper()
	dir := t.TempDir()
	v := New()
	v.Set("DB_HOST", "localhost")
	v.Set("DB_PORT", "5432")
	v.Set("SECRET", "s3cr3t")
	path := filepath.Join(dir, "vault.env")
	if err := v.Save(path, pass); err != nil {
		t.Fatalf("save: %v", err)
	}
	return path, v
}

func TestCompressAndDecompress(t *testing.T) {
	dir := t.TempDir()
	src, _ := makeCompressVault(t, "passphrase")
	gzPath := filepath.Join(dir, "vault.env.gz")
	restored := filepath.Join(dir, "restored.env")

	if err := CompressVault(src, gzPath, "passphrase"); err != nil {
		t.Fatalf("CompressVault: %v", err)
	}

	if _, err := os.Stat(gzPath); err != nil {
		t.Fatalf("compressed file missing: %v", err)
	}

	v, err := DecompressVault(gzPath, restored, "passphrase")
	if err != nil {
		t.Fatalf("DecompressVault: %v", err)
	}

	for _, key := range []string{"DB_HOST", "DB_PORT", "SECRET"} {
		val, ok := v.Get(key)
		if !ok {
			t.Errorf("key %q missing after decompress", key)
		}
		_ = val
	}
}

func TestCompressVaultWrongPassphrase(t *testing.T) {
	src, _ := makeCompressVault(t, "correct")
	gzPath := filepath.Join(t.TempDir(), "vault.env.gz")

	err := CompressVault(src, gzPath, "wrong")
	if err == nil {
		t.Fatal("expected error for wrong passphrase, got nil")
	}
}

func TestDecompressWrongPassphrase(t *testing.T) {
	dir := t.TempDir()
	src, _ := makeCompressVault(t, "correct")
	gzPath := filepath.Join(dir, "vault.env.gz")

	if err := CompressVault(src, gzPath, "correct"); err != nil {
		t.Fatalf("compress: %v", err)
	}

	_, err := DecompressVault(gzPath, filepath.Join(dir, "out.env"), "wrong")
	if err == nil {
		t.Fatal("expected error for wrong passphrase on decompress, got nil")
	}
}

func TestCompressedFileSmallerOrReasonable(t *testing.T) {
	dir := t.TempDir()
	src, _ := makeCompressVault(t, "passphrase")
	gzPath := filepath.Join(dir, "vault.env.gz")

	if err := CompressVault(src, gzPath, "passphrase"); err != nil {
		t.Fatalf("CompressVault: %v", err)
	}

	origStat, _ := os.Stat(src)
	gzStat, _ := os.Stat(gzPath)

	// Compressed file should be non-zero.
	if gzStat.Size() == 0 {
		t.Error("compressed file is empty")
	}
	// Sanity: compressed should not be more than 3x the original.
	if gzStat.Size() > origStat.Size()*3 {
		t.Errorf("compressed size %d seems unexpectedly large vs original %d", gzStat.Size(), origStat.Size())
	}
}
