package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func makeGroupVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "test.vault")
}

func TestLoadGroupRecordMissing(t *testing.T) {
	vp := makeGroupVault(t)
	rec, err := LoadGroupRecord(vp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rec.Groups) != 0 {
		t.Errorf("expected empty groups, got %v", rec.Groups)
	}
}

func TestAddToGroupNew(t *testing.T) {
	vp := makeGroupVault(t)
	if err := AddToGroup(vp, "backend", "DB_HOST"); err != nil {
		t.Fatalf("AddToGroup: %v", err)
	}
	keys, err := KeysForGroup(vp, "backend")
	if err != nil {
		t.Fatalf("KeysForGroup: %v", err)
	}
	if len(keys) != 1 || keys[0] != "DB_HOST" {
		t.Errorf("expected [DB_HOST], got %v", keys)
	}
}

func TestAddToGroupNoDuplicates(t *testing.T) {
	vp := makeGroupVault(t)
	_ = AddToGroup(vp, "backend", "DB_HOST")
	_ = AddToGroup(vp, "backend", "DB_HOST")
	keys, _ := KeysForGroup(vp, "backend")
	if len(keys) != 1 {
		t.Errorf("expected 1 key, got %d", len(keys))
	}
}

func TestAddToGroupSorted(t *testing.T) {
	vp := makeGroupVault(t)
	_ = AddToGroup(vp, "backend", "Z_KEY")
	_ = AddToGroup(vp, "backend", "A_KEY")
	keys, _ := KeysForGroup(vp, "backend")
	if keys[0] != "A_KEY" || keys[1] != "Z_KEY" {
		t.Errorf("expected sorted keys, got %v", keys)
	}
}

func TestRemoveFromGroup(t *testing.T) {
	vp := makeGroupVault(t)
	_ = AddToGroup(vp, "backend", "DB_HOST")
	_ = AddToGroup(vp, "backend", "DB_PORT")
	if err := RemoveFromGroup(vp, "backend", "DB_HOST"); err != nil {
		t.Fatalf("RemoveFromGroup: %v", err)
	}
	keys, _ := KeysForGroup(vp, "backend")
	if len(keys) != 1 || keys[0] != "DB_PORT" {
		t.Errorf("expected [DB_PORT], got %v", keys)
	}
}

func TestRemoveFromGroupDeletesEmpty(t *testing.T) {
	vp := makeGroupVault(t)
	_ = AddToGroup(vp, "backend", "DB_HOST")
	_ = RemoveFromGroup(vp, "backend", "DB_HOST")
	groups, _ := ListGroups(vp)
	if len(groups) != 0 {
		t.Errorf("expected no groups after removing last key, got %v", groups)
	}
}

func TestListGroups(t *testing.T) {
	vp := makeGroupVault(t)
	_ = AddToGroup(vp, "frontend", "API_URL")
	_ = AddToGroup(vp, "backend", "DB_HOST")
	groups, err := ListGroups(vp)
	if err != nil {
		t.Fatalf("ListGroups: %v", err)
	}
	if len(groups) != 2 || groups[0] != "backend" || groups[1] != "frontend" {
		t.Errorf("expected [backend frontend], got %v", groups)
	}
}

func TestGroupFilePermissions(t *testing.T) {
	vp := makeGroupVault(t)
	_ = AddToGroup(vp, "ops", "SECRET_KEY")
	path := groupFilePath(vp)
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}
