// Package crypto provides AES-GCM encryption and decryption utilities
// for envvault encrypted configuration files.
//
// Usage:
//
//	key := crypto.DeriveKey("my-passphrase")
//
//	ciphertext, err := crypto.Encrypt(key, []byte("MY_SECRET=value"))
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	plaintext, err := crypto.Decrypt(key, ciphertext)
//	if err != nil {
//		log.Fatal(err)
//	}
//
// Keys are derived deterministically from a passphrase using SHA-256,
// producing a 32-byte key suitable for AES-256-GCM.
//
// Each encryption call uses a randomly generated nonce, which is prepended
// to the ciphertext. This ensures that encrypting the same plaintext twice
// produces different ciphertexts, providing semantic security.
package crypto
