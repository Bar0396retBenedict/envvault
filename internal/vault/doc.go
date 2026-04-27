// Package vault implements the core vault data model for envvault.
//
// A Vault is an encrypted container that stores environment variables
// as key-value pairs. Vaults are persisted to disk as encrypted binary
// files using AES-256-GCM authenticated encryption (via the crypto package).
//
// # Basic Usage
//
// Create and populate a new vault:
//
//	v := vault.New()
//	v.Set("DATABASE_URL", "postgres://localhost/mydb")
//	v.Set("API_KEY", "my-secret-key")
//
// Save to disk with a passphrase:
//
//	if err := v.Save("production.vault", passphrase); err != nil {
//		log.Fatal(err)
//	}
//
// Load an existing vault:
//
//	v, err := vault.Load("production.vault", passphrase)
//	if err != nil {
//		log.Fatal(err)
//	}
//
// Export to .env format:
//
//	if err := v.Export(os.Stdout, vault.FormatDotEnv); err != nil {
//		log.Fatal(err)
//	}
//
// # File Format
//
// Vault files contain JSON-marshalled Vault structs encrypted with
// AES-256-GCM. The passphrase is stretched using scrypt before use.
// Files are written with mode 0600 to restrict access.
package vault
