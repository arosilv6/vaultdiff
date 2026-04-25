package audit

import (
	"context"
	"time"
)

// LintRecord captures a lint audit event.
type LintRecord struct {
	Timestamp  time.Time `json:"timestamp"`
	Operation  string    `json:"operation"`
	Path       string    `json:"path"`
	IssueCount int       `json:"issue_count"`
	HasErrors  bool      `json:"has_errors"`
}

// RecordLint writes a lint event to the audit logger.
func RecordLint(ctx context.Context, logger *Logger, path string, issueCount int, hasErrors bool) {
	if logger == nil {
		return
	}
	rec := LintRecord{
		Timestamp:  time.Now().UTC(),
		Operation:  "lint",
		Path:       path,
		IssueCount: issueCount,
		HasErrors:  hasErrors,
	}
	_ = logger.Record(ctx, rec)
}
