package vault

import (
	"testing"
)

func TestAddTagNew(t *testing.T) {
	var tags Tags
	tags = AddTag(tags, "prod", "DB_URL")
	keys, err := KeysForTag(tags, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 1 || keys[0] != "DB_URL" {
		t.Fatalf("expected [DB_URL], got %v", keys)
	}
}

func TestAddTagNoDuplicates(t *testing.T) {
	tags := make(Tags)
	tags = AddTag(tags, "prod", "API_KEY")
	tags = AddTag(tags, "prod", "API_KEY")
	keys, _ := KeysForTag(tags, "prod")
	if len(keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(keys))
	}
}

func TestAddTagSorted(t *testing.T) {
	tags := make(Tags)
	tags = AddTag(tags, "staging", "Z_VAR")
	tags = AddTag(tags, "staging", "A_VAR")
	tags = AddTag(tags, "staging", "M_VAR")
	keys, _ := KeysForTag(tags, "staging")
	expected := []string{"A_VAR", "M_VAR", "Z_VAR"}
	for i, k := range keys {
		if k != expected[i] {
			t.Fatalf("expected %v, got %v", expected, keys)
		}
	}
}

func TestRemoveTagKey(t *testing.T) {
	tags := make(Tags)
	tags = AddTag(tags, "prod", "DB_URL")
	tags = AddTag(tags, "prod", "API_KEY")
	tags = RemoveTag(tags, "prod", "DB_URL")
	keys, _ := KeysForTag(tags, "prod")
	if len(keys) != 1 || keys[0] != "API_KEY" {
		t.Fatalf("expected [API_KEY], got %v", keys)
	}
}

func TestRemoveTagDeletesEmptyTag(t *testing.T) {
	tags := make(Tags)
	tags = AddTag(tags, "dev", "SECRET")
	tags = RemoveTag(tags, "dev", "SECRET")
	if _, ok := tags["dev"]; ok {
		t.Fatal("expected tag 'dev' to be deleted")
	}
}

func TestKeysForTagMissing(t *testing.T) {
	tags := make(Tags)
	_, err := KeysForTag(tags, "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing tag")
	}
}

func TestListTags(t *testing.T) {
	tags := make(Tags)
	tags = AddTag(tags, "prod", "A")
	tags = AddTag(tags, "dev", "B")
	tags = AddTag(tags, "staging", "C")
	names := ListTags(tags)
	expected := []string{"dev", "prod", "staging"}
	for i, n := range names {
		if n != expected[i] {
			t.Fatalf("expected %v, got %v", expected, names)
		}
	}
}

func TestAddTagIgnoresEmpty(t *testing.T) {
	tags := make(Tags)
	tags = AddTag(tags, "", "KEY")
	tags = AddTag(tags, "prod", "")
	if len(tags) != 0 {
		t.Fatalf("expected empty tags map, got %v", tags)
	}
}
