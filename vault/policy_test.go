package vault

import (
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

func newTestPolicyChecker(t *testing.T) *PolicyChecker {
	t.Helper()
	cfg := vaultapi.DefaultConfig()
	cfg.Address = "http://127.0.0.1:8200"
	client, _ := vaultapi.NewClient(cfg)
	return NewPolicyChecker(client)
}

func TestNewPolicyChecker_NotNil(t *testing.T) {
	pc := newTestPolicyChecker(t)
	if pc == nil {
		t.Fatal("expected non-nil PolicyChecker")
	}
}

func TestPolicyChecker_Fields(t *testing.T) {
	pc := newTestPolicyChecker(t)
	if pc.client == nil {
		t.Error("expected client to be set")
	}
}

func TestCheckPath_NilClient(t *testing.T) {
	pc := &PolicyChecker{client: nil}
	_, err := pc.CheckPath(nil, "secret/data/foo") //nolint:staticcheck
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestCapabilityResult_Fields(t *testing.T) {
	r := CapabilityResult{
		Path:         "secret/data/foo",
		Capabilities: []string{"read", "list"},
	}
	if r.Path != "secret/data/foo" {
		t.Errorf("unexpected path: %s", r.Path)
	}
	if len(r.Capabilities) != 2 {
		t.Errorf("expected 2 capabilities, got %d", len(r.Capabilities))
	}
}

func TestHasCapability_True(t *testing.T) {
	r := CapabilityResult{
		Path:         "secret/data/foo",
		Capabilities: []string{"read", "update"},
	}
	if !r.HasCapability("read") {
		t.Error("expected HasCapability('read') to be true")
	}
}

func TestHasCapability_False(t *testing.T) {
	r := CapabilityResult{
		Path:         "secret/data/foo",
		Capabilities: []string{"read"},
	}
	if r.HasCapability("delete") {
		t.Error("expected HasCapability('delete') to be false")
	}
}
