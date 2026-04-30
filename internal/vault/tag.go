package vault

import (
	"fmt"
	"sort"
	"strings"
)

// Tags holds a mapping from tag name to a list of env var keys.
type Tags map[string][]string

// AddTag associates key with the given tag in the vault's tag index.
// If the key is already present under that tag it is not duplicated.
func AddTag(tags Tags, tag, key string) Tags {
	if tags == nil {
		tags = make(Tags)
	}
	tag = strings.TrimSpace(tag)
	key = strings.TrimSpace(key)
	if tag == "" || key == "" {
		return tags
	}
	for _, existing := range tags[tag] {
		if existing == key {
			return tags
		}
	}
	tags[tag] = append(tags[tag], key)
	sort.Strings(tags[tag])
	return tags
}

// RemoveTag removes key from the given tag. If the tag becomes empty it is
// deleted from the map entirely.
func RemoveTag(tags Tags, tag, key string) Tags {
	if tags == nil {
		return tags
	}
	keys := tags[tag]
	updated := keys[:0]
	for _, k := range keys {
		if k != key {
			updated = append(updated, k)
		}
	}
	if len(updated) == 0 {
		delete(tags, tag)
	} else {
		tags[tag] = updated
	}
	return tags
}

// KeysForTag returns the list of keys associated with tag, or an error if the
// tag does not exist.
func KeysForTag(tags Tags, tag string) ([]string, error) {
	keys, ok := tags[tag]
	if !ok {
		return nil, fmt.Errorf("tag %q not found", tag)
	}
	out := make([]string, len(keys))
	copy(out, keys)
	return out, nil
}

// ListTags returns all tag names in sorted order.
func ListTags(tags Tags) []string {
	names := make([]string, 0, len(tags))
	for t := range tags {
		names = append(names, t)
	}
	sort.Strings(names)
	return names
}
