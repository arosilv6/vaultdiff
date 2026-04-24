package vault

import (
	"testing"

	hashivault "github.com/hashicorp/vault/api"
)

func newTestPinner() *Pinner {
	cfg := hashivault.DefaultConfig()
	client, _ := hashivault.NewClient(cfg)
	return NewPinner(client, "secret")
}

func TestNewPinner_NotNil(t *testing.T) {
	p := newTestPinner()
	if p == nil {
		t.Fatal("expected non-nil Pinner")
	}
}

func TestPinner_Fields(t *testing.T) {
	p := newTestPinner()
	if p.mount != "secret" {
		t.Errorf("expected mount 'secret', got %q", p.mount)
	}
	if p.client == nil {
		t.Error("expected non-nil client")
	}
}

func TestPin_NilClient(t *testing.T) {
	p := &Pinner{client: nil, mount: "secret"}
	_, err := p.Pin(nil, "myapp/config", 2) //nolint:staticcheck
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestUnpin_NilClient(t *testing.T) {
	p := &Pinner{client: nil, mount: "secret"}
	_, err := p.Unpin(nil, "myapp/config") //nolint:staticcheck
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestPin_InvalidVersion(t *testing.T) {
	p := newTestPinner()
	_, err := p.Pin(nil, "myapp/config", 0) //nolint:staticcheck
	if err == nil {
		t.Fatal("expected error for version < 1")
	}
}

func TestPinResult_Fields(t *testing.T) {
	r := &PinResult{Path: "myapp/config", Version: 3, Pinned: true}
	if r.Path != "myapp/config" {
		t.Errorf("unexpected Path: %s", r.Path)
	}
	if r.Version != 3 {
		t.Errorf("unexpected Version: %d", r.Version)
	}
	if !r.Pinned {
		t.Error("expected Pinned to be true")
	}
}
