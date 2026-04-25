package vault

import (
	"regexp"
	"testing"
)

func newTestRedactor() *Redactor {
	return NewRedactor(nil)
}

func TestNewRedactor_NotNil(t *testing.T) {
	r := newTestRedactor()
	if r == nil {
		t.Fatal("expected non-nil Redactor")
	}
}

func TestRedactor_DefaultRules(t *testing.T) {
	r := newTestRedactor()
	if len(r.rules) == 0 {
		t.Fatal("expected default rules to be populated")
	}
}

func TestRedactor_AddRule(t *testing.T) {
	r := newTestRedactor()
	before := len(r.rules)
	r.AddRule(RedactRule{
		Name:    "custom",
		Pattern: regexp.MustCompile(`(?i)ssn`),
		Replace: "[SSN-REDACTED]",
	})
	if len(r.rules) != before+1 {
		t.Fatalf("expected %d rules, got %d", before+1, len(r.rules))
	}
}

func TestRedact_NilClient(t *testing.T) {
	r := newTestRedactor()
	_, err := r.Redact("secret/data/test", 1)
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestRedactResult_Fields(t *testing.T) {
	res := &RedactResult{
		Path:        "secret/data/myapp",
		Version:     3,
		Redacted:    map[string]string{"password": "[REDACTED]", "host": "localhost"},
		MatchedKeys: []string{"password"},
	}
	if res.Path != "secret/data/myapp" {
		t.Errorf("unexpected Path: %s", res.Path)
	}
	if res.Version != 3 {
		t.Errorf("unexpected Version: %d", res.Version)
	}
	if res.Redacted["password"] != "[REDACTED]" {
		t.Errorf("expected password to be redacted")
	}
	if res.Redacted["host"] != "localhost" {
		t.Errorf("expected host to remain unredacted")
	}
	if len(res.MatchedKeys) != 1 || res.MatchedKeys[0] != "password" {
		t.Errorf("unexpected MatchedKeys: %v", res.MatchedKeys)
	}
}

func TestRedactRule_PatternMatch(t *testing.T) {
	rule := RedactRule{
		Name:    "token",
		Pattern: regexp.MustCompile(`(?i)token`),
		Replace: "[REDACTED]",
	}
	result := rule.Pattern.ReplaceAllString("my-token-value", rule.Replace)
	if result != "my-[REDACTED]-value" {
		t.Errorf("unexpected replacement result: %s", result)
	}
}
