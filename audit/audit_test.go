package audit_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/user/vaultdiff/audit"
	"github.com/user/vaultdiff/diff"
)

func TestRecord_WritesJSONLine(t *testing.T) {
	f, err := os.CreateTemp("", "audit-*.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.Close()

	logger, err := audit.NewLogger(f.Name())
	if err != nil {
		t.Fatalf("NewLogger: %v", err)
	}
	changes := []diff.Change{
		{Key: "password", Type: diff.Modified, OldValue: "old", NewValue: "new"},
	}
	if err := logger.Record("secret/myapp", 1, 2, changes); err != nil {
		t.Fatalf("Record: %v", err)
	}
	logger.Close()

	data, err := os.ReadFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	var entry audit.Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if entry.Path != "secret/myapp" {
		t.Errorf("expected path secret/myapp, got %s", entry.Path)
	}
	if entry.VersionA != 1 || entry.VersionB != 2 {
		t.Errorf("unexpected versions: %d %d", entry.VersionA, entry.VersionB)
	}
	if len(entry.Changes) != 1 || entry.Changes[0].Key != "password" {
		t.Errorf("unexpected changes: %+v", entry.Changes)
	}
}

func TestNewLogger_Stdout(t *testing.T) {
	logger, err := audit.NewLogger("")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if err := logger.Close(); err != nil {
		t.Errorf("Close: %v", err)
	}
}

func TestNewLogger_InvalidPath(t *testing.T) {
	_, err := audit.NewLogger("/nonexistent/dir/audit.jsonl")
	if err == nil {
		t.Error("expected error for invalid path")
	}
}
