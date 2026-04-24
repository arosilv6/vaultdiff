package vault

import (
	"testing"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

func newTestArchiver() *Archiver {
	client, _ := vaultapi.NewClient(vaultapi.DefaultConfig())
	return NewArchiver(client, "secret")
}

func TestNewArchiver_NotNil(t *testing.T) {
	a := newTestArchiver()
	if a == nil {
		t.Fatal("expected non-nil Archiver")
	}
}

func TestArchiver_Fields(t *testing.T) {
	a := newTestArchiver()
	if a.mount != "secret" {
		t.Errorf("expected mount 'secret', got %q", a.mount)
	}
	if a.client == nil {
		t.Error("expected non-nil client")
	}
}

func TestArchive_NilClient(t *testing.T) {
	a := &Archiver{client: nil, mount: "secret"}
	_, err := a.Archive("myapp/config", 1, "archive/myapp/config-v1")
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestArchiveEntry_Fields(t *testing.T) {
	now := time.Now().UTC()
	entry := &ArchiveEntry{
		Path:       "archive/myapp/config-v1",
		Version:    1,
		Data:       map[string]interface{}{"key": "value"},
		ArchivedAt: now,
	}
	if entry.Path != "archive/myapp/config-v1" {
		t.Errorf("unexpected path: %s", entry.Path)
	}
	if entry.Version != 1 {
		t.Errorf("unexpected version: %d", entry.Version)
	}
	if entry.Data["key"] != "value" {
		t.Errorf("unexpected data value")
	}
	if entry.ArchivedAt != now {
		t.Errorf("unexpected archived_at time")
	}
}

func TestArchive_InvalidVersion(t *testing.T) {
	a := &Archiver{client: nil, mount: "secret"}
	_, err := a.Archive("myapp/config", 0, "archive/myapp/config-v0")
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}
