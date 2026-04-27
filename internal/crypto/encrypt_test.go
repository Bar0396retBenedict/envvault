package crypto

import (
	"bytes"
	"testing"
)

func TestDeriveKey(t *testing.T) {
	key := DeriveKey("my-secret-passphrase")
	if len(key) != 32 {
		t.Fatalf("expected key length 32, got %d", len(key))
	}

	// Same passphrase should produce same key
	key2 := DeriveKey("my-secret-passphrase")
	if !bytes.Equal(key, key2) {
		t.Fatal("expected deterministic key derivation")
	}

	// Different passphrase should produce different key
	key3 := DeriveKey("other-passphrase")
	if bytes.Equal(key, key3) {
		t.Fatal("different passphrases should produce different keys")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	key := DeriveKey("test-passphrase")
	plaintext := []byte("DATABASE_URL=postgres://localhost/mydb\nAPI_KEY=supersecret")

	ciphertext, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	if bytes.Equal(ciphertext, plaintext) {
		t.Fatal("ciphertext should not equal plaintext")
	}

	decrypted, err := Decrypt(key, ciphertext)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Fatalf("decrypted text does not match original: got %q", decrypted)
	}
}

func TestEncryptProducesUniqueOutput(t *testing.T) {
	key := DeriveKey("test-passphrase")
	plaintext := []byte("SECRET=value")

	c1, _ := Encrypt(key, plaintext)
	c2, _ := Encrypt(key, plaintext)

	// Due to random nonce, two encryptions of the same data should differ
	if bytes.Equal(c1, c2) {
		t.Fatal("expected different ciphertexts for same plaintext due to random nonce")
	}
}

func TestDecryptWithWrongKey(t *testing.T) {
	key := DeriveKey("correct-passphrase")
	wrongKey := DeriveKey("wrong-passphrase")
	plaintext := []byte("SECRET=value")

	ciphertext, _ := Encrypt(key, plaintext)

	_, err := Decrypt(wrongKey, ciphertext)
	if err == nil {
		t.Fatal("expected error when decrypting with wrong key")
	}
}

func TestDecryptShortCiphertext(t *testing.T) {
	key := DeriveKey("test-passphrase")
	_, err := Decrypt(key, []byte("short"))
	if err == nil {
		t.Fatal("expected error for short ciphertext")
	}
}
