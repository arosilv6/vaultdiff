package vault

import (
	"strings"
	"testing"
)

func newTestFormatter(format OutputFormat) *Formatter {
	return NewFormatter(format)
}

func TestNewFormatter_NotNil(t *testing.T) {
	f := newTestFormatter(FormatJSON)
	if f == nil {
		t.Fatal("expected non-nil Formatter")
	}
}

func TestFormatter_Fields(t *testing.T) {
	f := newTestFormatter(FormatTable)
	if f.format != FormatTable {
		t.Errorf("expected format %q, got %q", FormatTable, f.format)
	}
}

func TestFormatData_NilData(t *testing.T) {
	f := newTestFormatter(FormatJSON)
	_, err := f.FormatData(nil)
	if err == nil {
		t.Fatal("expected error for nil data")
	}
}

func TestFormatData_JSON(t *testing.T) {
	f := newTestFormatter(FormatJSON)
	data := map[string]interface{}{"key": "value"}
	out, err := f.FormatData(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "\"key\"") {
		t.Errorf("expected JSON output to contain key, got: %s", out)
	}
}

func TestFormatData_Table(t *testing.T) {
	f := newTestFormatter(FormatTable)
	data := map[string]interface{}{"username": "admin"}
	out, err := f.FormatData(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "KEY") || !strings.Contains(out, "VALUE") {
		t.Errorf("expected table header in output, got: %s", out)
	}
	if !strings.Contains(out, "username") {
		t.Errorf("expected key in table output, got: %s", out)
	}
}

func TestFormatData_YAML(t *testing.T) {
	f := newTestFormatter(FormatYAML)
	data := map[string]interface{}{"token": "abc123"}
	out, err := f.FormatData(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "token: abc123") {
		t.Errorf("expected YAML output, got: %s", out)
	}
}

func TestFormatData_UnsupportedFormat(t *testing.T) {
	f := newTestFormatter(OutputFormat("xml"))
	_, err := f.FormatData(map[string]interface{}{"k": "v"})
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}
