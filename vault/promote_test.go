package vault

import (
	"testing"

	"github.com/hashicorp/vault/api"
)

func newTestPromoter() *Promoter {
	cfg := api.DefaultConfig()
	cfg.Address = "http://127.0.0.1:8200"
	client, _ := api.NewClient(cfg)
	return NewPromoter(client, "secret")
}

func TestNewPromoter_NotNil(t *testing.T) {
	p := newTestPromoter()
	if p == nil {
		t.Fatal("expected non-nil Promoter")
	}
}

func TestPromoter_Fields(t *testing.T) {
	p := newTestPromoter()
	if p.mount != "secret" {
		t.Errorf("expected mount 'secret', got %q", p.mount)
	}
	if p.client == nil {
		t.Error("expected non-nil client")
	}
}

func TestPromote_NilClient(t *testing.T) {
	p := &Promoter{client: nil, mount: "secret"}
	_, err := p.Promote("myapp/config", 2)
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestPromote_InvalidVersion(t *testing.T) {
	p := newTestPromoter()
	_, err := p.Promote("myapp/config", 0)
	if err == nil {
		t.Fatal("expected error for invalid version 0")
	}
}
