package vault

import (
	"fmt"
	"regexp"
	"strings"
)

// RedactRule defines a pattern and its replacement.
type RedactRule struct {
	Name    string
	Pattern *regexp.Regexp
	Replace string
}

// RedactResult holds the output of a redaction pass.
type RedactResult struct {
	Path        string
	Version     int
	Redacted    map[string]string
	MatchedKeys []string
}

// Redactor applies redaction rules to secret data.
type Redactor struct {
	client *Client
	rules  []RedactRule
}

// NewRedactor creates a Redactor with the given client and built-in rules.
func NewRedactor(client *Client) *Redactor {
	defaultRules := []RedactRule{
		{
			Name:    "password",
			Pattern: regexp.MustCompile(`(?i)pass(word)?`),
			Replace: "[REDACTED]",
		},
		{
			Name:    "token",
			Pattern: regexp.MustCompile(`(?i)token|secret|key|api[_-]?key`),
			Replace: "[REDACTED]",
		},
		{
			Name:    "credit_card",
			Pattern: regexp.MustCompile(`\b\d{4}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}\b`),
			Replace: "[CC-REDACTED]",
		},
	}
	return &Redactor{client: client, rules: defaultRules}
}

// AddRule appends a custom redaction rule.
func (r *Redactor) AddRule(rule RedactRule) {
	r.rules = append(r.rules, rule)
}

// Redact reads a secret version and returns a redacted copy.
func (r *Redactor) Redact(path string, version int) (*RedactResult, error) {
	if r.client == nil {
		return nil, fmt.Errorf("redactor: nil client")
	}
	secret, err := r.client.Logical().Read(fmt.Sprintf("%s?version=%d", path, version))
	if err != nil {
		return nil, fmt.Errorf("redactor: read %s@%d: %w", path, version, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("redactor: no data at %s version %d", path, version)
	}
	data, _ := secret.Data["data"].(map[string]interface{})
	redacted := make(map[string]string, len(data))
	var matched []string
	for k, v := range data {
		val := fmt.Sprintf("%v", v)
		origVal := val
		for _, rule := range r.rules {
			if rule.Pattern.MatchString(k) {
				val = rule.Replace
				break
			}
			val = rule.Pattern.ReplaceAllString(val, rule.Replace)
		}
		redacted[k] = val
		if val != origVal || strings.Contains(val, "REDACTED") {
			matched = append(matched, k)
		}
	}
	return &RedactResult{
		Path:        path,
		Version:     version,
		Redacted:    redacted,
		MatchedKeys: matched,
	}, nil
}
