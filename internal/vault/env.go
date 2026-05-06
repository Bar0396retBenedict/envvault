package vault

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// EnvRecord maps environment variable names to their vault source paths.
type EnvRecord struct {
	Bindings map[string]string `json:"bindings"` // envKey -> vaultPath
}

func envFilePath(vaultPath string) string {
	dir := filepath.Dir(vaultPath)
	base := strings.TrimSuffix(filepath.Base(vaultPath), filepath.Ext(vaultPath))
	return filepath.Join(dir, base+".env.json")
}

// LoadEnvRecord reads the env binding record for the given vault.
// Returns an empty record if the file does not exist.
func LoadEnvRecord(vaultPath string) (EnvRecord, error) {
	rec := EnvRecord{Bindings: make(map[string]string)}
	path := envFilePath(vaultPath)
	if err := loadJSON(path, &rec); err != nil {
		if os.IsNotExist(err) {
			return rec, nil
		}
		return rec, fmt.Errorf("load env record: %w", err)
	}
	if rec.Bindings == nil {
		rec.Bindings = make(map[string]string)
	}
	return rec, nil
}

func saveEnvRecord(vaultPath string, rec EnvRecord) error {
	path := envFilePath(vaultPath)
	if err := saveJSON(path, rec); err != nil {
		return fmt.Errorf("save env record: %w", err)
	}
	return nil
}

// BindEnvVar associates an OS environment variable name with a key in the vault.
// When the binding is applied, the vault key's value is exported to the OS env.
func BindEnvVar(vaultPath, envKey, vaultKey string) error {
	if envKey == "" {
		return fmt.Errorf("env key must not be empty")
	}
	if vaultKey == "" {
		return fmt.Errorf("vault key must not be empty")
	}
	rec, err := LoadEnvRecord(vaultPath)
	if err != nil {
		return err
	}
	rec.Bindings[envKey] = vaultKey
	return saveEnvRecord(vaultPath, rec)
}

// UnbindEnvVar removes the binding for the given OS environment variable.
func UnbindEnvVar(vaultPath, envKey string) error {
	rec, err := LoadEnvRecord(vaultPath)
	if err != nil {
		return err
	}
	if _, ok := rec.Bindings[envKey]; !ok {
		return fmt.Errorf("env var %q is not bound", envKey)
	}
	delete(rec.Bindings, envKey)
	return saveEnvRecord(vaultPath, rec)
}

// ApplyEnvBindings loads the vault and sets OS environment variables
// for all configured bindings.
func ApplyEnvBindings(vaultPath, passphrase string) ([]string, error) {
	v, err := Load(vaultPath, passphrase)
	if err != nil {
		return nil, fmt.Errorf("apply env bindings: %w", err)
	}
	rec, err := LoadEnvRecord(vaultPath)
	if err != nil {
		return nil, err
	}
	var applied []string
	for envKey, vaultKey := range rec.Bindings {
		val, ok := v.Get(vaultKey)
		if !ok {
			return nil, fmt.Errorf("vault key %q not found for env binding %q", vaultKey, envKey)
		}
		if err := os.Setenv(envKey, val); err != nil {
			return nil, fmt.Errorf("setenv %q: %w", envKey, err)
		}
		applied = append(applied, envKey)
	}
	return applied, nil
}
