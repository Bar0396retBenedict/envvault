package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadProfilesMissing(t *testing.T) {
	dir := t.TempDir()
	cfg, err := LoadProfiles(dir)
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(cfg.Profiles) != 0 {
		t.Errorf("expected empty profiles, got %d", len(cfg.Profiles))
	}
}

func TestSaveAndLoadProfiles(t *testing.T) {
	dir := t.TempDir()
	cfg := &ProfileConfig{}
	cfg.AddProfile("local", filepath.Join(dir, "local.vault"))
	cfg.AddProfile("staging", filepath.Join(dir, "staging.vault"))

	if err := SaveProfiles(dir, cfg); err != nil {
		t.Fatalf("SaveProfiles: %v", err)
	}

	loaded, err := LoadProfiles(dir)
	if err != nil {
		t.Fatalf("LoadProfiles: %v", err)
	}
	if len(loaded.Profiles) != 2 {
		t.Errorf("expected 2 profiles, got %d", len(loaded.Profiles))
	}
}

func TestAddProfileUpdatesExisting(t *testing.T) {
	cfg := &ProfileConfig{}
	cfg.AddProfile("local", "/old/path")
	cfg.AddProfile("local", "/new/path")

	if len(cfg.Profiles) != 1 {
		t.Errorf("expected 1 profile after update, got %d", len(cfg.Profiles))
	}
	if cfg.Profiles[0].VaultPath != "/new/path" {
		t.Errorf("expected updated path, got %s", cfg.Profiles[0].VaultPath)
	}
}

func TestGetProfile(t *testing.T) {
	cfg := &ProfileConfig{}
	cfg.AddProfile("production", "/prod.vault")

	p, err := cfg.GetProfile("production")
	if err != nil {
		t.Fatalf("GetProfile: %v", err)
	}
	if p.VaultPath != "/prod.vault" {
		t.Errorf("unexpected vault path: %s", p.VaultPath)
	}

	_, err = cfg.GetProfile("nonexistent")
	if err == nil {
		t.Error("expected error for missing profile")
	}
}

func TestSetActive(t *testing.T) {
	cfg := &ProfileConfig{}
	cfg.AddProfile("staging", "/staging.vault")

	if err := cfg.SetActive("staging"); err != nil {
		t.Fatalf("SetActive: %v", err)
	}
	if cfg.Active != "staging" {
		t.Errorf("expected active=staging, got %s", cfg.Active)
	}

	if err := cfg.SetActive("missing"); err == nil {
		t.Error("expected error when setting missing profile active")
	}
}

func TestProfilesFilePermissions(t *testing.T) {
	dir := t.TempDir()
	cfg := &ProfileConfig{}
	cfg.AddProfile("local", "/local.vault")

	if err := SaveProfiles(dir, cfg); err != nil {
		t.Fatalf("SaveProfiles: %v", err)
	}

	path := filepath.Join(dir, ".envvault", "profiles.json")
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat profiles file: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected permissions 0600, got %v", info.Mode().Perm())
	}
}
