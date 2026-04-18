package vault

import (
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

func newTestRollbacker() *Rollbacker {
	client, _ := vaultapi.NewClient(vaultapi.DefaultConfig())
	return NewRollbacker(client, "secret")
}

func TestNewRollbacker_NotNil(t *testing.T) {
	r := newTestRollbacker()
	if r == nil {
		t.Fatal("expected non-nil Rollbacker")
	}
}

func TestRollbacker_Fields(t *testing.T) {
	r := newTestRollbacker()
	if r.mount != "secret" {
		t.Errorf("expected mount 'secret', got %s", r.mount)
	}
	if r.client == nil {
		t.Error("expected non-nil client")
	}
}

func TestRollback_NilClient(t *testing.T) {
	r := &Rollbacker{client: nil, mount: "secret"}
	_, err := r.Rollback(nil, "myapp/config", 2)
	if err == nil {
		t.Fatal("expected error with nil client")
	}
}

func TestRollbackResult_Fields(t *testing.T) {
	res := &RollbackResult{
		Path:        "myapp/config",
		FromVersion: 5,
		ToVersion:   2,
	}
	if res.Path != "myapp/config" {
		t.Errorf("unexpected path: %s", res.Path)
	}
	if res.FromVersion != 5 {
		t.Errorf("unexpected FromVersion: %d", res.FromVersion)
	}
	if res.ToVersion != 2 {
		t.Errorf("unexpected ToVersion: %d", res.ToVersion)
	}
}
