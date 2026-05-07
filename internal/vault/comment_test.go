package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func makeCommentVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "test.vault")
}

func TestLoadCommentRecordMissing(t *testing.T) {
	vp := makeCommentVault(t)
	rec, err := LoadCommentRecord(vp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rec.Comments) != 0 {
		t.Errorf("expected empty comments, got %v", rec.Comments)
	}
}

func TestSetAndGetComment(t *testing.T) {
	vp := makeCommentVault(t)
	if err := SetComment(vp, "API_KEY", "Production API key"); err != nil {
		t.Fatalf("SetComment: %v", err)
	}
	got, err := GetComment(vp, "API_KEY")
	if err != nil {
		t.Fatalf("GetComment: %v", err)
	}
	if got != "Production API key" {
		t.Errorf("expected 'Production API key', got %q", got)
	}
}

func TestSetCommentUpdatesExisting(t *testing.T) {
	vp := makeCommentVault(t)
	_ = SetComment(vp, "DB_PASS", "old comment")
	_ = SetComment(vp, "DB_PASS", "new comment")
	got, _ := GetComment(vp, "DB_PASS")
	if got != "new comment" {
		t.Errorf("expected 'new comment', got %q", got)
	}
}

func TestRemoveComment(t *testing.T) {
	vp := makeCommentVault(t)
	_ = SetComment(vp, "TOKEN", "some token")
	if err := RemoveComment(vp, "TOKEN"); err != nil {
		t.Fatalf("RemoveComment: %v", err)
	}
	got, _ := GetComment(vp, "TOKEN")
	if got != "" {
		t.Errorf("expected empty comment after removal, got %q", got)
	}
}

func TestGetCommentMissingKey(t *testing.T) {
	vp := makeCommentVault(t)
	got, err := GetComment(vp, "NONEXISTENT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("expected empty string for missing key, got %q", got)
	}
}

func TestListCommentsSorted(t *testing.T) {
	vp := makeCommentVault(t)
	_ = SetComment(vp, "Z_KEY", "last")
	_ = SetComment(vp, "A_KEY", "first")
	_ = SetComment(vp, "M_KEY", "middle")
	keys, comments, err := ListComments(vp)
	if err != nil {
		t.Fatalf("ListComments: %v", err)
	}
	if len(keys) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(keys))
	}
	if keys[0] != "A_KEY" || keys[1] != "M_KEY" || keys[2] != "Z_KEY" {
		t.Errorf("unexpected key order: %v", keys)
	}
	if comments[0] != "first" {
		t.Errorf("expected 'first', got %q", comments[0])
	}
}

func TestCommentFilePermissions(t *testing.T) {
	vp := makeCommentVault(t)
	_ = SetComment(vp, "KEY", "value")
	path := commentFilePath(vp)
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600 permissions, got %v", info.Mode().Perm())
	}
}
