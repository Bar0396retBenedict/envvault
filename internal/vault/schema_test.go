package vault

import (
	"testing"
)

func makeSchemaVault(t *testing.T, pairs map[string]string) *Vault {
	t.Helper()
	v := New()
	for k, val := range pairs {
		if err := v.Set(k, val); err != nil {
			t.Fatalf("Set(%q): %v", k, err)
		}
	}
	return v
}

func TestValidateSchemaNoPviolations(t *testing.T) {
	v := makeSchemaVault(t, map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"PORT":         "8080",
	})
	s := Schema{
		Rules: []SchemaRule{
			{KeyPattern: "^DATABASE_URL$", Required: true, ValueRegexp: "^postgres://"},
			{KeyPattern: "^PORT$", Required: true, ValueRegexp: `^\d+$`},
		},
	}
	violations, err := ValidateSchema(v, s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(violations))
	}
}

func TestValidateSchemaRequiredMissing(t *testing.T) {
	v := makeSchemaVault(t, map[string]string{
		"PORT": "8080",
	})
	s := Schema{
		Rules: []SchemaRule{
			{KeyPattern: "^DATABASE_URL$", Required: true},
		},
	}
	violations, err := ValidateSchema(v, s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "^DATABASE_URL$" {
		t.Errorf("unexpected violation key: %q", violations[0].Key)
	}
}

func TestValidateSchemaValueMismatch(t *testing.T) {
	v := makeSchemaVault(t, map[string]string{
		"PORT": "not-a-number",
	})
	s := Schema{
		Rules: []SchemaRule{
			{KeyPattern: "^PORT$", Required: false, ValueRegexp: `^\d+$`},
		},
	}
	violations, err := ValidateSchema(v, s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "PORT" {
		t.Errorf("expected violation for PORT, got %q", violations[0].Key)
	}
}

func TestValidateSchemaInvalidKeyPattern(t *testing.T) {
	v := makeSchemaVault(t, map[string]string{"FOO": "bar"})
	s := Schema{
		Rules: []SchemaRule{
			{KeyPattern: "[", Required: false},
		},
	}
	_, err := ValidateSchema(v, s)
	if err == nil {
		t.Fatal("expected error for invalid key_pattern, got nil")
	}
}

func TestValidateSchemaPatternMatchesMultipleKeys(t *testing.T) {
	v := makeSchemaVault(t, map[string]string{
		"AWS_ACCESS_KEY": "AKIA123",
		"AWS_SECRET_KEY": "short",
	})
	s := Schema{
		Rules: []SchemaRule{
			{KeyPattern: "^AWS_", Required: true, ValueRegexp: `.{6,}`},
		},
	}
	violations, err := ValidateSchema(v, s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation for short value, got %d", len(violations))
	}
	if violations[0].Key != "AWS_SECRET_KEY" {
		t.Errorf("expected violation for AWS_SECRET_KEY, got %q", violations[0].Key)
	}
}
