package vault

import (
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

func newTestBlamer() *Blamer {
	cfg := vaultapi.DefaultConfig()
	client, _ := vaultapi.NewClient(cfg)
	return NewBlamer(client, "secret")
}

func TestNewBlamer_NotNil(t *testing.T) {
	b := newTestBlamer()
	if b == nil {
		t.Fatal("expected non-nil Blamer")
	}
}

func TestBlamer_Fields(t *testing.T) {
	b := newTestBlamer()
	if b.mount != "secret" {
		t.Errorf("expected mount 'secret', got %q", b.mount)
	}
	if b.client == nil {
		t.Error("expected non-nil client")
	}
}

func TestBlame_NilClient(t *testing.T) {
	b := &Blamer{client: nil, mount: "secret"}
	_, err := b.Blame(nil, "myapp/config") //nolint:staticcheck
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestBlameEntry_Fields(t *testing.T) {
	e := BlameEntry{
		Version:   3,
		CreatedBy: "token:root",
		Operation: "create",
		Deleted:   false,
	}
	if e.Version != 3 {
		t.Errorf("expected version 3, got %d", e.Version)
	}
	if e.CreatedBy != "token:root" {
		t.Errorf("unexpected CreatedBy: %s", e.CreatedBy)
	}
	if e.Operation != "create" {
		t.Errorf("unexpected Operation: %s", e.Operation)
	}
	if e.Deleted {
		t.Error("expected Deleted to be false")
	}
}

func TestBlameEntry_DeletedFlag(t *testing.T) {
	e := BlameEntry{Deleted: true}
	if !e.Deleted {
		t.Error("expected Deleted to be true")
	}
}
