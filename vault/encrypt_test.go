package vault

import (
	"context"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

func newTestEncrypter(t *testing.T) *Encrypter {
	t.Helper()
	client, _ := vaultapi.NewClient(vaultapi.DefaultConfig())
	return NewEncrypter(client, "transit", "my-key")
}

func TestNewEncrypter_NotNil(t *testing.T) {
	e := newTestEncrypter(t)
	if e == nil {
		t.Fatal("expected non-nil Encrypter")
	}
}

func TestEncrypter_Fields(t *testing.T) {
	e := newTestEncrypter(t)
	if e.mount != "transit" {
		t.Errorf("expected mount 'transit', got %q", e.mount)
	}
	if e.keyName != "my-key" {
		t.Errorf("expected keyName 'my-key', got %q", e.keyName)
	}
}

func TestEncrypt_NilClient(t *testing.T) {
	e := &Encrypter{client: nil, mount: "transit", keyName: "k"}
	_, err := e.Encrypt(context.Background(), "dGVzdA==")
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestDecrypt_NilClient(t *testing.T) {
	e := &Encrypter{client: nil, mount: "transit", keyName: "k"}
	_, err := e.Decrypt(context.Background(), "vault:v1:abc")
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestEncryptResult_Fields(t *testing.T) {
	r := &EncryptResult{KeyName: "k", Ciphertext: "vault:v1:xyz"}
	if r.KeyName != "k" {
		t.Errorf("unexpected KeyName: %s", r.KeyName)
	}
	if r.Ciphertext != "vault:v1:xyz" {
		t.Errorf("unexpected Ciphertext: %s", r.Ciphertext)
	}
}

func TestDecryptResult_Fields(t *testing.T) {
	r := &DecryptResult{KeyName: "k", Plaintext: "dGVzdA=="}
	if r.KeyName != "k" {
		t.Errorf("unexpected KeyName: %s", r.KeyName)
	}
	if r.Plaintext != "dGVzdA==" {
		t.Errorf("unexpected Plaintext: %s", r.Plaintext)
	}
}
