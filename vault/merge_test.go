package vault

import (
	"testing"

	hashivault "github.com/hashicorp/vault/api"
)

func newTestMerger() *Merger {
	cfg := hashivault.DefaultConfig()
	client, _ := hashivault.NewClient(cfg)
	return NewMerger(client, "secret")
}

func TestNewMerger_NotNil(t *testing.T) {
	m := newTestMerger()
	if m == nil {
		t.Fatal("expected non-nil Merger")
	}
}

func TestMerger_Fields(t *testing.T) {
	m := newTestMerger()
	if m.client == nil {
		t.Error("expected client to be set")
	}
	if m.mount != "secret" {
		t.Errorf("expected mount 'secret', got %q", m.mount)
	}
}

func TestMerge_NilClient(t *testing.T) {
	m := &Merger{client: nil, mount: "secret"}
	_, err := m.Merge(nil, "myapp/config", 1, 2) //nolint:staticcheck
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestMergeResult_Fields(t *testing.T) {
	r := &MergeResult{
		Path:         "myapp/config",
		BaseVersion:  1,
		OtherVersion: 2,
		MergedKeys:   []string{"key1", "key2"},
		Conflicts:    []string{"key1"},
	}
	if r.Path != "myapp/config" {
		t.Errorf("unexpected Path: %s", r.Path)
	}
	if r.BaseVersion != 1 {
		t.Errorf("unexpected BaseVersion: %d", r.BaseVersion)
	}
	if r.OtherVersion != 2 {
		t.Errorf("unexpected OtherVersion: %d", r.OtherVersion)
	}
	if len(r.MergedKeys) != 2 {
		t.Errorf("expected 2 merged keys, got %d", len(r.MergedKeys))
	}
	if len(r.Conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d", len(r.Conflicts))
	}
}

func TestMerge_InvalidVersion(t *testing.T) {
	m := &Merger{client: nil, mount: "secret"}
	_, err := m.Merge(nil, "", 0, 0) //nolint:staticcheck
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}
