package vault

import (
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

func newTestSigner() *Signer {
	cfg := vaultapi.DefaultConfig()
	client, _ := vaultapi.NewClient(cfg)
	return NewSigner(client)
}

func TestNewSigner_NotNil(t *testing.T) {
	s := newTestSigner()
	if s == nil {
		t.Fatal("expected non-nil Signer")
	}
}

func TestSigner_Fields(t *testing.T) {
	s := newTestSigner()
	if s.client == nil {
		t.Error("expected client to be set")
	}
}

func TestSign_NilClient(t *testing.T) {
	s := &Signer{client: nil}
	_, err := s.Sign("secret", "myapp/config", 1, "hmac-key")
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestSignResult_Fields(t *testing.T) {
	r := &SignResult{
		Path:      "myapp/config",
		Version:   3,
		Signature: "abc123",
	}
	if r.Path != "myapp/config" {
		t.Errorf("unexpected Path: %s", r.Path)
	}
	if r.Version != 3 {
		t.Errorf("unexpected Version: %d", r.Version)
	}
	if r.Signature != "abc123" {
		t.Errorf("unexpected Signature: %s", r.Signature)
	}
}

func TestComputeSignature_Deterministic(t *testing.T) {
	data := map[string]interface{}{
		"username": "admin",
		"password": "s3cr3t",
	}
	sig1 := computeSignature(data, "test-key")
	sig2 := computeSignature(data, "test-key")
	if sig1 != sig2 {
		t.Errorf("signatures differ: %s vs %s", sig1, sig2)
	}
}

func TestComputeSignature_DifferentKeys(t *testing.T) {
	data := map[string]interface{}{"foo": "bar"}
	sig1 := computeSignature(data, "key-a")
	sig2 := computeSignature(data, "key-b")
	if sig1 == sig2 {
		t.Error("expected different signatures for different HMAC keys")
	}
}

func TestComputeSignature_OrderIndependent(t *testing.T) {
	data1 := map[string]interface{}{"a": "1", "b": "2"}
	data2 := map[string]interface{}{"b": "2", "a": "1"}
	if computeSignature(data1, "k") != computeSignature(data2, "k") {
		t.Error("signatures should be order-independent")
	}
}
