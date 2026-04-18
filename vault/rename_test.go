package vault

import (
	"testing"

	"github.com/hashicorp/vault/api"
)

func newTestRenamer(handler func(path string, data map[string]interface{}) (*api.Secret, error)) *Renamer {
	client, _ := api.NewClient(api.DefaultConfig())
	r := NewRenamer(client, "secret")
	return r
}

func TestNewRenamer_NotNil(t *testing.T) {
	client, _ := api.NewClient(api.DefaultConfig())
	r := NewRenamer(client, "secret")
	if r == nil {
		t.Fatal("expected non-nil Renamer")
	}
}

func TestRenamer_Fields(t *testing.T) {
	client, _ := api.NewClient(api.DefaultConfig())
	r := NewRenamer(client, "mymount")
	if r.mount != "mymount" {
		t.Errorf("expected mount mymount, got %s", r.mount)
	}
	if r.client != client {
		t.Error("expected client to be set")
	}
}

func TestRenamer_NilClientHandled(t *testing.T) {
	r := &Renamer{client: nil, mount: "secret"}
	if r == nil {
		t.Fatal("struct should be constructable")
	}
}
