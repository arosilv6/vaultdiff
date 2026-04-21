package vault

import (
	"testing"
)

func newTestValidator(requiredKeys []string) *Validator {
	return NewValidator(nil, "secret", requiredKeys)
}

func TestNewValidator_NotNil(t *testing.T) {
	v := newTestValidator([]string{"key1"})
	if v == nil {
		t.Fatal("expected non-nil Validator")
	}
}

func TestValidator_Fields(t *testing.T) {
	keys := []string{"username", "password"}
	v := NewValidator(nil, "kv", keys)

	if v.mount != "kv" {
		t.Errorf("expected mount kv, got %s", v.mount)
	}
	if len(v.requiredKeys) != 2 {
		t.Errorf("expected 2 required keys, got %d", len(v.requiredKeys))
	}
}

func TestValidate_NilClient(t *testing.T) {
	v := newTestValidator([]string{"host"})
	_, err := v.Validate("myapp/config", 1)
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestValidationResult_Fields(t *testing.T) {
	r := &ValidationResult{
		Path:    "myapp/db",
		Version: 3,
		Missing: []string{"password"},
		Valid:   false,
	}

	if r.Path != "myapp/db" {
		t.Errorf("unexpected path: %s", r.Path)
	}
	if r.Version != 3 {
		t.Errorf("unexpected version: %d", r.Version)
	}
	if r.Valid {
		t.Error("expected Valid to be false")
	}
	if len(r.Missing) != 1 || r.Missing[0] != "password" {
		t.Errorf("unexpected missing keys: %v", r.Missing)
	}
}

func TestValidationResult_Valid(t *testing.T) {
	r := &ValidationResult{
		Path:    "myapp/config",
		Version: 1,
		Missing: []string{},
		Valid:   true,
	}

	if !r.Valid {
		t.Error("expected Valid to be true")
	}
	if len(r.Missing) != 0 {
		t.Errorf("expected no missing keys, got %v", r.Missing)
	}
}
