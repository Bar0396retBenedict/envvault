package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envvault/internal/vault"
)

func init() {
	tagCmd := &cobra.Command{
		Use:   "tag",
		Short: "Manage tags for vault keys",
	}

	tagsFile := tagCmd.PersistentFlags().String("tags-file", "tags.json", "path to tags index file")

	addCmd := &cobra.Command{
		Use:   "add <tag> <key>",
		Short: "Add a key to a tag",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTagAdd(*tagsFile, args[0], args[1])
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <tag> <key>",
		Short: "Remove a key from a tag",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTagRemove(*tagsFile, args[0], args[1])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list [tag]",
		Short: "List tags or keys within a tag",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTagList(*tagsFile, args)
		},
	}

	tagCmd.AddCommand(addCmd, removeCmd, listCmd)
	rootCmd.AddCommand(tagCmd)
}

func loadTagsFile(path string) (vault.Tags, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return make(vault.Tags), nil
	}
	if err != nil {
		return nil, err
	}
	var tags vault.Tags
	if err := json.Unmarshal(data, &tags); err != nil {
		return nil, fmt.Errorf("parse tags file: %w", err)
	}
	return tags, nil
}

func saveTagsFile(path string, tags vault.Tags) error {
	data, err := json.MarshalIndent(tags, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func runTagAdd(tagsFile, tag, key string) error {
	tags, err := loadTagsFile(tagsFile)
	if err != nil {
		return err
	}
	tags = vault.AddTag(tags, tag, key)
	if err := saveTagsFile(tagsFile, tags); err != nil {
		return err
	}
	fmt.Printf("tagged %q with %q\n", key, tag)
	return nil
}

func runTagRemove(tagsFile, tag, key string) error {
	tags, err := loadTagsFile(tagsFile)
	if err != nil {
		return err
	}
	tags = vault.RemoveTag(tags, tag, key)
	if err := saveTagsFile(tagsFile, tags); err != nil {
		return err
	}
	fmt.Printf("removed %q from tag %q\n", key, tag)
	return nil
}

func runTagList(tagsFile string, args []string) error {
	tags, err := loadTagsFile(tagsFile)
	if err != nil {
		return err
	}
	if len(args) == 1 {
		keys, err := vault.KeysForTag(tags, args[0])
		if err != nil {
			return err
		}
		fmt.Println(strings.Join(keys, "\n"))
		return nil
	}
	for _, t := range vault.ListTags(tags) {
		fmt.Println(t)
	}
	return nil
}
