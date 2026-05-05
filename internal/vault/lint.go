package vault

import (
	"fmt"
	"regexp"
	"strings"
)

// LintIssue describes a single linting problem found in a vault.
type LintIssue struct {
	Key     string
	Message string
}

func (l LintIssue) String() string {
	return fmt.Sprintf("%s: %s", l.Key, l.Message)
}

var validKeyPattern = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// Lint inspects the vault's keys and values for common problems and returns
// a slice of LintIssues. An empty slice means the vault is clean.
//
// Rules enforced:
//   - Keys must match [A-Z][A-Z0-9_]* (POSIX convention).
//   - Values must not be empty.
//   - Values must not contain unescaped newlines.
func Lint(v *Vault) []LintIssue {
	var issues []LintIssue

	for _, key := range v.Keys() {
		val, _ := v.Get(key)

		if !validKeyPattern.MatchString(key) {
			issues = append(issues, LintIssue{
				Key:     key,
				Message: "key does not follow POSIX naming convention (expected [A-Z][A-Z0-9_]*)",
			})
		}

		if strings.TrimSpace(val) == "" {
			issues = append(issues, LintIssue{
				Key:     key,
				Message: "value is empty or blank",
			})
		}

		if strings.ContainsAny(val, "\n\r") {
			issues = append(issues, LintIssue{
				Key:     key,
				Message: "value contains a newline character",
			})
		}
	}

	return issues
}
