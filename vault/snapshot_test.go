package vault

import (
	"testing"

	"github.com/hashicorp/vault/api"
)

func newTestSnapshotter() *Snapshotter {
	client, _ := api.NewClient(api.DefaultConfig())
	return NewSnapshotter(client, "secret")
}

func TestNewSnapshotter_NotNil(t *testing.T) {
	s := newTestSnapshotter()
	if s == nil {
		t.Fatal("expected non-nil Snapshotter")
	}
}

func TestSnapshotter_Fields(t *testing.T) {
	s := newTestSnapshotter()
	if s.mount != "secret" {
		t.Errorf("expected mount 'secret', got %q", s.mount)
	}
	if s.client == nil {
		t.Error("expected non-nil client")
	}
}

func TestSnapshot_NilClient(t *testing.T) {
	s := &Snapshotter{client: nil, mount: "secret"}
	_, err := s.Snapshot("myapp")
	if err == nil {
		t.Error("expected error for nil client")
	}
}

func TestSnapshotEntry_Fields(t *testing.T) {
	e := SnapshotEntry{
		Path:    "myapp/config",
		Version: 3,
		Data:    map[string]interface{}{"key": "value"},
	}
	if e.Path != "myapp/config" {
		t.Errorf("unexpected path: %s", e.Path)
	}
	if e.Version != 3 {
		t.Errorf("unexpected version: %d", e.Version)
	}
	if e.Data["key"] != "value" {
		t.Errorf("unexpected data value")
	}
}
