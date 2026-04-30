package vault

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ImportFormat represents the format of an import source.
type ImportFormat string

const (
	ImportDotEnv ImportFormat = "dotenv"
	ImportShell  ImportFormat = "shell"
)

// ImportResult summarises what happened during an import.
type ImportResult struct {
	Added    int
	Skipped  int
	Overwritten int
}

// ImportFromFile reads key=value pairs from path and merges them into v.
// Existing keys are overwritten only when overwrite is true.
func ImportFromFile(v *Vault, path string, format ImportFormat, overwrite bool) (ImportResult, error) {
	f, err := os.Open(path)
	if err != nil {
		return ImportResult{}, fmt.Errorf("import: open file: %w", err)
	}
	defer f.Close()

	var result ImportResult
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if format == ImportShell {
			line = strings.TrimPrefix(line, "export ")
		}
		key, value, ok := parseLine(line)
		if !ok {
			continue
		}
		if _, exists := v.Get(key); exists && !overwrite {
			result.Skipped++
			continue
		}
		if _, exists := v.Get(key); exists {
			result.Overwritten++
		} else {
			result.Added++
		}
		v.Set(key, value)
	}
	if err := scanner.Err(); err != nil {
		return ImportResult{}, fmt.Errorf("import: scan: %w", err)
	}
	return result, nil
}

// parseLine splits a KEY=VALUE line, stripping optional surrounding quotes.
func parseLine(line string) (key, value string, ok bool) {
	idx := strings.IndexByte(line, '=')
	if idx < 1 {
		return "", "", false
	}
	key = strings.TrimSpace(line[:idx])
	value = strings.TrimSpace(line[idx+1:])
	if len(value) >= 2 {
		if (value[0] == '"' && value[len(value)-1] == '"') ||
			(value[0] == '\'' && value[len(value)-1] == '\'') {
			value = value[1 : len(value)-1]
		}
	}
	return key, value, true
}
