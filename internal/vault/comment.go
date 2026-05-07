package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// CommentRecord maps vault keys to their associated comments.
type CommentRecord struct {
	Comments map[string]string `json:"comments"`
}

func commentFilePath(vaultPath string) string {
	dir := filepath.Dir(vaultPath)
	base := filepath.Base(vaultPath)
	return filepath.Join(dir, "."+base+".comments.json")
}

// LoadCommentRecord loads the comment record for the given vault file.
// If the file does not exist, an empty record is returned.
func LoadCommentRecord(vaultPath string) (CommentRecord, error) {
	path := commentFilePath(vaultPath)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return CommentRecord{Comments: make(map[string]string)}, nil
	}
	if err != nil {
		return CommentRecord{}, fmt.Errorf("read comment record: %w", err)
	}
	var rec CommentRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return CommentRecord{}, fmt.Errorf("parse comment record: %w", err)
	}
	if rec.Comments == nil {
		rec.Comments = make(map[string]string)
	}
	return rec, nil
}

func saveCommentRecord(vaultPath string, rec CommentRecord) error {
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal comment record: %w", err)
	}
	path := commentFilePath(vaultPath)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("write comment record: %w", err)
	}
	return nil
}

// SetComment associates a comment with a vault key.
func SetComment(vaultPath, key, comment string) error {
	rec, err := LoadCommentRecord(vaultPath)
	if err != nil {
		return err
	}
	rec.Comments[key] = comment
	return saveCommentRecord(vaultPath, rec)
}

// RemoveComment removes the comment for the given key.
func RemoveComment(vaultPath, key string) error {
	rec, err := LoadCommentRecord(vaultPath)
	if err != nil {
		return err
	}
	delete(rec.Comments, key)
	return saveCommentRecord(vaultPath, rec)
}

// GetComment returns the comment for the given key, or empty string if none.
func GetComment(vaultPath, key string) (string, error) {
	rec, err := LoadCommentRecord(vaultPath)
	if err != nil {
		return "", err
	}
	return rec.Comments[key], nil
}

// ListComments returns all key→comment pairs sorted by key.
func ListComments(vaultPath string) ([]string, []string, error) {
	rec, err := LoadCommentRecord(vaultPath)
	if err != nil {
		return nil, nil, err
	}
	keys := make([]string, 0, len(rec.Comments))
	for k := range rec.Comments {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	comments := make([]string, len(keys))
	for i, k := range keys {
		comments[i] = rec.Comments[k]
	}
	return keys, comments, nil
}
