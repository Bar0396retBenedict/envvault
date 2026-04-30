package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Profile represents a named environment target (e.g. local, staging, production).
type Profile struct {
	Name     string `json:"name"`
	VaultPath string `json:"vault_path"`
}

// ProfileConfig holds all registered profiles for the project.
type ProfileConfig struct {
	Profiles []Profile `json:"profiles"`
	Active   string    `json:"active"`
}

// profileConfigPath returns the path to the profiles config file.
func profileConfigPath(dir string) string {
	return filepath.Join(dir, ".envvault", "profiles.json")
}

// LoadProfiles reads the profile configuration from the given directory.
// Returns an empty ProfileConfig if the file does not exist.
func LoadProfiles(dir string) (*ProfileConfig, error) {
	path := profileConfigPath(dir)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &ProfileConfig{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("load profiles: %w", err)
	}
	var cfg ProfileConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse profiles: %w", err)
	}
	return &cfg, nil
}

// SaveProfiles writes the profile configuration to the given directory.
func SaveProfiles(dir string, cfg *ProfileConfig) error {
	path := profileConfigPath(dir)
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("create profiles dir: %w", err)
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal profiles: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("write profiles: %w", err)
	}
	return nil
}

// AddProfile adds or updates a profile in the config.
func (c *ProfileConfig) AddProfile(name, vaultPath string) {
	for i, p := range c.Profiles {
		if p.Name == name {
			c.Profiles[i].VaultPath = vaultPath
			return
		}
	}
	c.Profiles = append(c.Profiles, Profile{Name: name, VaultPath: vaultPath})
}

// GetProfile returns the profile with the given name, or an error if not found.
func (c *ProfileConfig) GetProfile(name string) (Profile, error) {
	for _, p := range c.Profiles {
		if p.Name == name {
			return p, nil
		}
	}
	return Profile{}, fmt.Errorf("profile %q not found", name)
}

// SetActive marks the named profile as the active one.
func (c *ProfileConfig) SetActive(name string) error {
	for _, p := range c.Profiles {
		if p.Name == name {
			c.Active = name
			return nil
		}
	}
	return fmt.Errorf("profile %q not found", name)
}
