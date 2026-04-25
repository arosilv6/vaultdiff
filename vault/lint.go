package vault

import (
	"context"
	"fmt"
	"strings"

	vaultapi "github.com/hashicorp/vault/api"
)

// LintIssue represents a single linting problem found in a secret.
type LintIssue struct {
	Path    string
	Key     string
	Message string
	Severity string // "warn" or "error"
}

// LintResult holds all issues found during linting.
type LintResult struct {
	Path   string
	Issues []LintIssue
}

// Linter checks secret paths for common issues.
type Linter struct {
	client *vaultapi.Client
	mount  string
}

// NewLinter creates a new Linter.
func NewLinter(client *vaultapi.Client, mount string) *Linter {
	return &Linter{client: client, mount: mount}
}

// Lint reads the latest version of a secret and checks for issues.
func (l *Linter) Lint(ctx context.Context, path string) (*LintResult, error) {
	if l.client == nil {
		return nil, fmt.Errorf("vault client is nil")
	}

	secretPath := fmt.Sprintf("%s/data/%s", l.mount, path)
	secret, err := l.client.Logical().ReadWithContext(ctx, secretPath)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	result := &LintResult{Path: path}

	if secret == nil || secret.Data == nil {
		result.Issues = append(result.Issues, LintIssue{
			Path:     path,
			Message:  "secret not found or empty",
			Severity: "error",
		})
		return result, nil
	}

	data, _ := secret.Data["data"].(map[string]interface{})
	for k, v := range data {
		if strings.TrimSpace(k) != k {
			result.Issues = append(result.Issues, LintIssue{Path: path, Key: k, Message: "key has leading/trailing whitespace", Severity: "warn"})
		}
		if v == nil || v == "" {
			result.Issues = append(result.Issues, LintIssue{Path: path, Key: k, Message: "value is empty", Severity: "warn"})
		}
		if strings.ToLower(k) != k {
			result.Issues = append(result.Issues, LintIssue{Path: path, Key: k, Message: "key is not lowercase", Severity: "warn"})
		}
	}

	return result, nil
}
