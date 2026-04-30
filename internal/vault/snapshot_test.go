package vault

import (
	"path/filepath"
	"testing"
)

func makeSnapshotVault(t *testing.T, pass string, entries map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.vault")
	v := New(path, pass)
	for k, val := range entries {
		v.Set(k, val)
	}
	if err := v.Save(); err != nil {
		t.Fatalf("save vault: %v", err)
	}
	return path
}

func TestTakeSnapshot(t *testing.T) {
	path := makeSnapshotVault(t, "secret", map[string]string{"FOO": "bar", "BAZ": "qux"})
	snap, err := TakeSnapshot(path, "secret", "before-deploy")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap.Label != "before-deploy" {
		t.Errorf("expected label 'before-deploy', got %q", snap.Label)
	}
	if snap.Entries["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", snap.Entries["FOO"])
	}
	if snap.CreatedAt.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestListSnapshots(t *testing.T) {
	path := makeSnapshotVault(t, "secret", map[string]string{"KEY": "val"})

	if _, err := TakeSnapshot(path, "secret", "snap1"); err != nil {
		t.Fatalf("take snap1: %v", err)
	}
	if _, err := TakeSnapshot(path, "secret", "snap2"); err != nil {
		t.Fatalf("take snap2: %v", err)
	}

	snaps, err := ListSnapshots(path)
	if err != nil {
		t.Fatalf("list snapshots: %v", err)
	}
	if len(snaps) != 2 {
		t.Errorf("expected 2 snapshots, got %d", len(snaps))
	}
}

func TestListSnapshotsMissing(t *testing.T) {
	path := makeSnapshotVault(t, "secret", map[string]string{})
	snaps, err := ListSnapshots(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snaps != nil {
		t.Errorf("expected nil, got %v", snaps)
	}
}

func TestTakeSnapshotWrongPassphrase(t *testing.T) {
	path := makeSnapshotVault(t, "correct", map[string]string{"X": "1"})
	_, err := TakeSnapshot(path, "wrong", "label")
	if err == nil {
		t.Error("expected error for wrong passphrase")
	}
}

func TestSanitizeLabel(t *testing.T) {
	cases := []struct{ in, want string }{
		{"before-deploy", "before-deploy"},
		{"release 1.0", "release_1_0"},
		{"v2/prod", "v2_prod"},
	}
	for _, c := range cases {
		got := sanitizeLabel(c.in)
		if got != c.want {
			t.Errorf("sanitizeLabel(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
