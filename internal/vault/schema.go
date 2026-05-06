package vault

import (
	"fmt"
	"regexp"
)

// SchemaRule defines a validation rule for a specific key pattern.
type SchemaRule struct {
	KeyPattern  string `json:"key_pattern"`
	Required    bool   `json:"required"`
	ValueRegexp string `json:"value_regexp,omitempty"`
	Description string `json:"description,omitempty"`
}

// Schema holds a collection of rules used to validate vault contents.
type Schema struct {
	Rules []SchemaRule `json:"rules"`
}

// SchemaViolation describes a single rule violation found during validation.
type SchemaViolation struct {
	Key     string
	Rule    SchemaRule
	Message string
}

// ValidateSchema checks the vault against the provided schema and returns
// any violations found. An empty slice means the vault is compliant.
func ValidateSchema(v *Vault, s Schema) ([]SchemaViolation, error) {
	var violations []SchemaViolation

	for _, rule := range s.Rules {
		keyRe, err := regexp.Compile(rule.KeyPattern)
		if err != nil {
			return nil, fmt.Errorf("invalid key_pattern %q: %w", rule.KeyPattern, err)
		}

		var valRe *regexp.Regexp
		if rule.ValueRegexp != "" {
			valRe, err = regexp.Compile(rule.ValueRegexp)
			if err != nil {
				return nil, fmt.Errorf("invalid value_regexp %q: %w", rule.ValueRegexp, err)
			}
		}

		matched := false
		for _, key := range v.Keys() {
			if !keyRe.MatchString(key) {
				continue
			}
			matched = true
			if valRe != nil {
				val, _ := v.Get(key)
				if !valRe.MatchString(val) {
					violations = append(violations, SchemaViolation{
						Key:  key,
						Rule: rule,
						Message: fmt.Sprintf("value does not match pattern %q", rule.ValueRegexp),
					})
				}
			}
		}

		if rule.Required && !matched {
			violations = append(violations, SchemaViolation{
				Key:     rule.KeyPattern,
				Rule:    rule,
				Message: fmt.Sprintf("required key matching %q is missing", rule.KeyPattern),
			})
		}
	}

	return violations, nil
}
