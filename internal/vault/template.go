package vault

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// templateVarRe matches {{VAR_NAME}} placeholders in a template string.
var templateVarRe = regexp.MustCompile(`\{\{([A-Z_][A-Z0-9_]*)\}\}`)

// TemplateResult holds the rendered output and any keys that were missing
// from the vault during rendering.
type TemplateResult struct {
	Output  string
	Missing []string
}

// RenderTemplate reads a template file from templatePath, substitutes every
// {{KEY}} placeholder with the corresponding value from v, and returns a
// TemplateResult.  Missing keys are collected rather than treated as errors so
// callers can decide how strict to be.
func RenderTemplate(v *Vault, templatePath string) (TemplateResult, error) {
	raw, err := os.ReadFile(templatePath)
	if err != nil {
		return TemplateResult{}, fmt.Errorf("template: read %q: %w", templatePath, err)
	}

	var missing []string
	seen := map[string]bool{}

	result := templateVarRe.ReplaceAllStringFunc(string(raw), func(match string) string {
		key := strings.TrimSuffix(strings.TrimPrefix(match, "{{"), "}}")
		val, ok := v.Get(key)
		if !ok {
			if !seen[key] {
				missing = append(missing, key)
				seen[key] = true
			}
			return match // leave placeholder intact
		}
		return val
	})

	return TemplateResult{Output: result, Missing: missing}, nil
}
