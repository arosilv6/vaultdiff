package vault

import (
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

func newTestPatcher() *Patcher {
	client, _ := vaultapi.NewClient(vaultapi.DefaultConfig())
	return NewPatcher(client, "secret")
}

func TestNewPatcher_NotNil(t *testing.T) {
	p := newTestPatcher()
	if p == nil {
		t.Fatal("expected non-nil Patcher")
	}
}

func TestPatcher_Fields(t *testing.T) {
	p := newTestPatcher()
	if p.mount != "secret" {
		t.Errorf("expected mount 'secret', got %q", p.mount)
	}
	if p.client == nil {
		t.Error("expected non-nil client")
	}
}

func TestPatch_NilClient(t *testing.T) {
	p := &Patcher{client: nil, mount: "secret"}
	_, err := p.Patch(nil, "myapp/config", map[string]string{"key": "val"})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestPatch_EmptyUpdates(t *testing.T) {
	p := newTestPatcher()
	_, err := p.Patch(nil, "myapp/config", map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty updates")
	}
}

func TestPatchResult_Fields(t *testing.T) {
	r := &PatchResult{
		Path:        "myapp/config",
		Version:     3,
		KeysPatched: []string{"DB_PASS", "API_KEY"},
	}
	if r.Path != "myapp/config" {
		t.Errorf("unexpected Path: %s", r.Path)
	}
	if r.Version != 3 {
		t.Errorf("unexpected Version: %d", r.Version)
	}
	if len(r.KeysPatched) != 2 {
		t.Errorf("expected 2 keys patched, got %d", len(r.KeysPatched))
	}
}
