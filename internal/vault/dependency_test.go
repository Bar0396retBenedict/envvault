package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func makeDependencyVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "test.vault")
}

func TestLoadDependencyRecordMissing(t *testing.T) {
	vp := makeDependencyVault(t)
	rec, err := LoadDependencyRecord(vp)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(rec.Deps) != 0 {
		t.Fatalf("expected empty deps, got %v", rec.Deps)
	}
}

func TestAddDependency(t *testing.T) {
	vp := makeDependencyVault(t)
	if err := AddDependency(vp, "DATABASE_URL", "DB_HOST"); err != nil {
		t.Fatalf("AddDependency: %v", err)
	}
	deps, err := ListDependencies(vp, "DATABASE_URL")
	if err != nil {
		t.Fatalf("ListDependencies: %v", err)
	}
	if len(deps) != 1 || deps[0] != "DB_HOST" {
		t.Fatalf("expected [DB_HOST], got %v", deps)
	}
}

func TestAddDependencyNoDuplicates(t *testing.T) {
	vp := makeDependencyVault(t)
	_ = AddDependency(vp, "A", "B")
	_ = AddDependency(vp, "A", "B")
	deps, _ := ListDependencies(vp, "A")
	if len(deps) != 1 {
		t.Fatalf("expected 1 dep, got %d", len(deps))
	}
}

func TestAddDependencySorted(t *testing.T) {
	vp := makeDependencyVault(t)
	_ = AddDependency(vp, "KEY", "Z_DEP")
	_ = AddDependency(vp, "KEY", "A_DEP")
	deps, _ := ListDependencies(vp, "KEY")
	if deps[0] != "A_DEP" || deps[1] != "Z_DEP" {
		t.Fatalf("expected sorted deps, got %v", deps)
	}
}

func TestRemoveDependency(t *testing.T) {
	vp := makeDependencyVault(t)
	_ = AddDependency(vp, "KEY", "DEP1")
	_ = AddDependency(vp, "KEY", "DEP2")
	if err := RemoveDependency(vp, "KEY", "DEP1"); err != nil {
		t.Fatalf("RemoveDependency: %v", err)
	}
	deps, _ := ListDependencies(vp, "KEY")
	if len(deps) != 1 || deps[0] != "DEP2" {
		t.Fatalf("expected [DEP2], got %v", deps)
	}
}

func TestRemoveDependencyCleansEmptyEntry(t *testing.T) {
	vp := makeDependencyVault(t)
	_ = AddDependency(vp, "KEY", "ONLY")
	_ = RemoveDependency(vp, "KEY", "ONLY")
	rec, _ := LoadDependencyRecord(vp)
	if _, ok := rec.Deps["KEY"]; ok {
		t.Fatal("expected KEY to be removed from deps map")
	}
}

func TestDependents(t *testing.T) {
	vp := makeDependencyVault(t)
	_ = AddDependency(vp, "APP_URL", "BASE_URL")
	_ = AddDependency(vp, "API_URL", "BASE_URL")
	_ = AddDependency(vp, "OTHER", "UNRELATED")
	dependents, err := Dependents(vp, "BASE_URL")
	if err != nil {
		t.Fatalf("Dependents: %v", err)
	}
	if len(dependents) != 2 {
		t.Fatalf("expected 2 dependents, got %v", dependents)
	}
	if dependents[0] != "API_URL" || dependents[1] != "APP_URL" {
		t.Fatalf("unexpected dependents order: %v", dependents)
	}
}

func TestDependencyFilePermissions(t *testing.T) {
	vp := makeDependencyVault(t)
	_ = AddDependency(vp, "K", "D")
	path := dependencyFilePath(vp)
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("expected 0600, got %v", info.Mode().Perm())
	}
}
