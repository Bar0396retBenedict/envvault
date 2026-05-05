package vault

import (
	"testing"
)

func makeSearchVault(t *testing.T) *Vault {
	t.Helper()
	v := New()
	v.Set("DATABASE_URL", "postgres://localhost/mydb")
	v.Set("DATABASE_POOL", "10")
	v.Set("REDIS_URL", "redis://localhost:6379")
	v.Set("APP_SECRET", "supersecret")
	v.Set("APP_DEBUG", "true")
	return v
}

func TestSearchByKeyExact(t *testing.T) {
	v := makeSearchVault(t)
	results := Search(v, "DATABASE_URL", SearchOptions{})
	if len(results) != 1 || results[0].Key != "DATABASE_URL" {
		t.Fatalf("expected 1 result for DATABASE_URL, got %v", results)
	}
}

func TestSearchByKeyPartial(t *testing.T) {
	v := makeSearchVault(t)
	results := Search(v, "DATABASE", SearchOptions{})
	if len(results) != 2 {
		t.Fatalf("expected 2 results for DATABASE, got %d", len(results))
	}
	if results[0].Key != "DATABASE_POOL" || results[1].Key != "DATABASE_URL" {
		t.Errorf("unexpected order: %v", results)
	}
}

func TestSearchCaseInsensitive(t *testing.T) {
	v := makeSearchVault(t)
	results := Search(v, "database", SearchOptions{CaseSensitive: false})
	if len(results) != 2 {
		t.Fatalf("expected 2 results (case-insensitive), got %d", len(results))
	}
}

func TestSearchCaseSensitiveNoMatch(t *testing.T) {
	v := makeSearchVault(t)
	results := Search(v, "database", SearchOptions{CaseSensitive: true})
	if len(results) != 0 {
		t.Fatalf("expected 0 results with case-sensitive search, got %d", len(results))
	}
}

func TestSearchMatchValue(t *testing.T) {
	v := makeSearchVault(t)
	results := Search(v, "localhost", SearchOptions{MatchValue: true})
	if len(results) != 2 {
		t.Fatalf("expected 2 results matching value 'localhost', got %d", len(results))
	}
}

func TestSearchNoResults(t *testing.T) {
	v := makeSearchVault(t)
	results := Search(v, "NONEXISTENT", SearchOptions{})
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestSearchEmptyQuery(t *testing.T) {
	v := makeSearchVault(t)
	results := Search(v, "", SearchOptions{})
	if len(results) != 5 {
		t.Fatalf("empty query should match all keys, got %d", len(results))
	}
}
