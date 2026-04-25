package vault

import (
	"context"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

func newTestLinter() *Linter {
	cfg := vaultapi.DefaultConfig()
	client, _ := vaultapi.NewClient(cfg)
	return NewLinter(client, "secret")
}

func TestNewLinter_NotNil(t *testing.T) {
	l := newTestLinter()
	if l == nil {
		t.Fatal("expected non-nil Linter")
	}
}

func TestLinter_Fields(t *testing.T) {
	l := newTestLinter()
	if l.mount != "secret" {
		t.Errorf("expected mount 'secret', got %q", l.mount)
	}
	if l.client == nil {
		t.Error("expected non-nil client")
	}
}

func TestLint_NilClient(t *testing.T) {
	l := &Linter{client: nil, mount: "secret"}
	_, err := l.Lint(context.Background(), "myapp/config")
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestLintIssue_Fields(t *testing.T) {
	issue := LintIssue{
		Path:     "myapp/config",
		Key:      "API_KEY",
		Message:  "key is not lowercase",
		Severity: "warn",
	}
	if issue.Path != "myapp/config" {
		t.Errorf("unexpected Path: %s", issue.Path)
	}
	if issue.Severity != "warn" {
		t.Errorf("unexpected Severity: %s", issue.Severity)
	}
}

func TestLintResult_Fields(t *testing.T) {
	r := &LintResult{
		Path: "myapp/config",
		Issues: []LintIssue{
			{Key: "FOO", Message: "key is not lowercase", Severity: "warn"},
		},
	}
	if len(r.Issues) != 1 {
		t.Errorf("expected 1 issue, got %d", len(r.Issues))
	}
}
