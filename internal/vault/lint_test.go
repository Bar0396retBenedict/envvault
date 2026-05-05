package vault

import (
	"testing"
)

func makeLintVault(t *testing.T) *Vault {
	t.Helper()
	v, err := New(t.TempDir()+"/lint.vault", "passphrase")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return v
}

func TestLintCleanVault(t *testing.T) {
	v := makeLintVault(t)
	v.Set("DATABASE_URL", "postgres://localhost/db")
	v.Set("API_KEY", "abc123")

	issues := Lint(v)
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d: %v", len(issues), issues)
	}
}

func TestLintInvalidKeyName(t *testing.T) {
	v := makeLintVault(t)
	v.Set("invalid-key", "value")
	v.Set("1STARTS_WITH_DIGIT", "value")
	v.Set("lowercase", "value")

	issues := Lint(v)
	if len(issues) != 3 {
		t.Errorf("expected 3 issues, got %d: %v", len(issues), issues)
	}
}

func TestLintEmptyValue(t *testing.T) {
	v := makeLintVault(t)
	v.Set("EMPTY_VAR", "")
	v.Set("BLANK_VAR", "   ")

	issues := Lint(v)
	if len(issues) != 2 {
		t.Errorf("expected 2 issues, got %d: %v", len(issues), issues)
	}
	for _, iss := range issues {
		if iss.Message != "value is empty or blank" {
			t.Errorf("unexpected message: %s", iss.Message)
		}
	}
}

func TestLintNewlineInValue(t *testing.T) {
	v := makeLintVault(t)
	v.Set("MULTILINE", "line1\nline2")

	issues := Lint(v)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "MULTILINE" {
		t.Errorf("expected key MULTILINE, got %s", issues[0].Key)
	}
}

func TestLintMultipleProblemsOnSameKey(t *testing.T) {
	v := makeLintVault(t)
	// bad name AND newline in value
	v.Set("bad-key", "val\nue")

	issues := Lint(v)
	if len(issues) != 2 {
		t.Errorf("expected 2 issues for bad-key, got %d: %v", len(issues), issues)
	}
}

func TestLintIssueString(t *testing.T) {
	issue := LintIssue{Key: "FOO", Message: "something wrong"}
	want := "FOO: something wrong"
	if issue.String() != want {
		t.Errorf("String() = %q, want %q", issue.String(), want)
	}
}
